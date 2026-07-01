package repo_test

import (
	"application/internal/datasource"
	"application/internal/draw/biz"
	"application/internal/draw/entity"
	"application/internal/draw/repo"
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

// newRamsqlDB spins up an isolated ramsql DB with the draws + tickets tables.
// ramsql cannot parse the production transactional SQL used by the Run
// happy-path (parameterized LIMIT/OFFSET is inlined here, but the winner UPDATE
// / crypto selection is only fully integration-tested against real Postgres and
// covered by the biz mock tests). ramsql is used here for the parseable read
// paths and the pre-UPDATE failure branches of Run.
func newRamsqlDB(t *testing.T) *datasource.PostgresDB {
	t.Helper()

	db, err := sql.Open("ramsql", "draw_"+uuid.NewString())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()

	stmts := []string{
		`CREATE TABLE draws (
			id TEXT PRIMARY KEY, competition_id TEXT,
			winner_user_id TEXT, winner_ticket_id TEXT, prize TEXT, status TEXT,
			drawn_at TIMESTAMP, created_at TIMESTAMP, updated_at TIMESTAMP)`,
		`CREATE TABLE tickets (id TEXT PRIMARY KEY, competition_id TEXT, user_id TEXT)`,
	}
	for _, s := range stmts {
		_, err := db.ExecContext(ctx, s)
		require.NoError(t, err)
	}

	return &datasource.PostgresDB{DB: db}
}

func seedDraw(t *testing.T, db *datasource.PostgresDB, id, competitionID uuid.UUID, prize, status string) {
	t.Helper()

	now := time.Now().UTC()

	// winner_* and drawn_at are passed as nil bound params (→ SQL NULL) rather
	// than NULL literals, so scanDraw exercises the nullable columns via Get.
	_, err := db.ExecContext(context.Background(),
		`INSERT INTO draws
			(id, competition_id, winner_user_id, winner_ticket_id, prize, status, drawn_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		id.String(), competitionID.String(), nil, nil, prize, status, nil, now, now,
	)
	require.NoError(t, err)
}

func newRepo(db *datasource.PostgresDB) biz.Repository {
	return repo.NewDraw(slog.New(slog.NewTextHandler(io.Discard, nil)), db)
}

func TestRepo_Get(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedDraw(t, db, id, uuid.New(), "Audi RS3", "pending")

	got, err := r.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
	require.Equal(t, entity.StatusPending, got.Status)
	require.Nil(t, got.WinnerUserID)
	require.Nil(t, got.DrawnAt)
}

func TestRepo_Get_NotFound(t *testing.T) {
	db := newRamsqlDB(t)
	r := newRepo(db)

	_, err := r.Get(context.Background(), uuid.New())
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

func TestRepo_List(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	for i := 0; i < 3; i++ {
		seedDraw(t, db, uuid.New(), uuid.New(), "Prize", "pending")
	}

	page, err := r.List(ctx, biz.ListFilter{Limit: 2, Offset: 0})
	require.NoError(t, err)
	require.Equal(t, 3, page.Total)
	require.Len(t, page.Draws, 2)
}

func TestRepo_Run_AlreadyDrawn(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedDraw(t, db, id, uuid.New(), "Prize", "drawn")

	_, err := r.Run(ctx, id)
	require.ErrorIs(t, err, biz.ErrAlreadyDrawn)
}

func TestRepo_Run_NoTickets(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := newRepo(db)

	id := uuid.New()
	seedDraw(t, db, id, uuid.New(), "Prize", "pending")

	_, err := r.Run(ctx, id)
	require.ErrorIs(t, err, biz.ErrNoTickets)
}
