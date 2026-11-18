package models

import "testing"

func TestIsFreeHostnameBase(t *testing.T) {
	for name, want := range map[string]bool{
		"203-0-113-7.dotnexthosted.com":   true,
		"203-0-113-7-b.dotnexthosted.com": true,
		"203-0-113-7.DotNextHosted.com.": true,
		"dotnexthosted.com":               false,
		".dotnexthosted.com":              false,
		"a.b.dotnexthosted.com":           false, // a preview host, not a base
		"evil-dotnexthosted.com":          false,
		"dotnexthosted.com.evil.com":      false,
		"preview.host.tld":               false,
		"":                               false,
	} {
		if got := IsFreeHostnameBase(name); got != want {
			t.Errorf("IsFreeHostnameBase(%q) = %v, want %v", name, got, want)
		}
	}
}
