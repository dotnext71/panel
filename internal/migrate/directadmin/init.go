package directadmin

import (
	"github.com/dotnext-panel/dotnext/internal/migrate"
	"github.com/dotnext-panel/dotnext/internal/models"
)

func init() {
	migrate.Register(models.MigrationSourceDirectAdmin, func() migrate.Discoverer {
		return New()
	})
}
