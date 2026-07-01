package repo_test

import (
	"application/internal/datasource"
	"application/internal/user/biz"
	"application/internal/user/repo"
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

func newRamsqlDB(t *testing.T) *datasource.PostgresDB {
	t.Helper()

	db, err := sql.Open("ramsql", "user_"+uuid.NewString())
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()

	stmts := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY, name TEXT, email TEXT,
			tickets_owned BIGINT, total_spent_pence BIGINT, created_at TIMESTAMP)`,
		`CREATE TABLE tickets (
			id TEXT PRIMARY KEY, competition_id TEXT, user_id TEXT, purchased_at TIMESTAMP)`,
		`CREATE TABLE competitions (id TEXT PRIMARY KEY, ticket_price_pence BIGINT)`,
	}
	for _, s := range stmts {
		_, err := db.ExecContext(ctx, s)
		require.NoError(t, err)
	}

	return &datasource.PostgresDB{DB: db}
}

func seedUser(t *testing.T, db *datasource.PostgresDB, id uuid.UUID, name, email string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO users (id, name, email, tickets_owned, total_spent_pence, created_at)
			VALUES ($1, $2, $3, 0, 0, $4)`,
		id.String(), name, email, time.Now().UTC(),
	)
	require.NoError(t, err)
}

func userRepo(db *datasource.PostgresDB) biz.RepositoryUser {
	return repo.NewUser(slog.New(slog.NewTextHandler(io.Discard, nil)), db)
}

func ticketRepo(db *datasource.PostgresDB) biz.RepositoryTicket {
	return repo.NewTicket(slog.New(slog.NewTextHandler(io.Discard, nil)), db)
}

func TestUser_Get(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := userRepo(db)

	id := uuid.New()
	seedUser(t, db, id, "Olivia", "olivia@example.com")

	got, err := r.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "olivia@example.com", got.Email)
}

func TestUser_Get_NotFound(t *testing.T) {
	db := newRamsqlDB(t)
	r := userRepo(db)

	_, err := r.Get(context.Background(), uuid.New())
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

func TestUser_List(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	r := userRepo(db)

	for i := 0; i < 3; i++ {
		seedUser(t, db, uuid.New(), "User", uuid.NewString()+"@example.com")
	}

	page, err := r.List(ctx, biz.UserListFilter{Limit: 2, Offset: 0})
	require.NoError(t, err)
	require.Equal(t, 3, page.Total)
	require.Len(t, page.Users, 2) // limited to page size
}

// NOTE: the happy-path Purchase transaction uses `SET col = col + $n` (an atomic,
// race-free increment) which ramsql's lexer cannot parse. That path is exercised
// through the biz mock tests (orchestration) and is only fully integration-tested
// against real Postgres. Here we cover the DB-touching failure path that ramsql
// *can* run: transaction begin + competition price lookup + rollback.
func TestTicket_Purchase_CompetitionNotFound(t *testing.T) {
	ctx := context.Background()
	db := newRamsqlDB(t)
	tr := ticketRepo(db)

	userID := uuid.New()
	seedUser(t, db, userID, "Buyer", "buyer@example.com")

	_, err := tr.Purchase(ctx, biz.PurchaseInput{CompetitionID: uuid.New(), UserID: userID, Quantity: 1})
	require.ErrorIs(t, err, biz.ErrCompetitionNotFound)
}
