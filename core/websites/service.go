// Package websites provides domain normalization, validation, and document-root
// derivation for dotNext web hosting domains.
package websites

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrInvalidDomain = errors.New("websites: invalid domain name")

var domainLabelPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// NormalizeDomain trims whitespace, strips any trailing DNS dot, and lowercases.
func NormalizeDomain(domain string) string {
	d := strings.TrimSpace(strings.ToLower(domain))
	return strings.TrimSuffix(d, ".")
}

// ValidateDomain verifies RFC-safe FQDN syntax without shell metacharacters.
func ValidateDomain(domain string) error {
	d := NormalizeDomain(domain)
	if len(d) < 3 || len(d) > 253 || !strings.Contains(d, ".") {
		return ErrInvalidDomain
	}
	labels := strings.Split(d, ".")
	for _, lbl := range labels {
		if !domainLabelPattern.MatchString(lbl) {
			return ErrInvalidDomain
		}
	}
	return nil
}

// DefaultDocRoot builds the canonical web document root under a tenant's home.
func DefaultDocRoot(username, domain string) (string, error) {
	if err := ValidateDomain(domain); err != nil {
		return "", err
	}
	u := strings.TrimSpace(username)
	if u == "" || strings.ContainsAny(u, "/\\.") {
		return "", errors.New("websites: invalid username for docroot")
	}
	return filepath.Join("/home", u, "web", NormalizeDomain(domain), "public_html"), nil
}
