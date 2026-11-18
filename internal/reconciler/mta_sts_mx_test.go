package reconciler

import (
	"testing"

	"github.com/dotnext-panel/dotnext/internal/models"
)

// TestMtaStsMXForDomain pins the policy mx: host to the domain's real MX
// (GH #441): DotNext-managed mail publishes MX mail.<domain>, so the MTA-STS
// policy must use that — never the panel hostname. External providers are not
// owned by DotNext, so MTA-STS must not be applied for them.
func TestMtaStsMXForDomain(t *testing.T) {
	cases := []struct {
		provider string
		wantMX   string
		wantOK   bool
	}{
		{"dotnext", "mail.example.com", true},
		{"", "mail.example.com", true}, // default is dotnext
		{"m365", "", false},
		{"google", "", false},
		{"none", "", false},
	}
	for _, c := range cases {
		d := &models.Domain{Name: "example.com", MailProvider: c.provider}
		mx, ok := mtaStsMXForDomain(d)
		if ok != c.wantOK || mx != c.wantMX {
			t.Errorf("provider %q: got (%q,%v) want (%q,%v)", c.provider, mx, ok, c.wantMX, c.wantOK)
		}
	}
}
