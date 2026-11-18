// Package files provides dotNext filesystem boundary constructors built on
// internal/filesafe so file-manager and migration callers share consistent
// symlink-safe docroot and admin denylist confinement.
package files

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/dotnext-panel/dotnext/internal/filesafe"
)

var (
	ErrUnsafePath     = errors.New("files: path contains forbidden traversal or control characters")
	ErrEmptyMutable   = errors.New("files: admin write scope requires at least one mutable root")
)

// OpenTenantScope constructs a docroot-confined Scope for a hosting account.
func OpenTenantScope(userID, username string, docroots []string) (*filesafe.Scope, error) {
	return filesafe.NewScope(userID, username, docroots)
}

// OpenAdminScope constructs a system-wide admin Scope guarded by dotNext's
// DefaultAdminDeniedPrefixes (/etc/dotnext, /run/dotnext, /var/lib/dotnext, etc.).
func OpenAdminScope(userID, username string, mutableRoots []string) (*filesafe.Scope, error) {
	if len(mutableRoots) == 0 {
		return nil, ErrEmptyMutable
	}
	return filesafe.NewAdminScope(userID, username, filesafe.DefaultAdminDeniedPrefixes, mutableRoots)
}

// ValidateRelativePath normalizes a user-supplied relative path under a
// docroot and rejects NUL bytes, leading slashes, or parent traversal ("..").
func ValidateRelativePath(rel string) (string, error) {
	if strings.ContainsRune(rel, '\x00') || strings.ContainsAny(rel, "\r\n") {
		return "", ErrUnsafePath
	}
	trimmed := strings.TrimSpace(rel)
	if trimmed == "" || trimmed == "." {
		return ".", nil
	}
	if strings.HasPrefix(trimmed, "/") || strings.HasPrefix(trimmed, `\`) {
		return "", ErrUnsafePath
	}
	cleaned := filepath.ToSlash(filepath.Clean(trimmed))
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", ErrUnsafePath
	}
	return cleaned, nil
}
