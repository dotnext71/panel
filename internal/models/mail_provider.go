package models

// Mail provider enum (GH#181 / ADR-0120). MailProvider on a Domain is the
// single source of truth for the domain's mail posture; EmailEnabled and
// SkipAutoSAN are derived from it via DeriveMailFlags so the two booleans
// can never drift from operator intent.
const (
	MailProviderDotNext = "dotnext" // DotNext is the MTA (default; today's behavior)
	MailProviderNone    = "none"    // no mail anywhere — no mail records, no mail SANs
	MailProviderM365    = "m365"    // Microsoft 365
	MailProviderGoogle  = "google"  // Google Workspace
	// MailProviderCustom (GH #1627) is the posture of a domain created from a
	// custom DNS template: mail (and other services) live wherever the template's
	// records point. Like m365/google it is external — DotNext is not the MTA and
	// no DotNext mail SANs go on the cert. Unlike them the reconciler asserts
	// NOTHING for it: reconcileMailProviderRecords hands the zone off entirely so
	// the template's own apex MX/SPF/DKIM survive as tenant-owned records. It is
	// NEVER caller-supplied — createDomainOp sets it only when a dns_template_id
	// is chosen, and rejects a request that names it directly.
	MailProviderCustom = "custom"
)

// ValidMailProvider reports whether p is a recognised provider value.
// MailProviderCustom is recognised (so a persisted 'custom' row validates
// everywhere), but createDomainOp rejects it as CALLER input — 'custom' is only
// reachable by selecting a dns_template_id (GH #1627).
func ValidMailProvider(p string) bool {
	switch p {
	case MailProviderDotNext, MailProviderNone, MailProviderM365, MailProviderGoogle, MailProviderCustom:
		return true
	default:
		return false
	}
}

// MailProviderIsExternal reports whether the provider hosts mail somewhere
// other than DotNext (M365 / Google) — i.e. DotNext is not the MTA but the
// domain still has mail (so provider DNS records are published).
func MailProviderIsExternal(p string) bool {
	return p == MailProviderM365 || p == MailProviderGoogle || p == MailProviderCustom
}

// DeriveMailFlags returns the (EmailEnabled, SkipAutoSAN) pair implied by a
// provider. This is the ONLY place those two booleans are computed from
// provider intent — callers must use it on create/edit so the columns the
// reconciler + cert paths read stay consistent with the enum.
//
//	dotnext                 -> (true,  false)  DotNext MTA + dotnext mail SANs on cert
//	none/m365/google/custom-> (false, true)   not the DotNext MTA; no dotnext mail SANs
func DeriveMailFlags(provider string) (emailEnabled, skipAutoSAN bool) {
	if provider == MailProviderDotNext {
		return true, false
	}
	return false, true
}

// MailHostedElsewhere reports whether this server must NOT act as the mail
// server for a domain: email is off, or the mail provider is anything but
// DotNext (an empty provider is a row from before the enum and means DotNext).
//
// The mail server keeps such a domain in its registry while any account on it
// is left (every domain has a send-only noreply@ for its website's mail), and
// a registered domain counts as local. The reconciler lists these domains to
// the agent (mail.elsewhere.sync), which keeps them in Stalwart's in-memory
// lookup "dotnext-mail-elsewhere": the installer's route sends mail for them to
// their MX, and unauthenticated mail for them is refused.
func MailHostedElsewhere(emailEnabled bool, provider string) bool {
	if !emailEnabled {
		return true
	}
	return provider != "" && provider != MailProviderDotNext
}
