package repo_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"
	"time"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/competition/repo"
	"application/internal/datasource"

	"github.com/google/uuid"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

// newRamsqlDB spins up an isolated ramsql DB with the competitions + media
// tables. ramsql lacks Postgres-isms (RETURNING, uuid type), so the schema uses
// portable types and these tests cover the read paths (Get/List) plus the media
// join; write paths (RETURNING) are covered by the biz mock tests.
func newRamsqlDB(t *testing.T) *datasource.PostgresDB {
	t.Helper()

	db, err := sql.Open("ramsql", "competition_"+uuid.NewString())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()

	_, err = db.ExecContext(ctx, `CREATE TABLE competitions (
		id TEXT PRIMARY KEY,
		title TEXT,
		slug TEXT,
		description TEXT,
		prize TEXT,
		ticket_price_pence BIGINT,
		tickets_total BIGINT,
		tickets_sold BIGINT,
		category_id TEXT,
		status TEXT,
		starts_at TIMESTAMP,
		ends_at TIMESTAMP,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `CREATE TABLE media (
		id TEXT PRIMARY KEY,
		owner_type TEXT,
		owner_id TEXT,
		kind TEXT,
		bucket TEXT,
		object_key TEXT,
		content_type TEXT,
		size_bytes BIGINT,
		width INT,
		height INT,
		duration_seconds DECIMAL,
		position INT,
		created_at TIMESTAMP
	)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `CREATE TABLE categories (
		id TEXT PRIMARY KEY,
		name TEXT,
		slug TEXT,
		created_at TIMESTAMP
	)`)
	require.NoError(t, err)

	// tickets + draws let the delete-safety guard be exercised end to end.
	_, err = db.ExecContext(ctx, `CREATE TABLE tickets (
		id TEXT PRIMARY KEY,
		competition_id TEXT,
		user_id TEXT
	)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `CREATE TABLE draws (
		id TEXT PRIMARY KEY,
		competition_id TEXT,
		status TEXT
	)`)
	require.NoError(t, err)

	return &datasource.PostgresDB{DB: db}
}

func seedCompetition(t *testing.T, db *datasource.PostgresDB, id uuid.UUID, slug, status string) {
	t.Helper()

	now := time.Now().UTC()

	_, err := db.ExecContext(
		context.Background(), `INSERT INTO competitions
		(id, title, slug, description, prize, ticket_price_pence, tickets_total,
		 tickets_sold, category_id, status, starts_at, ends_at, created_at, updated_at)
		VALUES ($1,$2,$3,'','Prize',125,1000,0,NULL,$4,$5,$6,$7,$8)`,
		id.String(), "Title "+slug, slug, status, now, now, now, now,
	)
	require.NoError(t, err)
}

func seedMedia(t *testing.T, db *datasource.PostgresDB, ownerID uuid.UUID, position int) {
	t.Helper()

	_, err := db.ExecContext(
		context.Background(), `INSERT INTO media
		(id, owner_type, owner_id, kind, bucket, object_key, content_type,
		 size_bytes, width, height, duration_seconds, position, created_at)
		VALUES ($1,'competition',$2,'image','botb-media','k','image/png',1,0,0,NULL,$3, NOW())`,
		uuid.NewString(), ownerID.String(), position,
	)
	require.NoError(t, err)
}

func newRepo(db *datasource.PostgresDB) biz.Repository {
	return repo.NewCompetition(slog.New(slog.NewTextHandler(io.Discard, nil)), db)
}

func TestRepo_Get_WithMedia(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedCompetition(t, db, id, "audi-rs3", "live")
	seedMedia(t, db, id, 1)
	seedMedia(t, db, id, 0)
	seedMedia(t, db, uuid.New(), 0) // different owner — must not leak

	got, err := r.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
	require.Equal(t, entity.StatusLive, got.Status)
	require.Len(t, got.Media, 2)
	require.Equal(t, 0, got.Media[0].Position) // ordered by position ASC
}

func TestRepo_Get_NotFound(t *testing.T) {
	db := newRamsqlDB(t)
	r := newRepo(db)

	_, err := r.Get(context.Background(), uuid.New())
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

func TestRepo_List_All(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	seedCompetition(t, db, uuid.New(), "one", "live")
	seedCompetition(t, db, uuid.New(), "two", "draft")

	got, err := r.List(ctx, biz.ListFilter{})
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestRepo_List_StatusFilter(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	live := uuid.New()
	seedCompetition(t, db, live, "live-one", "live")
	seedCompetition(t, db, uuid.New(), "draft-one", "draft")
	seedMedia(t, db, live, 0)

	status := entity.StatusLive
	got, err := r.List(ctx, biz.ListFilter{Status: &status})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "live-one", got[0].Slug)
	require.Len(t, got[0].Media, 1)
}

func seedTicket(t *testing.T, db *datasource.PostgresDB, competitionID uuid.UUID) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO tickets (id, competition_id, user_id) VALUES ($1, $2, $3)`,
		uuid.NewString(), competitionID.String(), uuid.NewString())
	require.NoError(t, err)
}

func seedDrawRow(t *testing.T, db *datasource.PostgresDB, competitionID uuid.UUID) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO draws (id, competition_id, status) VALUES ($1, $2, 'pending')`,
		uuid.NewString(), competitionID.String())
	require.NoError(t, err)
}

// A clean competition (no tickets, no draws) deletes, and its media object
// keys come back for the object-storage purge; the media rows are gone.
func TestRepo_Delete_CleanReturnsMediaKeys(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedCompetition(t, db, id, "deletable", "draft")
	seedMedia(t, db, id, 0)
	seedMedia(t, db, id, 1)

	keys, err := r.Delete(ctx, id)
	require.NoError(t, err)
	require.Len(t, keys, 2)

	// Competition + its media rows are gone.
	_, err = r.Get(ctx, id)
	require.ErrorIs(t, err, biz.ErrResourceNotFound)

	var mediaLeft int
	require.NoError(t, db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM media WHERE owner_id = $1`, id.String()).Scan(&mediaLeft))
	require.Equal(t, 0, mediaLeft)
}

// A competition with sold tickets cannot be deleted.
func TestRepo_Delete_BlockedByTickets(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedCompetition(t, db, id, "has-tickets", "live")
	seedTicket(t, db, id)

	_, err := r.Delete(ctx, id)
	require.ErrorIs(t, err, biz.ErrCompetitionHasEntrants)

	// Untouched.
	_, err = r.Get(ctx, id)
	require.NoError(t, err)
}

// A competition with an existing draw cannot be deleted.
func TestRepo_Delete_BlockedByDraw(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedCompetition(t, db, id, "has-draw", "closed")
	seedDrawRow(t, db, id)

	_, err := r.Delete(ctx, id)
	require.ErrorIs(t, err, biz.ErrCompetitionHasEntrants)
}
