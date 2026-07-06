package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"application/internal/datasource"
	"application/internal/user/biz"
	"application/internal/user/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ticket struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryTicket = (*ticket)(nil)

// NewTicket constructs the pgx-backed ticket repository.
func NewTicket(logger *slog.Logger, db *datasource.PostgresDB) *ticket {
	return &ticket{
		logger: logger.With("layer", "TicketRepo"),
		tracer: otel.Tracer("TicketRepo"),
		db:     db,
	}
}

// Purchase runs the buy flow atomically: read the competition's ticket price
// from the shared `competitions` table, insert `quantity` ticket rows, and bump
// the user's tickets_owned + total_spent. Design note: we read the price from
// the shared DB (consistent with the single-datasource template) and do NOT
// write competitions.tickets_sold here — that column is owned by the competition
// service; keeping it in sync belongs in an event/job (JetStream is available).
func (r *ticket) Purchase(ctx context.Context, in biz.PurchaseInput) (biz.PurchaseResult, error) {
	logger := r.logger.With("method", "Purchase")

	ctx, span := r.tracer.Start(ctx, "Purchase")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return biz.PurchaseResult{}, err
	}

	defer func() { _ = tx.Rollback() }()

	// Suspension guard lives INSIDE the purchase transaction so a concurrent
	// suspend cannot race a purchase past it.
	var isActive bool

	activeErr := tx.QueryRowContext(
		ctx,
		`SELECT is_active FROM users WHERE id = $1`, in.UserID,
	).Scan(&isActive)

	if errors.Is(activeErr, sql.ErrNoRows) {
		return biz.PurchaseResult{}, biz.ErrResourceNotFound
	}

	if activeErr != nil {
		logger.WarnContext(ctx, "failed to read user active flag", "error", activeErr)

		return biz.PurchaseResult{}, activeErr
	}

	if !isActive {
		return biz.PurchaseResult{}, biz.ErrUserSuspended
	}

	var pricePence int64

	priceErr := tx.QueryRowContext(
		ctx,
		`SELECT ticket_price_pence FROM competitions WHERE id = $1`, in.CompetitionID,
	).Scan(&pricePence)

	if errors.Is(priceErr, sql.ErrNoRows) {
		return biz.PurchaseResult{}, biz.ErrCompetitionNotFound
	}

	if priceErr != nil {
		logger.WarnContext(ctx, "failed to read competition price", "error", priceErr)

		return biz.PurchaseResult{}, priceErr
	}

	tickets, err := insertTickets(ctx, tx, in, time.Now().UTC())
	if err != nil {
		logger.WarnContext(ctx, "failed to insert tickets", "error", err)

		return biz.PurchaseResult{}, err
	}

	totalCost := pricePence * int64(in.Quantity)

	updatedUser, err := bumpUserTotals(ctx, tx, in.UserID, in.Quantity, totalCost)
	if err != nil {
		return biz.PurchaseResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return biz.PurchaseResult{}, err
	}

	return biz.PurchaseResult{Tickets: tickets, User: updatedUser, TotalCostPence: totalCost}, nil
}

func insertTickets(ctx context.Context, tx *sql.Tx, in biz.PurchaseInput, at time.Time) ([]entity.Ticket, error) {
	tickets := make([]entity.Ticket, 0, in.Quantity)

	for range in.Quantity {
		t := entity.Ticket{
			ID:            uuid.New(),
			CompetitionID: in.CompetitionID,
			UserID:        in.UserID,
			PurchasedAt:   at,
		}

		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO tickets (id, competition_id, user_id, purchased_at) VALUES ($1, $2, $3, $4)`,
			t.ID, t.CompetitionID, t.UserID, at,
		)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, t)
	}

	return tickets, nil
}

// bumpUserTotals atomically increments the user's counters (SET col = col + n is
// race-free without an explicit row lock), then reads the row back for the
// response.
func bumpUserTotals(
	ctx context.Context,
	tx *sql.Tx,
	userID uuid.UUID,
	quantity int,
	totalCost int64,
) (entity.User, error) {
	res, err := tx.ExecContext(
		ctx,
		`UPDATE users
			SET tickets_owned = tickets_owned + $1, total_spent_pence = total_spent_pence + $2
			WHERE id = $3`,
		quantity, totalCost, userID,
	)
	if err != nil {
		return entity.User{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.User{}, err
	}

	if affected == 0 {
		return entity.User{}, biz.ErrResourceNotFound
	}

	row := tx.QueryRowContext(
		ctx,
		`SELECT id, name, email, tickets_owned, total_spent_pence, is_active, created_at FROM users WHERE id = $1`,
		userID,
	)

	var u entity.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.TicketsOwned, &u.TotalSpentPence, &u.IsActive, &u.CreatedAt); err != nil {
		return entity.User{}, err
	}

	return u, nil
}

// ListByUser returns a user's tickets, newest first.
func (r *ticket) ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.Ticket, error) {
	logger := r.logger.With("method", "ListByUser")

	query := `SELECT id, competition_id, user_id, purchased_at FROM tickets
		WHERE user_id = $1 ORDER BY purchased_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		logger.WarnContext(ctx, "failed to query tickets", "error", err)

		return nil, err
	}
	defer rows.Close()

	tickets := []entity.Ticket{}

	for rows.Next() {
		var t entity.Ticket
		if err := rows.Scan(&t.ID, &t.CompetitionID, &t.UserID, &t.PurchasedAt); err != nil {
			logger.WarnContext(ctx, "failed to scan ticket", "error", err)

			continue
		}

		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}
