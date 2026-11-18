package cyberpanel

import (
	"github.com/dotnext-panel/dotnext/internal/migrate"
	"github.com/dotnext-panel/dotnext/internal/models"
)

func init() {
	migrate.Register(models.MigrationSourceCyberPanel, func() migrate.Discoverer {
		return New()
	})
}
