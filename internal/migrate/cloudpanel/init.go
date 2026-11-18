package cloudpanel

import (
	"github.com/dotnext-panel/dotnext/internal/migrate"
	"github.com/dotnext-panel/dotnext/internal/models"
)

func init() {
	migrate.Register(models.MigrationSourceCloudPanel, func() migrate.Discoverer {
		return New()
	})
}
