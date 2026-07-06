// Package audit records admin mutations into the shared admin_audit_log table
// (migration 000006). Every service that exposes admin write endpoints calls
// Recorder.Record after a successful mutation, attributing the action to the
// admin identified by the verified JWT claims in the request context.
package audit

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"application/pkg/middlewares"

	"github.com/google/uuid"
)

// DB is the minimal database seam (satisfied by *datasource.PostgresDB and
// *sql.DB/*sql.Tx) so the recorder stays decoupled from the datasource layer.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Entry is one audit record. ActorID/ActorEmail default from context claims.
type Entry struct {
	ActorID    string
	ActorEmail string
	Action     string // e.g. "draw.void", "competition.update"
	EntityType string // e.g. "draw", "competition", "user", "media"
	EntityID   string
	Reason     string // free-text reason/summary; required for sensitive actions by the caller
}

// Recorder writes audit entries. A nil Recorder is a no-op so tests and
// services without a DB can skip wiring it.
type Recorder struct {
	logger *slog.Logger
	db     DB
}

// NewRecorder constructs a Recorder.
func NewRecorder(logger *slog.Logger, db DB) *Recorder {
	return &Recorder{logger: logger.With("layer", "AuditRecorder"), db: db}
}

// Record inserts an audit row. Failures are logged, never returned: an audit
// insert must not roll back or mask the admin mutation it describes — losing
// one audit line is preferable to failing the action after it already
// happened. Callers that need atomicity pass their *sql.Tx as db instead.
func (r *Recorder) Record(ctx context.Context, e Entry) {
	if r == nil || r.db == nil {
		return
	}

	if claims, ok := middlewares.ClaimsFromContext(ctx); ok {
		if e.ActorID == "" {
			e.ActorID = claims.Subject
		}

		if e.ActorEmail == "" {
			e.ActorEmail = claims.Email
		}
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO admin_audit_log (id, actor_id, actor_email, action, entity_type, entity_id, reason, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		uuid.New(), e.ActorID, e.ActorEmail, e.Action, e.EntityType, e.EntityID, e.Reason, time.Now().UTC(),
	)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to write audit entry",
			"error", err, "action", e.Action, "entity_type", e.EntityType, "entity_id", e.EntityID)
	}
}
