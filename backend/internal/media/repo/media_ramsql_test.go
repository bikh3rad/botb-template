package repo_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"

	"application/internal/datasource"
	"application/internal/media/biz"
	"application/internal/media/entity"
	"application/internal/media/repo"

	"github.com/google/uuid"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

// newRamsqlDB spins up an isolated in-memory ramsql database with the media
// table, wrapped in the datasource.PostgresDB the repo expects. ramsql does not
// implement Postgres-specific features (RETURNING, DEFAULT now(), uuid type), so
// the schema uses portable types and the read paths (Get/ListByOwner) are what
// we exercise here; Create's RETURNING clause is covered by the biz mocks.
func newRamsqlDB(t *testing.T) *datasource.PostgresDB {
	t.Helper()

	db, err := sql.Open("ramsql", "media_"+uuid.NewString())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.ExecContext(context.Background(), `CREATE TABLE media (
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

	return &datasource.PostgresDB{DB: db}
}

func seedMedia(t *testing.T, db *datasource.PostgresDB, m entity.Media) {
	t.Helper()

	_, err := db.ExecContext(
		context.Background(), `INSERT INTO media
		(id, owner_type, owner_id, kind, bucket, object_key, content_type,
		 size_bytes, width, height, duration_seconds, position, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12, NOW())`,
		m.ID.String(), m.OwnerType, m.OwnerID.String(), string(m.Kind), m.Bucket,
		m.ObjectKey, m.ContentType, m.SizeBytes, m.Width, m.Height, nil, m.Position,
	)
	require.NoError(t, err)
}

func TestRepo_Get(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := repo.NewMedia(slog.New(slog.NewTextHandler(io.Discard, nil)), db)

	id := uuid.New()
	ownerID := uuid.New()
	seedMedia(t, db, entity.Media{
		ID: id, OwnerType: "competition", OwnerID: ownerID, Kind: entity.KindImage,
		Bucket: "botb-media", ObjectKey: "competition/x/y.png", ContentType: "image/png",
		SizeBytes: 2048, Position: 0,
	})

	got, err := r.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
	require.Equal(t, "competition", got.OwnerType)
	require.Equal(t, entity.KindImage, got.Kind)
	require.Equal(t, int64(2048), got.SizeBytes)
}

func TestRepo_Get_NotFound(t *testing.T) {
	db := newRamsqlDB(t)
	r := repo.NewMedia(slog.New(slog.NewTextHandler(io.Discard, nil)), db)

	_, err := r.Get(context.Background(), uuid.New())
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

func TestRepo_ListByOwner(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := repo.NewMedia(slog.New(slog.NewTextHandler(io.Discard, nil)), db)

	ownerID := uuid.New()
	for i := 0; i < 3; i++ {
		seedMedia(t, db, entity.Media{
			ID: uuid.New(), OwnerType: "competition", OwnerID: ownerID, Kind: entity.KindImage,
			Bucket: "botb-media", ObjectKey: "k", ContentType: "image/png", SizeBytes: 1, Position: i,
		})
	}
	// A different owner's media must not leak into the result.
	seedMedia(t, db, entity.Media{
		ID: uuid.New(), OwnerType: "competition", OwnerID: uuid.New(), Kind: entity.KindImage,
		Bucket: "botb-media", ObjectKey: "other", ContentType: "image/png", SizeBytes: 1, Position: 0,
	})

	got, err := r.ListByOwner(ctx, "competition", ownerID)
	require.NoError(t, err)
	require.Len(t, got, 3)
}
