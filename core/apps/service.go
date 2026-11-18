// Package apps provides dotNext boundary validation for managed web
// applications (WordPress, Docker Compose apps, and Python WSGI/ASGI apps).
package apps

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

var (
	ErrInvalidSlug    = errors.New("apps: slug must match ^[a-z0-9][a-z0-9-]{1,47}$")
	ErrInvalidSubpath = errors.New("apps: mount path must start with / without traversal")
	ErrInvalidMemory  = errors.New("apps: memory_mb must be between 64 and 65536")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,47}$`)

// AppDeploySpec holds normalized deployment attributes checked before runtime provisioning.
type AppDeploySpec struct {
	Slug     string
	MountURL string
	MemoryMB int
}

// ValidateDeploySpec validates and normalizes application deployment parameters.
func ValidateDeploySpec(slug, mountURL string, memoryMB int) (AppDeploySpec, error) {
	normSlug := strings.ToLower(strings.TrimSpace(slug))
	if !slugPattern.MatchString(normSlug) {
		return AppDeploySpec{}, ErrInvalidSlug
	}
	normMount, err := NormalizeMountURL(mountURL)
	if err != nil {
		return AppDeploySpec{}, err
	}
	if memoryMB != 0 && (memoryMB < 64 || memoryMB > 65536) {
		return AppDeploySpec{}, ErrInvalidMemory
	}
	return AppDeploySpec{
		Slug:     normSlug,
		MountURL: normMount,
		MemoryMB: memoryMB,
	}, nil
}

// NormalizeMountURL validates an HTTP subpath like "/" or "/blog".
func NormalizeMountURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "/", nil
	}
	if strings.ContainsAny(trimmed, " \t\r\n?#\\") || strings.Contains(trimmed, "..") {
		return "", ErrInvalidSubpath
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	cleaned := path.Clean(trimmed)
	if !strings.HasPrefix(cleaned, "/") {
		return "", ErrInvalidSubpath
	}
	return cleaned, nil
}
