// Package system provides dotNext system runtime endpoints including
// supervised background worker telemetry.
package system

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dotnext-panel/dotnext/internal/platform/jobs"
)

// RegisterWorkerStatusRoute mounts GET /admin/system/workers on the supplied
// admin-gated RouterGroup so operators can inspect background worker health.
func RegisterWorkerStatusRoute(adminGroup *gin.RouterGroup, sup *jobs.Supervisor) {
	if adminGroup == nil || sup == nil {
		return
	}
	adminGroup.GET("/system/workers", func(c *gin.Context) {
		workers := sup.Snapshot()
		c.JSON(http.StatusOK, gin.H{
			"data":  workers,
			"total": len(workers),
		})
	})
}
