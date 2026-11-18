// Package httpserver defines consistent HTTP response envelopes, typed
// API errors, and pagination helpers for dotNext REST endpoints.
package httpserver

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIError is a structured, non-leaking HTTP error payload.
type APIError struct {
	Status    int    `json:"-"`
	ErrorCode string `json:"error"`
	Detail    string `json:"detail,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func (e APIError) Error() string {
	if e.Detail == "" {
		return e.ErrorCode
	}
	return e.ErrorCode + ": " + e.Detail
}

// NewError builds an APIError with explicit status, machine code, and user-safe detail.
func NewError(status int, code, detail string) APIError {
	return APIError{Status: status, ErrorCode: code, Detail: detail}
}

// RespondError writes a uniform JSON error envelope populated with the
// request's correlation ID when present.
func RespondError(c *gin.Context, err APIError) {
	if err.Status == 0 {
		err.Status = http.StatusInternalServerError
	}
	if err.ErrorCode == "" {
		err.ErrorCode = "internal_error"
	}
	if err.RequestID == "" {
		err.RequestID = c.GetString("request_id")
		if err.RequestID == "" {
			err.RequestID = c.GetHeader("X-Request-ID")
		}
	}
	c.AbortWithStatusJSON(err.Status, err)
}

// ListEnvelope matches the canonical paginated response contract expected by
// the UI and API clients (`{data, total, page, page_size}`).
type ListEnvelope[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// RespondList writes a paginated ListEnvelope, normalizing nil slices to `[]`.
func RespondList[T any](c *gin.Context, items []T, total int64, page, pageSize int) {
	if items == nil {
		items = []T{}
	}
	c.JSON(http.StatusOK, ListEnvelope[T]{
		Data:     items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Pagination holds validated query parameters for list endpoints.
type Pagination struct {
	Page     int
	PageSize int
	Offset   int
	Search   string
	Sort     string
	Order    string
}

// ParsePagination reads `page`, `page_size`, `q`, `sort`, and `order` with
// safe defaults and upper bounds.
func ParsePagination(c *gin.Context, defaultPageSize, maxPageSize int) Pagination {
	if defaultPageSize <= 0 {
		defaultPageSize = 50
	}
	if maxPageSize <= 0 {
		maxPageSize = 200
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if size < 1 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	order := strings.ToLower(strings.TrimSpace(c.DefaultQuery("order", "desc")))
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	return Pagination{
		Page:     page,
		PageSize: size,
		Offset:   (page - 1) * size,
		Search:   strings.TrimSpace(c.Query("q")),
		Sort:     strings.TrimSpace(c.Query("sort")),
		Order:    order,
	}
}
