package phases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dotnext-panel/dotnext/internal/agent"
	"github.com/dotnext-panel/dotnext/internal/autoresponderops"
	"github.com/dotnext-panel/dotnext/internal/models"
	"github.com/dotnext-panel/dotnext/internal/repository"
)

// DEAD CODE — NOT registered anywhere (only NewDisclaimerPhase is). Superseded
// by mailbox.sieve.apply / forwarderops.Converge (GH #1795): the autoresponder
// now compiles into the single active dotnext-managed standard SieveScript, NOT a
// separate native VacationResponse. Do NOT wire this phase — pushing
// "autoresponder.set" would recreate and activate a "vacation" script, which
// deactivates dotnext-managed and breaks external forwarding fleet-wide.
//
// autoresponderPhase converges dotnext email_autoresponders → Stalwart
// VacationResponse. DB is truth (ADR-0051); operator changes via Stalwart
// admin are overwritten on next tick.
type autoresponderPhase struct {
	agent          agent.AgentInterface
	autoresponders repository.EmailAutoresponderRepository
}

func NewAutoresponderPhase(ag agent.AgentInterface, ar repository.EmailAutoresponderRepository) Phase {
	return &autoresponderPhase{agent: ag, autoresponders: ar}
}

func (p *autoresponderPhase) Name() string { return "autoresponder" }

func (p *autoresponderPhase) ReconcileDomain(ctx context.Context, _ *models.Domain, _ map[string]interface{}) error {
	return nil // mailbox-scoped feature
}

func (p *autoresponderPhase) ReconcileMailbox(ctx context.Context, mb *models.Mailbox, dom *models.Domain, _ map[string]interface{}) error {
	if p.agent == nil || p.autoresponders == nil || dom == nil {
		return nil
	}
	ar, err := p.autoresponders.FindByMailboxID(ctx, mb.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Nothing to reconcile.
			return nil
		}
		return fmt.Errorf("autoresponder: find %s: %w", mb.ID, err)
	}
	// JAB-346: the agent payload comes from the one canonical projection the
	// HTTP + CLI Set path also uses, so all three callers push byte-identical
	// parameters.
	params := autoresponderops.AgentParams(mb.LocalPart+"@"+dom.Name, ar)
	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := p.agent.Call(callCtx, "autoresponder.set", params); err != nil {
		return fmt.Errorf("autoresponder: agent push for %s: %w", mb.ID, err)
	}
	return nil
}
