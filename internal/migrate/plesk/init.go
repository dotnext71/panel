package plesk

import (
	"github.com/dotnext-panel/dotnext/internal/migrate"
	"github.com/dotnext-panel/dotnext/internal/models"
)

func init() {
	migrate.Register(models.MigrationSourcePlesk, func() migrate.Discoverer {
		return New()
	})
}
