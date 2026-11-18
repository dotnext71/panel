# Domains & Subdomains

Every hostname dotNext touches, what it is for, where it is configured in the
code, and what breaks when it is missing.

Two rules of thumb:

1. **Only the "Panel-owned" section needs you to buy a domain.** Everything
   under "Automatic (panel-generated)" is created by the panel per tenant
   domain and needs no purchase.
2. A missing panel-owned hostname **does not stop the panel from running** —
   each feature that uses one is either opt-in or already has a documented
   fallback. The "if missing" column says exactly what degrades.

Current state: no domain is registered, so the repository points at GitHub
(`github.com/dotnext71/panel`) everywhere a *project link* would otherwise go.
The entries marked **needed to enable a feature** are the ones that will need
a domain once that feature should go live.

---

## 1. Panel-owned hostnames

### 1.1 Needed to enable a feature

| Hostname | Purpose | Configured in | If missing |
|---|---|---|---|
| `support.<domain>` | "Ask dotNext" support chat — a Chatwoot widget loaded into the panel page. Off by default; an admin turns it on per panel. | `frontend/src/lib/supportChat.ts` (`SUPPORT_CHAT_BASE_URL`), CSP allow-list in `agent/commands/nginx_panel_csp.go` (`supportChatOrigin`), copy in `frontend/src/locales/*/common.json`, setting fields in `internal/models/server_settings.go` | Chat is never loaded. The panel works; the toggle would point at a dead origin |
| `enclosed.<domain>` | Diagnostics upload target. The agent redacts and encrypts a support bundle, uploads it, and returns a `#password` link the operator can hand to support. | `agent/commands/system_diagnostic_report.go` (`defaultEnclosedBaseURL`), `frontend/src/hooks/useSupport.ts`, `frontend/src/shells/admin/support/DiagnosticReportModal.tsx`, `support-claim/` (the sibling Enclosed service) | "Send diagnostic report" fails; operators would share logs manually |
| `claims.<domain>` | Hosted-service claim endpoint (a customer proves control of a host to attach it to a hosted account). Served by `support-claim/`. | `support-claim/main.go` (the service itself is deployed at this host) | Hosted-account claiming is unavailable |

### 1.2 Optional convenience hostnames

| Hostname | Purpose | Configured in | If missing |
|---|---|---|---|
| `get.<domain>` | Short install URL — `curl -fsSL https://get.<domain> | bash`. A 302 to the raw `bootstrap.sh`. | `deploy/get-installer/` (`_redirects`, `nginx-get.<domain>.conf`, README) | **Already accounted for:** the README and `bootstrap.sh` use the GitHub raw URL, which is the supported install path. This directory is ready for when the domain exists |
| `translate.<domain>` | Translation portal (Weblate). | `README.md` "Translations" section (currently states the portal is not published) | Translators contribute locale JSON through the repo instead |
| `<domain>` (apex) | Marketing / documentation site. | `frontend/src/shells/{admin,user}/QuickStartModal.tsx` | **Already accounted for:** both modals link to the repository README |

### 1.3 Operational hostnames (only if you run that infrastructure)

| Hostname | Purpose | Configured in | If missing |
|---|---|---|---|
| `git.<domain>` | The retired self-hosted forge. Kept only as a *detection* string so `dotnext update` repoints an old `origin` at GitHub. | `internal/panelcmd/update.go`, `agent/commands/system_update_check.go` | A host whose git remote still points at the old forge keeps a frozen remote until it is repointed by hand |

---

## 2. Automatic (panel-generated) — no purchase required

These are created by the panel for each tenant domain. You never buy them;
they appear in DNS when the tenant adds a domain and enables the matching
service.

| Pattern | Created when | Code |
|---|---|---|
| `mail.<tenant-domain>` | Mail is on for the domain (webmail vhost) | `infrastructure/nginx/dotnext-mail-vhost.conf.tmpl`, mail reconciler |
| `mx.<tenant-domain>` | The panel is the mail exchanger for the domain | DNS reconciler (`internal/reconciler/`) |
| `mta-sts.<tenant-domain>` | MTA-STS policy is published for the domain | `agent/commands/` mail policy handlers |
| `autodiscover.<tenant-domain>`, `autoconfig.<tenant-domain>` | Client auto-configuration records | DNS reconciler |
| `<selector>._domainkey.<tenant-domain>` | DKIM key for the domain (selector is generated per domain) | DKIM handlers |
| `_dmarc.<tenant-domain>`, SPF TXT | Deliverability policy | DNS reconciler |

Panel-generated names carry no branding: they are the customer's own domain.

---

## 3. Local / development-only names

| Name | Purpose |
|---|---|
| `*.dotnext-panel.local` (`mx.`, `mail.`, `enclosed.`) | Loopback test hostnames in ADRs, plans and dev runbooks. Not routable, nothing to buy |

---

## 4. When a domain is bought, this is the checklist

1. Register the domain (any registrar).
2. Point these at whatever runs the matching service, in the order the
   features are needed:
   - `get.` → the installer redirect (see `deploy/get-installer/README.md`)
   - `support.` → the Chatwoot instance
   - `enclosed.` → the Enclosed instance (`support-claim/` documents the sibling service)
   - `claims.` → the `support-claim` binary
   - `translate.` → the Weblate instance
   - apex → the documentation/marketing site
3. Update the three constants listed in §1.1 and restart `dotnext-api` /
   `dotnext-agent`.
4. Re-enable the translation badge in `README.md`.
5. Re-point the project links that currently use GitHub:
   `install.sh` (`REPO_URL`), `bootstrap.sh` (`REPO_WEB`, `DL_BASE`,
   `API_BASE`), `internal/panelcmd/update_repo_owner.go`
   (`canonicalRepoURL`), and the README install commands.

Step 5 is optional: the GitHub repository can stay the source of truth
indefinitely.
