package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/dotnext-panel/dotnext/internal/auth"
	"github.com/dotnext-panel/dotnext/internal/ginctx"
	"github.com/dotnext-panel/dotnext/internal/repository"
)

type stubUserRepo struct{ repository.UserRepository }

func TestMePhpCli_RejectsBadVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		ginctx.SetClaims(c, &auth.AccessClaims{UserID: "u1"})
		c.Next()
	})
	RegisterMePhpCliRoutes(v1, MePhpCliConfig{Users: &stubUserRepo{}})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/me/php-cli-version",
		strings.NewReader(`{"version":"9; rm -rf"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "invalid_version") {
		t.Fatalf("want 400 invalid_version, got %d %s", rec.Code, rec.Body.String())
	}
}
