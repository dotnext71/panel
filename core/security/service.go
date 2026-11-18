// Package security provides dotNext boundary validation for network ACLs,
// CIDR ranges, and SSH public key types before privileged host enforcement.
package security

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"
)

var (
	ErrInvalidCIDROrIP = errors.New("security: invalid IP address or CIDR prefix")
	ErrInvalidSSHKey   = errors.New("security: unsupported or malformed SSH public key")
)

var allowedSSHKeyTypes = map[string]struct{}{
	"ssh-ed25519":                        {},
	"sk-ssh-ed25519@openssh.com":         {},
	"ecdsa-sha2-nistp256":                {},
	"ecdsa-sha2-nistp384":                {},
	"ecdsa-sha2-nistp521":                {},
	"sk-ecdsa-sha2-nistp256@openssh.com": {},
	"ssh-rsa":                            {},
}

// ParseIPOrCIDR normalizes either a single IP (into /32 or /128) or an explicit CIDR prefix.
func ParseIPOrCIDR(raw string) (netip.Prefix, error) {
	s := strings.TrimSpace(raw)
	if strings.Contains(s, "/") {
		p, err := netip.ParsePrefix(s)
		if err != nil {
			return netip.Prefix{}, fmt.Errorf("%w: %v", ErrInvalidCIDROrIP, err)
		}
		return p.Masked(), nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("%w: %v", ErrInvalidCIDROrIP, err)
	}
	return netip.PrefixFrom(addr, addr.BitLen()), nil
}

// ValidateSSHPublicKey verifies single-line AuthorizedKeys format and allowed key types.
func ValidateSSHPublicKey(raw string) (string, error) {
	line := strings.TrimSpace(raw)
	if line == "" || strings.ContainsAny(line, "\r\n\x00") {
		return "", ErrInvalidSSHKey
	}
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", ErrInvalidSSHKey
	}
	if _, ok := allowedSSHKeyTypes[fields[0]]; !ok {
		return "", fmt.Errorf("%w: algorithm %q not allowed", ErrInvalidSSHKey, fields[0])
	}
	return line, nil
}
