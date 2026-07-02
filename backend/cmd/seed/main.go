// Command seed idempotently loads the SAMPLE dataset (package
// application/internal/seeddata) into Postgres + MinIO so that a fresh
// `docker compose up` always lands the exact same competitions, media, users,
// tickets and drawn draws.
//
// Idempotency is achieved by deriving every row id deterministically from its
// natural key with a fixed UUIDv5 namespace (seeddata.Namespace). Re-running the
// seeder therefore upserts the same rows and re-uploads objects to the same
// deterministic MinIO keys — it never duplicates.
//
// It never mutates seeddata: that package is the single source of truth and is
// imported read-only. Missing image files are logged and skipped (soft failure);
// the process still exits non-zero if any upload failed, but DB errors are hard
// failures that abort immediately.
package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"application/internal/seeddata"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib" // register the "pgx" database/sql driver
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// drawStatusDrawn is the draws.status value for a completed draw (a winner has
// been picked). Kept local so the seeder does not depend on the draw entity pkg.
const drawStatusDrawn = "drawn"

// ownerTypeCompetition / ownerTypeUser are the media.owner_type discriminators.
const (
	ownerTypeCompetition = "competition"
	ownerTypeUser        = "user"
)

// seedTimeout bounds the whole run so a hung DB/MinIO connection cannot wedge
// the one-shot seed container forever.
const seedTimeout = 5 * time.Minute

// config holds the resolved runtime configuration read from the environment.
type config struct {
	postgresDSN    string
	minioEndpoint  string
	minioAccessKey string
	minioSecretKey string
	minioBucket    string
	minioUseSSL    bool
	assetsDir      string
}

func main() {
	// log with timestamps disabled — the compose logs already timestamp lines.
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatalf("seed FAILED: %v", err)
	}
}

// run performs the full seed. It returns a non-nil error on any hard failure
// (bad config, DB error) or if one or more image uploads failed.
func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), seedTimeout)
	defer cancel()

	// --- Postgres ---
	db, err := sql.Open("pgx", cfg.postgresDSN)
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	// --- MinIO ---
	mc, err := minio.New(cfg.minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.minioAccessKey, cfg.minioSecretKey, ""),
		Secure: cfg.minioUseSSL,
	})
	if err != nil {
		return fmt.Errorf("new minio client: %w", err)
	}

	if err := ensureBucket(ctx, mc, cfg.minioBucket); err != nil {
		return fmt.Errorf("ensure bucket: %w", err)
	}

	ns, err := uuid.Parse(seeddata.Namespace)
	if err != nil {
		return fmt.Errorf("parse seeddata namespace %q: %w", seeddata.Namespace, err)
	}

	s := &seeder{cfg: cfg, db: db, mc: mc, ns: ns}

	if err := s.seedCompetitions(ctx); err != nil {
		return err
	}

	if err := s.seedWinners(ctx); err != nil {
		return err
	}

	log.Printf(
		"seeded: %d competitions, %d media, %d users, %d tickets, %d draws",
		s.competitions, s.media, s.users, s.tickets, s.draws,
	)

	// Soft-fail: individual missing/failed uploads do not abort the run, but the
	// overall exit is non-zero so callers (compose, CI) notice.
	if s.uploadFailures > 0 {
		return fmt.Errorf("%d image upload(s) failed", s.uploadFailures)
	}

	log.Print("SEED OK")

	return nil
}

// seeder carries the shared clients plus running counters for the summary line.
type seeder struct {
	cfg *config
	db  *sql.DB
	mc  *minio.Client
	ns  uuid.UUID

	competitions   int
	media          int
	users          int
	tickets        int
	draws          int
	uploadFailures int
}

// seedCompetitions upserts every sample competition and uploads its hero image
// as a single position-0 media row.
func (s *seeder) seedCompetitions(ctx context.Context) error {
	now := time.Now().UTC()

	for _, comp := range seeddata.Competitions {
		id := s.competitionID(comp.Slug)

		// starts_at is fixed a day in the past; ends_at is derived from the
		// sample's EndsInHours so the frontend "ENDS …" badge stays live.
		startsAt := now.Add(-24 * time.Hour)
		endsAt := now.Add(time.Duration(comp.EndsInHours) * time.Hour)

		if err := s.upsertCompetition(ctx, competitionRow{
			id:               id,
			title:            comp.Title,
			slug:             comp.Slug,
			description:      comp.Description,
			prize:            comp.Prize,
			ticketPricePence: comp.TicketPricePence,
			ticketsTotal:     comp.TicketsTotal,
			ticketsSold:      comp.TicketsSold,
			status:           comp.Status,
			startsAt:         &startsAt,
			endsAt:           &endsAt,
		}); err != nil {
			return fmt.Errorf("upsert competition %q: %w", comp.Slug, err)
		}
		s.competitions++

		// The hero image lives under <assets>/images/comps/<Image>.
		srcPath := filepath.Join(s.cfg.assetsDir, "images", "comps", comp.Image)
		if err := s.uploadMedia(ctx, ownerTypeCompetition, id, comp.Image, srcPath, 0); err != nil {
			return err
		}
	}

	return nil
}

// seedWinners creates the CLOSED winners-archive competition and, for each
// sample winner, a user + ticket + avatar media + drawn draw hanging off it.
func (s *seeder) seedWinners(ctx context.Context) error {
	now := time.Now().UTC()

	// The archive is a single closed competition owning every draw so winners do
	// not pollute the live competition grids. It has no image and no end date.
	archiveID := s.competitionID(seeddata.WinnersArchiveSlug)
	startsAt := now.Add(-24 * time.Hour)

	if err := s.upsertCompetition(ctx, competitionRow{
		id:          archiveID,
		title:       "BOTB Winners",
		slug:        seeddata.WinnersArchiveSlug,
		description: "Archive of past BOTB winners.",
		prize:       "Past winners",
		status:      seeddata.StatusClosed,
		startsAt:    &startsAt,
		endsAt:      nil,
	}); err != nil {
		return fmt.Errorf("upsert winners archive: %w", err)
	}
	s.competitions++

	for _, w := range seeddata.Winners {
		userID := s.userID(w.Email)
		ticketID := uuid.NewSHA1(s.ns, []byte("ticket:"+w.Email))
		drawID := uuid.NewSHA1(s.ns, []byte("draw:"+w.Email))

		// user — tickets_owned 1, total_spent = the winning ticket's cost.
		if err := s.upsertUser(ctx, userID, w.Name, w.Email, 1, w.WonForPence); err != nil {
			return fmt.Errorf("upsert user %q: %w", w.Email, err)
		}
		s.users++

		// ticket — one entry on the archive competition (insert-once).
		if err := s.insertTicket(ctx, ticketID, archiveID, userID); err != nil {
			return fmt.Errorf("insert ticket for %q: %w", w.Email, err)
		}
		s.tickets++

		// avatar — <assets>/images/winners/<Image>, owned by the user.
		srcPath := filepath.Join(s.cfg.assetsDir, "images", "winners", w.Image)
		if err := s.uploadMedia(ctx, ownerTypeUser, userID, w.Image, srcPath, 0); err != nil {
			return err
		}

		// draw — a completed draw linking user + ticket + prize.
		if err := s.upsertDraw(ctx, drawRow{
			id:             drawID,
			competitionID:  archiveID,
			winnerUserID:   &userID,
			winnerTicketID: &ticketID,
			prize:          w.Prize,
			status:         drawStatusDrawn,
			drawnAt:        &now,
		}); err != nil {
			return fmt.Errorf("upsert draw for %q: %w", w.Email, err)
		}
		s.draws++
	}

	return nil
}

// -------------------- deterministic id derivation --------------------

func (s *seeder) competitionID(slug string) uuid.UUID {
	return uuid.NewSHA1(s.ns, []byte("competition:"+slug))
}

func (s *seeder) userID(email string) uuid.UUID {
	return uuid.NewSHA1(s.ns, []byte("user:"+email))
}

func (s *seeder) mediaID(ownerType string, ownerID uuid.UUID, filename string) uuid.UUID {
	return uuid.NewSHA1(s.ns, []byte("media:"+ownerType+":"+ownerID.String()+":"+filename))
}

// -------------------- media upload --------------------

// uploadMedia uploads srcPath to MinIO under a deterministic object key and
// upserts the matching media row. A missing source file is a SOFT failure: it is
// logged, counted, and skipped so the rest of the seed still completes.
func (s *seeder) uploadMedia(
	ctx context.Context,
	ownerType string,
	ownerID uuid.UUID,
	filename, srcPath string,
	position int,
) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		log.Printf("WARN: skipping media %q (owner %s/%s): %v", filename, ownerType, ownerID, err)
		s.uploadFailures++

		return nil
	}

	kind, contentType := classify(filename)

	// Deterministic key: re-running overwrites the same object rather than
	// creating a new one, keeping MinIO in lockstep with the derived media id.
	objectKey := fmt.Sprintf("%s/%s/%s", ownerType, ownerID, filename)

	_, err = s.mc.PutObject(
		ctx, s.cfg.minioBucket, objectKey,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Printf("WARN: upload failed for %q → %s: %v", srcPath, objectKey, err)
		s.uploadFailures++

		return nil
	}

	if err := s.upsertMedia(ctx, mediaRow{
		id:          s.mediaID(ownerType, ownerID, filename),
		ownerType:   ownerType,
		ownerID:     ownerID,
		kind:        kind,
		bucket:      s.cfg.minioBucket,
		objectKey:   objectKey,
		contentType: contentType,
		sizeBytes:   int64(len(data)),
		position:    position,
	}); err != nil {
		// A failed media row IS a hard error — the object is now orphaned and the
		// dataset is inconsistent, so abort rather than silently continue.
		return fmt.Errorf("upsert media %q: %w", objectKey, err)
	}
	s.media++

	return nil
}

// classify maps a filename extension to a media kind + MIME content type.
// Unknown extensions fall back to a generic binary image row.
func classify(filename string) (kind, contentType string) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".webp":
		return "image", "image/webp"
	case ".png":
		return "image", "image/png"
	case ".jpg", ".jpeg":
		return "image", "image/jpeg"
	case ".mp4":
		return "video", "video/mp4"
	case ".webm":
		return "video", "video/webm"
	default:
		return "image", "application/octet-stream"
	}
}

// -------------------- row structs + upserts --------------------

type competitionRow struct {
	id               uuid.UUID
	title            string
	slug             string
	description      string
	prize            string
	ticketPricePence int64
	ticketsTotal     int64
	ticketsSold      int64
	status           string
	startsAt         *time.Time
	endsAt           *time.Time
}

func (s *seeder) upsertCompetition(ctx context.Context, c competitionRow) error {
	const q = `INSERT INTO competitions
		(id, title, slug, description, prize, ticket_price_pence, tickets_total,
		 tickets_sold, status, starts_at, ends_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET
			title              = EXCLUDED.title,
			slug               = EXCLUDED.slug,
			description        = EXCLUDED.description,
			prize              = EXCLUDED.prize,
			ticket_price_pence = EXCLUDED.ticket_price_pence,
			tickets_total      = EXCLUDED.tickets_total,
			tickets_sold       = EXCLUDED.tickets_sold,
			status             = EXCLUDED.status,
			starts_at          = EXCLUDED.starts_at,
			ends_at            = EXCLUDED.ends_at,
			updated_at         = NOW()`

	_, err := s.db.ExecContext(ctx, q,
		c.id, c.title, c.slug, c.description, c.prize, c.ticketPricePence,
		c.ticketsTotal, c.ticketsSold, c.status, c.startsAt, c.endsAt,
	)

	return err
}

type mediaRow struct {
	id          uuid.UUID
	ownerType   string
	ownerID     uuid.UUID
	kind        string
	bucket      string
	objectKey   string
	contentType string
	sizeBytes   int64
	position    int
}

func (s *seeder) upsertMedia(ctx context.Context, m mediaRow) error {
	const q = `INSERT INTO media
		(id, owner_type, owner_id, kind, bucket, object_key, content_type,
		 size_bytes, width, height, duration_seconds, position)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,0,0,NULL,$9)
		ON CONFLICT (id) DO UPDATE SET
			owner_type   = EXCLUDED.owner_type,
			owner_id     = EXCLUDED.owner_id,
			kind         = EXCLUDED.kind,
			bucket       = EXCLUDED.bucket,
			object_key   = EXCLUDED.object_key,
			content_type = EXCLUDED.content_type,
			size_bytes   = EXCLUDED.size_bytes,
			position     = EXCLUDED.position`

	_, err := s.db.ExecContext(ctx, q,
		m.id, m.ownerType, m.ownerID, m.kind, m.bucket, m.objectKey,
		m.contentType, m.sizeBytes, m.position,
	)

	return err
}

func (s *seeder) upsertUser(
	ctx context.Context,
	id uuid.UUID,
	name, email string,
	ticketsOwned, totalSpentPence int64,
) error {
	const q = `INSERT INTO users
		(id, name, email, tickets_owned, total_spent_pence)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (id) DO UPDATE SET
			name              = EXCLUDED.name,
			email             = EXCLUDED.email,
			tickets_owned     = EXCLUDED.tickets_owned,
			total_spent_pence = EXCLUDED.total_spent_pence`

	_, err := s.db.ExecContext(ctx, q, id, name, email, ticketsOwned, totalSpentPence)

	return err
}

// insertTicket is insert-once: one deterministic ticket per winner is enough, so
// a conflict simply means it already exists.
func (s *seeder) insertTicket(ctx context.Context, id, competitionID, userID uuid.UUID) error {
	const q = `INSERT INTO tickets (id, competition_id, user_id)
		VALUES ($1,$2,$3)
		ON CONFLICT (id) DO NOTHING`

	_, err := s.db.ExecContext(ctx, q, id, competitionID, userID)

	return err
}

type drawRow struct {
	id             uuid.UUID
	competitionID  uuid.UUID
	winnerUserID   *uuid.UUID
	winnerTicketID *uuid.UUID
	prize          string
	status         string
	drawnAt        *time.Time
}

func (s *seeder) upsertDraw(ctx context.Context, d drawRow) error {
	const q = `INSERT INTO draws
		(id, competition_id, winner_user_id, winner_ticket_id, prize, status, drawn_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET
			competition_id   = EXCLUDED.competition_id,
			winner_user_id   = EXCLUDED.winner_user_id,
			winner_ticket_id = EXCLUDED.winner_ticket_id,
			prize            = EXCLUDED.prize,
			status           = EXCLUDED.status,
			drawn_at         = EXCLUDED.drawn_at,
			updated_at       = NOW()`

	_, err := s.db.ExecContext(ctx, q,
		d.id, d.competitionID, uuidPtr(d.winnerUserID), uuidPtr(d.winnerTicketID),
		d.prize, d.status, d.drawnAt,
	)

	return err
}

// uuidPtr renders a nullable uuid pointer as a driver value ([16]byte or nil).
// database/sql cannot handle *uuid.UUID directly, so nil → SQL NULL.
func uuidPtr(id *uuid.UUID) any {
	if id == nil {
		return nil
	}

	return *id
}

// -------------------- config + minio bootstrap --------------------

// ensureBucket creates the target bucket if it does not already exist so the
// seeder is self-sufficient even before the compose `mc` init container runs.
func ensureBucket(ctx context.Context, mc *minio.Client, bucket string) error {
	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		return mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}

	return nil
}

// loadConfig resolves the env contract, applying defaults and failing fast on a
// missing required DSN.
func loadConfig() (*config, error) {
	dsn := os.Getenv("APP_DATASOURCE_POSTGRES_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("APP_DATASOURCE_POSTGRES_DSN is required")
	}

	// ParseBool tolerates empty ("" → false) via the default branch.
	useSSL := false
	if raw := os.Getenv("SEED_MINIO_USE_SSL"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, fmt.Errorf("SEED_MINIO_USE_SSL %q: %w", raw, err)
		}
		useSSL = parsed
	}

	return &config{
		postgresDSN:    dsn,
		minioEndpoint:  envDefault("SEED_MINIO_ENDPOINT", "minio:9000"),
		minioAccessKey: envDefault("SEED_MINIO_ACCESS_KEY", "minioadmin"),
		minioSecretKey: envDefault("SEED_MINIO_SECRET_KEY", "minioadmin"),
		minioBucket:    envDefault("SEED_MINIO_BUCKET", "botb-media"),
		minioUseSSL:    useSSL,
		assetsDir:      envDefault("SEED_ASSETS_DIR", "/assets"),
	}, nil
}

// envDefault returns the env var value or fallback when it is unset/empty.
func envDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

// printConfig echoes the resolved configuration with secrets masked.
func printConfig(cfg *config) {
	log.Print("seed configuration:")
	log.Printf("  APP_DATASOURCE_POSTGRES_DSN = %s", maskDSN(cfg.postgresDSN))
	log.Printf("  SEED_MINIO_ENDPOINT         = %s", cfg.minioEndpoint)
	log.Printf("  SEED_MINIO_ACCESS_KEY       = %s", cfg.minioAccessKey)
	log.Printf("  SEED_MINIO_SECRET_KEY       = %s", mask(cfg.minioSecretKey))
	log.Printf("  SEED_MINIO_BUCKET           = %s", cfg.minioBucket)
	log.Printf("  SEED_MINIO_USE_SSL          = %t", cfg.minioUseSSL)
	log.Printf("  SEED_ASSETS_DIR             = %s", cfg.assetsDir)
}

// mask hides all but the first and last character of a secret.
func mask(secret string) string {
	if len(secret) <= 2 {
		return "****"
	}

	return secret[:1] + "****" + secret[len(secret)-1:]
}

// maskDSN replaces the password component of a postgres URL/keyword DSN so it is
// safe to print in logs.
func maskDSN(dsn string) string {
	// URL form: scheme://user:PASSWORD@host/db
	if at := strings.LastIndex(dsn, "@"); at != -1 {
		if sep := strings.Index(dsn, "://"); sep != -1 {
			creds := dsn[sep+3 : at]
			if colon := strings.Index(creds, ":"); colon != -1 {
				return dsn[:sep+3] + creds[:colon] + ":****" + dsn[at:]
			}
		}
	}

	return dsn
}
