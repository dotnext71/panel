package hestiacp

import (
	"github.com/dotnext-panel/dotnext/internal/migrate"
	"github.com/dotnext-panel/dotnext/internal/models"
)

func init() {
	migrate.Register(models.MigrationSourceHestia, func() migrate.Discoverer {
		return New()
	})
}
