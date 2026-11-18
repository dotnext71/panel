// Package dns provides DNS record conflict checks (including RFC 1034 CNAME
// exclusivity) and zone validation rules for dotNext.
package dns

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedRecordType = errors.New("dns: unsupported record type")
	ErrCNAMEConflict         = errors.New("dns: CNAME cannot coexist with other records at the same name")
)

var allowedTypes = map[string]struct{}{
	"A": {}, "AAAA": {}, "CNAME": {}, "MX": {}, "TXT": {},
	"NS": {}, "SRV": {}, "CAA": {}, "PTR": {}, "SOA": {},
}

// RecordSummary represents an existing record in a zone for conflict checking.
type RecordSummary struct {
	ID   string
	Name string
	Type string
}

// ValidateRecordType checks whether rtype is supported by dotNext/PowerDNS.
func ValidateRecordType(rtype string) (string, error) {
	up := strings.ToUpper(strings.TrimSpace(rtype))
	if _, ok := allowedTypes[up]; !ok {
		return "", ErrUnsupportedRecordType
	}
	return up, nil
}

// CheckCNAMEConflict enforces RFC 1034 §3.6.2: a CNAME at ownerName cannot
// coexist with any non-CNAME record (or another CNAME) at the same ownerName.
func CheckCNAMEConflict(ownerName, newType string, excludeID string, existing []RecordSummary) error {
	normName := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(ownerName)), ".")
	upType := strings.ToUpper(strings.TrimSpace(newType))

	for _, r := range existing {
		if excludeID != "" && r.ID == excludeID {
			continue
		}
		rName := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(r.Name)), ".")
		if rName != normName {
			continue
		}
		rType := strings.ToUpper(strings.TrimSpace(r.Type))
		if upType == "CNAME" || rType == "CNAME" {
			return ErrCNAMEConflict
		}
	}
	return nil
}
