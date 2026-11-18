package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRespondErrorAndListEnvelope(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/err", func(c *gin.Context) {
		RespondError(c, NewError(http.StatusUnprocessableEntity, "invalid_domain", "domain is invalid"))
	})
	r.GET("/list", func(c *gin.Context) {
		p := ParsePagination(c, 25, 100)
		var empty []string
		RespondList(c, empty, 0, p.Page, p.PageSize)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
	var errBody map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("unmarshal err: %v", err)
	}
	if errBody["error"] != "invalid_domain" || errBody["detail"] != "domain is invalid" {
		t.Fatalf("unexpected error body: %v", errBody)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/list?page=2&page_size=500&order=ASC", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
	var listBody ListEnvelope[string]
	if err := json.Unmarshal(w2.Body.Bytes(), &listBody); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}
	if listBody.Data == nil || len(listBody.Data) != 0 || listBody.Page != 2 || listBody.PageSize != 100 {
		t.Fatalf("unexpected list envelope: %+v", listBody)
	}
}
