package biz

import (
	"context"
	"log/slog"

	"application/internal/competition/entity"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultAuditLimit = 20
	maxAuditLimit     = 100
)

// UsecaseAudit reads the shared admin audit log (recent admin actions across
// ALL services — they all write to the one admin_audit_log table). Read-only.
type UsecaseAudit interface {
	Recent(ctx context.Context, limit int) ([]entity.AuditEntry, error)
}

// RepositoryAudit reads audit rows.
type RepositoryAudit interface {
	Recent(ctx context.Context, limit int) ([]entity.AuditEntry, error)
}

type auditReader struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   RepositoryAudit
}

var _ UsecaseAudit = (*auditReader)(nil)

// NewAuditReader constructs the audit-log read use case.
func NewAuditReader(logger *slog.Logger, repo RepositoryAudit) *auditReader {
	return &auditReader{
		logger: logger.With("layer", "AuditReader"),
		tracer: otel.Tracer("AuditReaderUseCase"),
		repo:   repo,
	}
}

func (uc *auditReader) Recent(ctx context.Context, limit int) ([]entity.AuditEntry, error) {
	if limit <= 0 {
		limit = defaultAuditLimit
	}

	if limit > maxAuditLimit {
		limit = maxAuditLimit
	}

	return uc.repo.Recent(ctx, limit)
}
