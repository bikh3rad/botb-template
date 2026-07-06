package repo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log/slog"
	"math/big"
	"strconv"
	"time"

	"application/internal/datasource"
	"application/internal/draw/biz"
	"application/internal/draw/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type draw struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.Repository = (*draw)(nil)

// NewDraw constructs the pgx-backed draw repository.
func NewDraw(logger *slog.Logger, db *datasource.PostgresDB) *draw {
	return &draw{
		logger: logger.With("layer", "DrawRepo"),
		tracer: otel.Tracer("DrawRepo"),
		db:     db,
	}
}

const drawColumns = `id, competition_id, winner_user_id, winner_ticket_id, prize,
	status, void_reason, drawn_at, created_at, updated_at`

// Create inserts a pending draw (id pre-generated) and returns the stored row.
func (r *draw) Create(ctx context.Context, d entity.Draw) (entity.Draw, error) {
	logger := r.logger.With("method", "Create")

	query := `INSERT INTO draws (id, competition_id, prize, status)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query, d.ID, d.CompetitionID, d.Prize, string(d.Status))
	if err := row.Scan(&d.CreatedAt, &d.UpdatedAt); err != nil {
		logger.WarnContext(ctx, "failed to insert draw", "error", err)

		return entity.Draw{}, err
	}

	return d, nil
}

// Get returns a draw by id, mapping a missing row to ErrResourceNotFound.
func (r *draw) Get(ctx context.Context, id uuid.UUID) (entity.Draw, error) {
	logger := r.logger.With("method", "Get")

	query := `SELECT ` + drawColumns + ` FROM draws WHERE id = $1`

	d, err := scanDraw(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Draw{}, biz.ErrResourceNotFound
		}

		logger.WarnContext(ctx, "failed to scan draw", "error", err)

		return entity.Draw{}, err
	}

	return d, nil
}

// List returns a page of draws (optionally filtered by a prize substring) plus
// the total match count.
func (r *draw) List(ctx context.Context, filter biz.ListFilter) (biz.DrawPage, error) {
	logger := r.logger.With("method", "List")

	where := ""
	args := []any{}

	if filter.Query != "" {
		where = ` WHERE prize ILIKE $1`

		args = append(args, "%"+filter.Query+"%")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM draws`+where, args...).Scan(&total); err != nil {
		logger.WarnContext(ctx, "failed to count draws", "error", err)

		return biz.DrawPage{}, err
	}

	// Limit/Offset are bounded ints (capped in the use case), so they are inlined
	// as literals rather than bound parameters — portable across drivers.
	query := `SELECT ` + drawColumns + ` FROM draws` + where +
		` ORDER BY created_at DESC LIMIT ` + strconv.Itoa(filter.Limit) +
		` OFFSET ` + strconv.Itoa(filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WarnContext(ctx, "failed to query draws", "error", err)

		return biz.DrawPage{}, err
	}
	defer rows.Close()

	draws := []entity.Draw{}

	for rows.Next() {
		d, err := scanDraw(rows)
		if err != nil {
			logger.WarnContext(ctx, "failed to scan draw row", "error", err)

			continue
		}

		draws = append(draws, d)
	}

	if err := rows.Err(); err != nil {
		return biz.DrawPage{}, err
	}

	return biz.DrawPage{Draws: draws, Total: total}, nil
}

// Run picks a winning ticket and marks the draw drawn, atomically. The status
// guard lives inside the transaction and the UPDATE is conditional on the draw
// still being pending, so two concurrent runs can never both succeed.
//
// Winner selection reads the competition's tickets from the shared `tickets`
// table (owned by the user service) — consistent with the single-datasource
// template. Design note: we deliberately do NOT mark the competition closed
// here; that field is owned by the competition service, and a real system would
// emit a JetStream event to sync it (mirroring the ticket-purchase note).
func (r *draw) Run(ctx context.Context, id uuid.UUID) (entity.Draw, error) {
	logger := r.logger.With("method", "Run")

	ctx, span := r.tracer.Start(ctx, "Run")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Draw{}, err
	}

	defer func() { _ = tx.Rollback() }()

	competitionID, err := loadPendingCompetition(ctx, tx, id)
	if err != nil {
		return entity.Draw{}, err
	}

	winnerTicketID, winnerUserID, err := pickWinner(ctx, tx, competitionID)
	if err != nil {
		return entity.Draw{}, err
	}

	drawnAt := time.Now().UTC()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE draws
			SET winner_user_id = $1, winner_ticket_id = $2, status = 'drawn',
			    drawn_at = $3, updated_at = $4
			WHERE id = $5 AND status = 'pending'`,
		winnerUserID, winnerTicketID, drawnAt, drawnAt, id,
	)
	if err != nil {
		logger.WarnContext(ctx, "failed to update draw", "error", err)

		return entity.Draw{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.Draw{}, err
	}

	// Lost the race to a concurrent run.
	if affected == 0 {
		return entity.Draw{}, biz.ErrAlreadyDrawn
	}

	updated, err := scanDraw(tx.QueryRowContext(ctx, `SELECT `+drawColumns+` FROM draws WHERE id = $1`, id))
	if err != nil {
		return entity.Draw{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.Draw{}, err
	}

	return updated, nil
}

// loadPendingCompetition reads the draw's competition id, guarding that the draw
// exists and is still pending (the guard lives inside the caller's transaction).
func loadPendingCompetition(ctx context.Context, tx *sql.Tx, id uuid.UUID) (uuid.UUID, error) {
	var (
		competitionID uuid.UUID
		status        string
	)

	err := tx.QueryRowContext(
		ctx,
		`SELECT competition_id, status FROM draws WHERE id = $1`, id,
	).Scan(&competitionID, &status)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, biz.ErrResourceNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	if entity.Status(status) != entity.StatusPending {
		return uuid.Nil, biz.ErrAlreadyDrawn
	}

	return competitionID, nil
}

// pickWinner reads the competition's tickets and returns a uniformly-random
// (crypto/rand) winning ticket + its owner.
func pickWinner(ctx context.Context, tx *sql.Tx, competitionID uuid.UUID) (ticketID, userID uuid.UUID, err error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, user_id FROM tickets WHERE competition_id = $1`, competitionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	defer rows.Close()

	type ref struct{ ticket, user uuid.UUID }

	var tickets []ref

	for rows.Next() {
		var t ref
		if err := rows.Scan(&t.ticket, &t.user); err != nil {
			return uuid.Nil, uuid.Nil, err
		}

		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	if len(tickets) == 0 {
		return uuid.Nil, uuid.Nil, biz.ErrNoTickets
	}

	idx, err := randomIndex(len(tickets))
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return tickets[idx].ticket, tickets[idx].user, nil
}

// randomIndex returns a uniformly-distributed index in [0, n) using crypto/rand.
func randomIndex(n int) (int, error) {
	if n <= 0 {
		return 0, biz.ErrNoTickets
	}

	idx, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}

	return int(idx.Int64()), nil
}

// UpdatePrize edits the prize text (UPDATE-then-SELECT, no RETURNING, for
// ramsql compatibility).
func (r *draw) UpdatePrize(ctx context.Context, id uuid.UUID, prize string) (entity.Draw, error) {
	logger := r.logger.With("method", "UpdatePrize")

	res, err := r.db.ExecContext(ctx,
		`UPDATE draws SET prize = $1, updated_at = $2 WHERE id = $3`,
		prize, time.Now().UTC(), id)
	if err != nil {
		logger.WarnContext(ctx, "failed to update draw prize", "error", err)

		return entity.Draw{}, err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return entity.Draw{}, err
	} else if affected == 0 {
		return entity.Draw{}, biz.ErrResourceNotFound
	}

	return r.Get(ctx, id)
}

// Void marks a pending or drawn draw void with the reason. Optimistic CAS:
// the UPDATE is conditional on the status we just read, so a concurrent
// void/run makes it a no-op and we report the conflict instead of clobbering.
func (r *draw) Void(ctx context.Context, id uuid.UUID, reason string) (entity.Draw, error) {
	logger := r.logger.With("method", "Void")

	current, err := r.Get(ctx, id)
	if err != nil {
		return entity.Draw{}, err
	}

	if current.Status == entity.StatusVoid {
		return entity.Draw{}, biz.ErrInvalidState
	}

	res, err := r.db.ExecContext(ctx,
		`UPDATE draws SET status = $1, void_reason = $2, updated_at = $3
		 WHERE id = $4 AND status = $5`,
		string(entity.StatusVoid), reason, time.Now().UTC(), id, string(current.Status))
	if err != nil {
		logger.WarnContext(ctx, "failed to void draw", "error", err)

		return entity.Draw{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.Draw{}, err
	}

	if affected == 0 {
		// Lost a race to a concurrent void/run.
		return entity.Draw{}, biz.ErrInvalidState
	}

	return r.Get(ctx, id)
}

// Reassign swaps a DRAWN draw's winner to another ticket of the same
// competition, validating everything inside one transaction.
func (r *draw) Reassign(ctx context.Context, id uuid.UUID, ticketID uuid.UUID) (entity.Draw, error) {
	logger := r.logger.With("method", "Reassign")

	ctx, span := r.tracer.Start(ctx, "Reassign")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Draw{}, err
	}

	defer func() { _ = tx.Rollback() }()

	var (
		competitionID uuid.UUID
		status        string
	)

	err = tx.QueryRowContext(
		ctx,
		`SELECT competition_id, status FROM draws WHERE id = $1`, id,
	).Scan(&competitionID, &status)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Draw{}, biz.ErrResourceNotFound
	}

	if err != nil {
		return entity.Draw{}, err
	}

	if entity.Status(status) != entity.StatusDrawn {
		return entity.Draw{}, biz.ErrInvalidState
	}

	var (
		ticketCompetition uuid.UUID
		ticketOwner       uuid.UUID
	)

	err = tx.QueryRowContext(
		ctx,
		`SELECT competition_id, user_id FROM tickets WHERE id = $1`, ticketID,
	).Scan(&ticketCompetition, &ticketOwner)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Draw{}, biz.ErrTicketMismatch
	}

	if err != nil {
		return entity.Draw{}, err
	}

	if ticketCompetition != competitionID {
		return entity.Draw{}, biz.ErrTicketMismatch
	}

	if _, err := tx.ExecContext(ctx,
		`UPDATE draws SET winner_ticket_id = $1, winner_user_id = $2, updated_at = $3 WHERE id = $4`,
		ticketID, ticketOwner, time.Now().UTC(), id); err != nil {
		logger.WarnContext(ctx, "failed to reassign winner", "error", err)

		return entity.Draw{}, err
	}

	updated, err := scanDraw(tx.QueryRowContext(ctx, `SELECT `+drawColumns+` FROM draws WHERE id = $1`, id))
	if err != nil {
		return entity.Draw{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.Draw{}, err
	}

	return updated, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanDraw(s scanner) (entity.Draw, error) {
	var (
		d            entity.Draw
		status       string
		winnerUser   uuid.NullUUID
		winnerTicket uuid.NullUUID
		drawnAt      sql.NullTime
	)

	if err := s.Scan(
		&d.ID, &d.CompetitionID, &winnerUser, &winnerTicket, &d.Prize,
		&status, &d.VoidReason, &drawnAt, &d.CreatedAt, &d.UpdatedAt,
	); err != nil {
		return entity.Draw{}, err
	}

	d.Status = entity.Status(status)

	if winnerUser.Valid {
		u := winnerUser.UUID
		d.WinnerUserID = &u
	}

	if winnerTicket.Valid {
		tk := winnerTicket.UUID
		d.WinnerTicketID = &tk
	}

	if drawnAt.Valid {
		at := drawnAt.Time
		d.DrawnAt = &at
	}

	return d, nil
}
