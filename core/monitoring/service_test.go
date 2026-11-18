package system

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dotnext-panel/dotnext/internal/platform/jobs"
)

func TestRegisterWorkerStatusRoute(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	sup := jobs.NewSupervisor(nil)
	sup.Register("reconciler", time.Second, func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	sup.Start(ctx)
	defer func() {
		cancel()
		sup.Wait()
	}()

	r := gin.New()
	RegisterWorkerStatusRoute(r.Group("/api/v1/admin"), sup)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/workers", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var payload struct {
		Total int                `json:"total"`
		Data  []jobs.WorkerState `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Total != 1 || len(payload.Data) != 1 || payload.Data[0].Name != "reconciler" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}
