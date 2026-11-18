# dotnext-panel.com Docs — Source Set

Comprehensive English docs for the current (Go) generation of dotNext. Designed to be the source-of-truth for the dotnext-panel.com documentation site.

## Contents

| Page | What |
|---|---|
| [index.md](./index.md) | What DotNext does, who it's for, where to start. |
| [quickstart.md](./quickstart.md) | New host → first user → first SSL in 7 steps. |
| [installation.md](./installation.md) | Requirements, install command, env vars, what runs. |
| [admin.md](./admin.md) | The `/dotnext-admin` shell, all pages. |
| [user.md](./user.md) | The `/dotnext-api` (tenant) shell. |
| [domains.md](./domains.md) | Per-domain settings + lifecycle. |
| [ssl.md](./ssl.md) | Let's Encrypt automation. |
| [dns.md](./dns.md) | PowerDNS auth + recursor + DNSSEC operator notes. |
| [mail.md](./mail.md) | Stalwart mail stack, mailboxes, deliverability. |
| [php.md](./php.md) | Per-user pools, version management, extensions. |
| [wordpress.md](./wordpress.md) | 1-click WP install / clone / admin SSO. |
| [applications.md](./applications.md) | 16-app registry, install pipeline. |
| [files.md](./files.md) | In-panel AntD File Manager. |
| [sftp.md](./sftp.md) | SFTP via OpenSSH `Match Group`, sandboxed SSH shell (M13 bubblewrap), SSH key vault. |
| [cron.md](./cron.md) | systemd-user timers + command allowlist. |
| [databases.md](./databases.md) | MariaDB / PostgreSQL + phpMyAdmin SSO + admin ops. |
| [backups.md](./backups.md) | account_full + system_backup + destinations + schedules. |
| [security.md](./security.md) | CrowdSec, AppSec, AppArmor, Snuffleupagus, AIDE, egress, malware. |
| [firewall.md](./firewall.md) | UFW (port baseline) + CrowdSec (IP trust). |
| [resource-limits.md](./resource-limits.md) | Quota + cgroup v2 + nginx limit_req. |
| [notifications.md](./notifications.md) | Channels + event sources + routing. |
| [server-status.md](./server-status.md) | The Server Status dashboard. |
| [updates.md](./updates.md) | `dotnext-cli update` flow. |
| [support.md](./support.md) | Encrypted diag bundle. |
| [ip-manager.md](./ip-manager.md) | Per-domain listen IP + apex DNS. |
| [migrations.md](./migrations.md) | cPanel / DirectAdmin / Hestia / WHM ingest. |
| [operations.md](./operations.md) | Day-2 runbook. |
| [troubleshooting.md](./troubleshooting.md) | Symptom → cause → fix. |
| [removed-features.md](./removed-features.md) | What's intentionally not here (OIDC, ModSec, filebrowser, impersonation, Refine). |
| [platform/stack.md](./platform/stack.md) | Process / data / reconciler model. |
| [platform/agent.md](./platform/agent.md) | Privileged-agent contract. |
| [platform/cli.md](./platform/cli.md) | `dotnext` CLI cheat-sheet — common commands, hand-written. |
| [platform/cli-reference.md](./platform/cli-reference.md) | Full generated `dotnext` command reference (every subcommand + flag). |
| [platform/dnssec.md](./platform/dnssec.md) | DNSSEC architecture. |
| [platform/mail-autoconfig.md](./platform/mail-autoconfig.md) | Thunderbird / Outlook / Apple autoconfig. |
| [platform/health-monitor.md](./platform/health-monitor.md) | `/health`, `/health/agent` and the automation status endpoints. |
| [platform/monitoring.md](./platform/monitoring.md) | Audit + notifications + Prometheus. |
| [admin/](./admin/) | **57 admin subpages** — one per route / feature on the `/dotnext-admin/*` URL tree. Mirrors the legacy site's `/docs/admin/*` granularity. |
| [user/](./user/) | **34 tenant subpages** — one per route / feature on the `/dotnext-api/*` URL tree. Mirrors the legacy site's `/docs/user/*` granularity. |
| [INTEGRATION.md](./INTEGRATION.md) | How the docs team lands this content into the dotnext-panel.com Astro site: `sync-docs.mjs` MAPPINGS, ui.ts placeholder bootstrap, Astro page wiring, translation flow, sidebar IA. |
| [MAPPINGS-snippet.txt](./MAPPINGS-snippet.txt) | Copy-paste-ready full MAPPINGS list for `sync-docs.mjs` covering every top-level page, platform page, admin subpage, and user subpage. |

## How to land these on dotnext-panel.com

The site (dotnext-panel.com) stores doc HTML inline in `src/i18n/ui.ts` (7 languages) sourced from markdown via `scripts/sync-docs.mjs`. Two paths:

### Option A — extend the existing sync pipeline

1. Drop these files into `/home/shuki/projects/dotnext/docs/` (or repoint sync-docs.mjs at this directory).
2. Extend `MAPPINGS` in `~/projects/dotnext-panel.com/scripts/sync-docs.mjs` to cover the new pages.
3. Run `node scripts/sync-docs.mjs` — updates the `en` `docsHtml.*` keys.
4. Translate the new content into other locales as needed.
5. Add new `*.astro` pages under `~/projects/dotnext-panel.com/src/pages/docs/` mirroring the keys in `ui.ts.docsHtml`.

### Option B — direct paste

For pages where the source MD doesn't need to live anywhere else, render to HTML and paste into `ui.ts.docsHtml.<key>` directly:

```bash
pandoc -f markdown -t html5 --no-highlight --wrap=none site/<page>.md
```

(The `<h1>` is stripped + backticks escaped by sync-docs.mjs — match that processing if pasting by hand.)

## Conventions

- Each page is self-contained: it doesn't assume the reader has read another doc, except where it explicitly links one.
- Domain language matches `/home/shuki/projects/dotnext2/CONTEXT.md`: **Agent**, **Reconciler**, **Privileged DB Admin Action**, **Cron Job**, **Cron Job Intake**, **SSO Token**, **Shadow Account**, **SSO Token Resolution**.
- All "this is intentionally not here" notes for old-panel features live in `removed-features.md` so user docs stay focused on what's shipped.
- No internal codes in the published pages: no milestones (M43), JAB- tickets, ADR numbers, or GitHub issue and PR numbers. Say what the feature does and where it is instead. dotnext-panel.com removes them when it merges a page, so a code added here comes back on the website with the next merge. `TestDocsSiteHasNoInternalCodes` (`internal/panelcmd`) fails on one. Product names such as M365 are fine.
- File paths, systemd unit names, and CLI verbs are quoted exact and verified against the current code.
