package repo

import (
	"application/internal/competition/biz"

	"github.com/google/wire"
)

// ProviderSet wires the competition, category and content repositories to
// their biz interfaces.
var ProviderSet = wire.NewSet(
	NewCompetition,
	wire.Bind(new(biz.Repository), new(*competition)),

	NewCategory,
	wire.Bind(new(biz.RepositoryCategory), new(*category)),

	NewContent,
	wire.Bind(new(biz.RepositoryContent), new(*content)),

	NewAuditRepo,
	wire.Bind(new(biz.RepositoryAudit), new(*auditRepo)),
)
