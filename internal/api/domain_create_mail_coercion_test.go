package api

import (
	"context"
	"testing"

	"github.com/dotnext-panel/dotnext/internal/models"
)

// TestCreateDomainOp_MailModuleCoercion pins that createDomainOp (the REST/GUI/
// automation entrypoint) routes the resolved mail provider through the shared
// domainops.MailProviderForServer leaf (JAB-279 / GH #1409): a DotNext provider is
// coerced to none — email disabled — when the server mail module is off, and kept
// when it is on. This is the api-side proof that REST still routes through the
// leaf after the coercion moved out of this package; the CLI runs the identical
// leaf. Falsify: pass a literal true at the call site → the module-off case
// reddens (mail stays dotnext, email enabled).
func TestCreateDomainOp_MailModuleCoercion(t *testing.T) {
	uname := "alice"
	owner := &models.User{ID: "u-alice", Email: "alice@example.com", Username: &uname}

	run := func(mailModuleEnabled bool) *models.Domain {
		dom := newDCDomains()
		h := &domainHandler{cfg: DomainHandlerConfig{
			Users:          newAbUsers(owner),
			Domains:        dom,
			ServerSettings: &fakeStatusSettingsRepo{s: &models.ServerSettings{MailEnabled: mailModuleEnabled}},
		}}
		d, oerr := createDomainOp(context.Background(), h, createDomainInput{
			OwnerID:       owner.ID,
			Name:          "shop.example.com",
			MailProvider:  models.MailProviderDotNext,
			SSLMode:       models.SSLModeLE,
			SkipInlineSSL: true,
		})
		if oerr != nil {
			t.Fatalf("unexpected error (mailModuleEnabled=%v): %v", mailModuleEnabled, oerr)
		}
		return d
	}

	t.Run("module off coerces dotnext to none, email disabled", func(t *testing.T) {
		d := run(false)
		if d.MailProvider != models.MailProviderNone {
			t.Errorf("mail-less server must coerce dotnext to none, got %q", d.MailProvider)
		}
		if d.EmailEnabled {
			t.Error("a coerced domain must not have email enabled")
		}
	})

	t.Run("module on keeps dotnext, email enabled", func(t *testing.T) {
		d := run(true)
		if d.MailProvider != models.MailProviderDotNext {
			t.Errorf("mail-enabled server must keep dotnext, got %q", d.MailProvider)
		}
		if !d.EmailEnabled {
			t.Error("an uncoerced dotnext domain must have email enabled")
		}
	})
}
