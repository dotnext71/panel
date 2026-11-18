<div align="center">

  <p>
    <a href="#gh-light-mode-only">
      <img src="panel-ui/public/images/dotnext_logo.svg" alt="dotNext" width="320">
    </a>
    <a href="#gh-dark-mode-only">
      <img src="panel-ui/public/images/dotnext_logo_dark.svg" alt="dotNext" width="320">
    </a>
  </p>

# dotNext

**A modern, multi-tenant web hosting control panel built with Go and React.**

<p>
  <a href="#installation">Install</a>
  &nbsp;|&nbsp;
  <a href="#architecture">Architecture</a>
  &nbsp;|&nbsp;
  <a href="#cli">CLI</a>
  &nbsp;|&nbsp;
  <a href="docs/LICENSING_AUDIT.md">Licensing</a>
</p>

<p>
  <img src="https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white" alt="Go 1.27">
  <img src="https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black" alt="React 19">
  <img src="https://img.shields.io/badge/Ant_Design-6-0170FE?logo=antdesign&logoColor=white" alt="Ant Design 6">
  <img src="https://img.shields.io/badge/MariaDB-11-003545?logo=mariadb&logoColor=white" alt="MariaDB 11">
  <img src="https://img.shields.io/badge/Debian-12%20%7C%2013-A81D33?logo=debian&logoColor=white" alt="Debian 12/13">
  <img src="https://img.shields.io/badge/AlmaLinux%20%2F%20Rocky-8%20%7C%209%20%7C%2010-0F4266?logo=almalinux&logoColor=white" alt="AlmaLinux / Rocky 8/9/10">
  <img src="https://img.shields.io/badge/License-MIT-yellow" alt="MIT">
</p>

<p>
  Multi-tenant isolation. Root-safe automation. Database-driven reconciliation.
  dotNext separates its unprivileged HTTP/API &amp; SPA process from a root-owned
  Unix-socket host agent (`dotnext-agent`) and administrative CLI (`dotnext-cli`).
</p>

</div>

## Demo mode (build-tag gated, on `main`)

**Demo mode** lives on `main`, gated **out of production artifacts at compile time**
(see [ADR-0160](docs/adr/0160-build-tag-demo-mode.md)). It is not a runtime toggle
and there is no separate demo branch to maintain.

What demo mode adds (present only in a demo build):

- write-blocking middleware — every non-idempotent `/api/v1/*` request
  (POST/PUT/PATCH/DELETE) returns `403 {"error":"demo_mode"}`, so visitors can
  browse every read endpoint without ever reaching the agent or a DB write;
- a `/info` endpoint that exposes demo showcase credentials;
- a fixed **DEMO banner** + role quick-entry controls.

**Why this is safe on `main`.** The demo Go code carries `//go:build demo` and the
demo UI is gated behind `import.meta.env.VITE_DEMO === "1"`, so a **production build
contains none of it** — the write-block, the credential-exposing `/info`, and the
banner/login override are physically absent from a non-demo binary and SPA bundle.
CI verifies this on every change (`make demo-guard`).

**Operating demo mode:** mark the host with `echo demo > /etc/dotnext/deploy-profile`.
Production hosts leave the file absent and never carry demo code.

## Installation

Install on a fresh supported Debian 12/13, Ubuntu 22.04/24.04, or AlmaLinux / Rocky Linux 8, 9, 10 server:

```
curl -fsSL https://raw.githubusercontent.com/dotnext71/panel/main/bootstrap.sh | sudo bash
```

`bootstrap.sh` downloads the signed release tarball, extracts `dotnext-installer`
and `install.sh`, and starts the interactive TUI installer. Add installer args after `-s --`,
e.g. `--dry-run`. `DOTNEXT_RELEASE_CHANNEL=development` installs the newest main
build instead and sets the server to the development channel.

Classic / scripted install (no TUI — the proven bash installer directly):

```
curl -fsSL https://raw.githubusercontent.com/dotnext71/panel/main/install.sh | sudo bash
```

It builds the current stable release from source, with that release's own
`install.sh`. To build the newest main instead, and set the server to the
development channel, run it as `sudo DOTNEXT_RELEASE_CHANNEL=development bash`.

Both run the same engine. `install.sh` fetches Go 1.27, builds the panel + agent
binaries, builds the SPA with Vite, writes systemd units, provisions MariaDB +
Redis + PowerDNS + Stalwart + Bulwark + CrowdSec (per the selected modules), and
smoke-tests `/health`. Idempotent — safe to re-run; update an installed server
with `dotnext update`. Set `DOTNEXT_MODULES=…`
(comma list) or leave the TUI to build it for a minimal / modular install. The
profiles and module keys are in [Choosing modules](docs/site/installation.md#choosing-modules).

Optional flags:

- `--debug` show full output instead of spinner
- `--hostname <fqdn>` override auto-detected hostname
- `DOTNEXT_HOSTNAME=<fqdn>` env-var equivalent for unattended installs

Uninstall (rolls back system packages, optionally keeps `/home`):

```
curl -fsSL https://raw.githubusercontent.com/dotnext71/panel/main/bootstrap.sh | sudo bash -s -- --uninstall
```

After install:

- Admin panel: `https://your-host:8443/dotnext-admin`
- User panel: `https://your-host:8443/dotnext-panel`
- Webmail: `https://mail.your-domain/`

The panel API listens on a Unix socket; nginx terminates TLS on `:8443` and
proxies upstream. If nginx goes down, the panel and agent stay running so
operators can recover via `dotnext` CLI without losing in-flight state.

## Highlights

- Per-user Linux accounts with per-user PHP-FPM master + cgroup v2 + POSIX quota
- SSH shell access always inside a mandatory sandbox — default `bubblewrap`
  (namespace + bind-mount jail; hides other tenants, sockets, and PIDs);
  fail-closed to `nologin`, never a bare host shell
- FTP/SFTP subaccounts (opt-in) with same-uid or true separate-uid isolation,
  per-subaccount WebDAV, and an opt-in PAM-gated vsftpd FTP(S) module
- Root agent for SSL, mail, DNS, backups, migrations — fronted by a typed
  NDJSON RPC contract over `/run/dotnext-agent/agent.sock` (no shelling out
  from the panel)
- DB-as-truth model: a 60s reconciler reads the panel DB and converges nginx
  vhosts, PHP pools, mailboxes, DKIM, DNS, SSL, mta-sts, and per-user limits
- cPanel and WHM migrations (analyse → fix-perms → validate → restore) with
  preserved MySQL users + password hashes
- IMAP sync for migrating mail from external servers
- Stalwart Mail Server with browser-trusted IMAPS/465/587 (LE cert pushed
  into Stalwart Certificate object) and self-deleting SSO file (Installatron
  pattern) for one-click webmail
- Per-mailbox forwarders, autoresponders, shared folders, disclaimers
- Bulwark (Next.js JMAP) webmail with same-origin per-tenant routing
  (nginx `sub_filter` rewrites panel hostname → `$host` so the SPA stays
  same-origin on `mail.<tenant>`)
- PowerDNS authoritative + recursor with native DNSSEC (per-domain toggle)
- Per-domain listen-IP binding (M24 IP Manager) with reserved-word-safe
  migrations against MariaDB 11.x
- Per-domain opt-in nginx FastCGI micro-cache with safe-bypass for cart /
  admin / authenticated cookies
- Restic backups (account_full + system_backup) with encryption, dedup,
  SFTP / S3 / B2 / Azure / GCS / REST destinations, scheduled + on-demand,
  per-destination slot cap + circuit breaker, opt-in auto daily backups
- Portable Full Server container backup + restore-from-upload, tenant
  self-service restore (files / DB / mail), single-domain docroot restore,
  and per-database Download + Restore-from-file for MariaDB **and**
  PostgreSQL (chunked async upload; pgAdmin custom/tar formats accepted)
- WordPress 1-click install / delete / clone (M10) — 15-app catalogue (M19)
  incl. Moodle / Joomla / NextCloud / OpenCart / Mautic / Drupal
- Per-user resource limits: cgroups v2 slice drop-in, nginx limit_req,
  POSIX quota — admin toggleable, reconciler-converged
- Integrated security suite: CrowdSec parsers + AppSec WAF + per-user
  egress firewall (nftables + cgroupv2-vmap) + ModSecurity replaced by
  CrowdSec AppSec (ADR-0060) + LMD + ClamAV-on-demand + YARA + Tetragon
  for malware detection + dotnext quarantine + M14 notifications dispatch
- CrowdSec 1.8 AppSec bot-detection challenge (default OFF; per-domain
  opt-in / opt-out; tenant self-service on own domains), step-up (MFA)
  auth for the admin File Manager + Root Terminal, and secret-rotation
  tooling (DB app-user / JWT / foundation)
- 6-channel notifications: Discord, ntfy, Web Push (VAPID), SMS,
  Email, Webhook, Slack, in-app bell — 4 event sources incl. cert renew,
  disk full, service down, CrowdSec spike
- One-time login tokens (CLI + dashboard) with IP binding
- Magic-link-free SSO between panel and webmail (no Hydra / OIDC overhead —
  ADR-0040 supersedes the M16 Hydra rollback)
- Audit logs, account activity feed, encrypted diagnostic-log sharing to
  support, in-app updates + support tabs

## Feature Map

### Admin Panel

- Dashboard with stats, health, recent activity, notifications bell
- User management with suspension, packages, quotas, impersonation,
  in-place username rename, at-a-glance resource + bandwidth counts
- Server settings (hostname, nameservers, public IPs, panel cert)
- Service manager for systemd services + start/stop/restart
- PHP version and per-user pool management (server-wide extensions tab)
- DNS zones, templates, DNSSEC, secondary NS
- SSL issuance and renewals (panel cert + per-domain)
- IP address assignments (managed IP pool, per-domain bind)
- Backups: account_full + system_backup, local + remote (SFTP/S3),
  schedules, encrypted destinations
- Migrations (cPanel restore, WHM downloads, IMAP sync)
- Security: CrowdSec allowlists / alerts / console / captcha + UFW + AppSec
  geoblock + AppSec bot-detection (default OFF, per-domain scoped) + per-user
  egress firewall + malware quarantine + secret-rotation tooling
- Updates + Support tabs: live `dotnext update` with transient systemd units,
  enclosed-encrypted diagnostic sharing to webmaster
- Server status (CPU / mem / disk / queues / 5s polling)
- Database admin ops (curated tuner, root password, processlist,
  pmaAdmin SSO)
- Email queue, throttles, MTA-STS, outbound reports
- Audit logs, notifications dispatcher, malware / quarantine events
- Notification channel admin (test-send, scopes, throttles)

### User Panel

- Domains with independent Web / Mail / DNS service flags, redirects,
  custom nginx rules, listen-IP, FastCGI micro-cache, reverse-proxy mode,
  per-domain environment variables, custom document root
- DNS records editor with conflict detection (CNAME exclusivity per
  RFC 1034 §3.6.2) and dedup
- Mail Domains drill-down: mailboxes, forwarders, autoresponders,
  catch-all, disclaimers (HTML), shared folders, CalDAV/CardDAV (with
  per-domain server override), per-domain logs + statistics, mail-only delete
- IMAP sync (single + bulk)
- Webmail SSO (Bulwark, Next.js JMAP)
- WordPress (install, update, scan, SSO) + 14 other 1-click apps
- File manager (AntD-native, chunked upload with cancel + live speed,
  async copy/move/extract progress) + SFTP + SSH keys
- SSH shell access inside a mandatory bubblewrap sandbox
- FTP/SFTP subaccounts (opt-in) + WebDAV
- Databases (MariaDB + Postgres, tabbed) with phpMyAdmin (MariaDB) /
  Adminer (Postgres) SSO + per-database Download + Restore-from-file
- PHP settings per account + per-domain (OPcache/JIT, Xdebug, Composer,
  extra extensions, env vars)
- SSL management (Let's Encrypt + custom certificates)
- Cron jobs (systemd-user timers + allowlist: php / wp / python / node)
- Backups + restore (account_full) + restore-from-upload
- Tenant-owned notification channels (when enabled by the admin)
- Logs, statistics, bandwidth usage (daily nginx-log sync)
- Support access link generator (one-time IP-bound tokens)
- Notification preferences (Discord, ntfy, Web Push, SMS, Email)

### Platform

- Root-level agent (`panel-agent`) with typed NDJSON RPC handler registry
- Reconciler (60s tick) converges domain.create / SSL / DKIM / vhosts /
  PHP pools / nginx rate-limits / mailboxes / mtasts / ssh keys / cron
- Job queue: async backup + migration steps + WordPress install
- Health monitor with notification dispatch on service down / cert near-
  expiry / disk-full / CrowdSec spike / queue depth
- Redis (Unix socket, ACL-scoped) for cache, sessions, notifications
  dispatcher streams
- Per-domain opt-in FastCGI micro-cache + manual purge
- Opt-in read-only Automation API on :443 for firewalled billing hosts
- Disaster-recovery standby sync + promote (status card, stall alerts)
- Zabbix 7 monitoring template; OS security auto-updates on by default
- Side-nav badge counts (tenant + admin); admin header Health badge
- Multi-language UI (en default; i18n harness ready)

## Architecture

- **Control plane**: Go binary `panel-api` (Gin) listening on
  `/run/dotnext-panel/api.sock`. Embeds the SPA and serves it from `/`
- **Data plane**: Go binary `panel-agent` running as root, listening on
  `/run/dotnext-agent/agent.sock` (0660, group `dotnext`). Typed NDJSON RPC
  registry — every privileged op (nginx reload, certbot, systemctl,
  mysql DDL, file ops) is a named handler the panel calls by name. No
  shelling out from the panel itself
- **State plane**: MariaDB `dotnext_panel` (single DB, single writer = panel-
  api). Reconciler reads the DB every 60s and converges the host
- **Job plane**: Redis Streams dispatcher (notifications, backups, mail-
  scan)
- **Frontend**: React 19 + Ant Design 5 + TanStack Query, built by Vite,
  served from the Go binary's embedded FS — single deploy unit, no Node
  runtime on the host
- **Webmail**: Bulwark (Next.js JMAP) at `/opt/dotnext-webmail`, served on
  `mail.<tenant>` per-domain via nginx → Unix socket
- **SSH shell**: mandatory sandbox per login — default `bubblewrap`
  (namespace + bind-mount jail via setuid `bwrap`); an `nspawn` mode is
  selectable but not yet wired. Fail-closed to `nologin`, never a bare shell
- **Security**: CrowdSec parsers + AppSec (nginx-bouncer Lua, WAF) +
  per-user egress firewall (nftables + cgroupv2-vmap, ADR-0084) + LMD +
  ClamAV-on-demand + YARA + Tetragon
- **Logging**: structured JSON via slog; nginx access logs feed CrowdSec
- **Server metrics**: live `/proc` reads, no Prometheus exporter dependency

Service stack (single-node default):

- **panel-api** (Go, Unix socket, embedded SPA)
- **panel-agent** (Go, Unix socket, root)
- **nginx** (TLS terminator on `:8443`, user vhosts on `:80`/`:443`,
  per-domain mail vhost on `:443`, FastCGI cache keyzone shared, AppSec
  bouncer Lua)
- **MariaDB** (Unix socket only — `skip-networking`)
- **Redis** (Unix socket, mode 0660, group `dotnext-sockets`)
- **PowerDNS authoritative** (split-port :5300, MySQL backend) +
  **pdns-recursor** (loopback :53, resolver chain)
- **Stalwart Mail Server** (SMTP / IMAP / 465 / 587 / 993 / JMAP /
  ManageSieve, LE-cert pushed into Certificate object)
- **Bulwark** (Next.js JMAP webmail, Unix socket, served per-tenant)
- **Kratos** (Unix socket admin + public, sole auth source — M20)
- **CrowdSec** (LAPI socket + AppSec :7422 + nginx-bouncer Lua)
- **Restic** (encrypted, dedup, backup destinations)
- **systemd-user** (cron jobs as user-scope timers)

## Requirements

- Fresh Debian 12/13, Ubuntu 22.04/24.04, or AlmaLinux / Rocky 8/9/10 install (no pre-existing web or mail stack)
- 2 GB RAM minimum, 4 GB recommended — the same guidance as cPanel (on a ≤4 GB host the installer adds swap and caps the on-box SPA/Go build memory so it doesn't OOM; on 2 GB the build is slower but completes)
- A domain for panel + mail (with glue records if hosting DNS)
- PTR (reverse DNS) for mail hostname
- Open ports: 22, 80, 443, 8443, 25, 465, 587, 993, 995, 53

## Security Hardening

See [`docs/adr/`](docs/adr/) for the full architectural-decision record
(110+ ADRs covering every load-bearing design choice). Highlights:

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `DOTNEXT_HOSTNAME` | Override auto-detected panel hostname during install | (auto) |
| `DOTNEXT_PANEL_BIND` | Override panel-api listen socket | `/run/dotnext-panel/api.sock` |
| `DOTNEXT_AGENT_SOCKET` | Override agent RPC socket | `/run/dotnext-agent/agent.sock` |
| `DOTNEXT_TEST_DATABASE_URL` | Real MariaDB DSN for integration tests | (unset) |
| `DOTNEXT_LOG_LEVEL` | Slog level (debug / info / warn / error) | `info` |
| `TLS_CERT` / `TLS_KEY` | Cleaned from panel.env on update — nginx terminates | (auto-cleaned) |

### Key Security Features

- Panel never runs as root. Every privileged op crosses the agent Unix
  socket as a typed RPC call; agent verifies caller via `SO_PEERCRED`
- Shell arguments validated + escaped per-handler (no `sh -c $arg`
  patterns); domain names validated against `validateDomainNameForShell`
- DKIM private keys + SSO tokens + mailbox plaintexts encrypted at rest
  via AES-GCM with a per-host SSO key (`/etc/dotnext-panel/sso.key`)
- One-time admin SSO tokens are 256-bit, single-use, 5-minute TTL,
  reaped every 30s by systemd timer (ADR-0040 webmail SSO file pattern)
- Stalwart Certificate object pushed from LE-renewed cert on each
  certbot deploy-hook — IMAPS / 465 / 587 always serve browser-trusted
  cert (no rcgen self-signed fallback)
- CrowdSec AppSec WAF + per-user egress firewall (cgroupv2-vmap, ADR-0084)
- Self-healing reconciler — config drift on disk is reverted on next
  tick; operator hand-edits to nginx vhosts are lost-by-design
- CSP, HSTS, SameSite cookies, X-Forwarded-Proto handled by nginx
- Migrations are schema-only (no app-populated tables seeded by SQL)
- Audit log on every privileged mutation; actions performed while impersonating are attributed to the target user
- Pre-commit + CI gates: `go vet`, `go test -race ./...`, `npx tsc -b`,
  `bash -n install.sh`, Playwright E2E, AppSec geoblock golden tests

## Updates

Update the panel (code, dependencies, DB migrations, infrastructure):

```
dotnext update
```

This pulls the latest code, rebuilds the panel + agent binaries, rebuilds
the SPA, applies golang-migrate migrations, syncs nginx vhosts + systemd
units + PHP config + CrowdSec acquis, and restarts the panel + agent.
Safe to run on a live server — the reconciler tolerates a brief panel
restart and converges state on the next tick.

Self-heal a broken install (7 detectors, --diagnose default,
--auto safe, --all --yes destructive):

```
dotnext repair --diagnose       # report only
dotnext repair --auto           # fix safe issues
dotnext repair --all --yes      # destructive recovery
```

## CLI

The `dotnext` command uses a noun:verb pattern. All commands support
`--json` for machine-readable output. Commands that confirm before a risky
step take `--yes`, which skips the question or, where the command does a dry
run first, applies the change.

<!-- BEGIN generated: go test ./internal/panelcmd -run TestREADMECLICommands -update -->
```
dotnext-cli active-tasks          # Show running/queued tasks (updates, backups, malware scans)
dotnext-cli admin                 backfill-usernames|purge-orphan-identities|
                                  rebuild-kratos|relabel-identifiers|
                                  slice-cutover
dotnext-cli aide                  check|rebuild|status
dotnext-cli app                   cache|cache-doctor|cache-doctor-run|clone|
                                  delete|e2e|get|install|list|magic-link|
                                  refresh-cache-plugin|registry|scan
dotnext-cli apparmor              flip-mature|status
dotnext-cli appsec                exclusion|explain|flarum-sync|host-mode|
                                  render-config
dotnext-cli audit                 prune|query|verify
dotnext-cli automation-token      list|mint|revoke
dotnext-cli backup                account-restore|copy|destination|retention|
                                  schedule|scheduler
dotnext-cli cron                  add|delete|http-trigger|list|run-now|update
dotnext-cli crowdsec              alerts|allowlists|blocklists|bouncers|captcha|
                                  country-exempt|decisions|geoblock|hub|metrics|
                                  profiles|status
dotnext-cli db                    backup|chown|config|create|delete|kill|list|
                                  maintenance|maintenance-status|postgres|
                                  processes|restore|root-password|sso|user
dotnext-cli disk-usage            refresh|refresh-all|show
dotnext-cli dns                   prune-orphan-records|prune-service-records|
                                  record|zone
dotnext-cli docker                disable|enable|enable-tenant|status
dotnext-cli docker-app            backup-create|backups|catalog|delete|edit|env|
                                  exec|install|list|logs|maintenance|patch|
                                  rebuild|restart|restore|start|status|stop|
                                  update
dotnext-cli domain                bandwidth|bot-exclusion|browse|catchall|chown|
                                  create|delete|dir-privacy|disable|disclaimer|
                                  email-disable|email-dkim-rotate|email-enable|
                                  enable|fix-perms|htaccess-preview|ip-acl|list|
                                  mta-sts|ownership|php-settings|php-version|
                                  prune-orphans|set|show|skip-auto-san|whois
dotnext-cli dr                    feed|pair|promote|status|unfeed|unpair
dotnext-cli files                 archive|chmod|copy|delete|download|extract|
                                  list|mkdir|move|read|rename|stat|upload|write
dotnext-cli impersonation         create|end|list
dotnext-cli ip                    add|delete|list|update
dotnext-cli limits                apply|check|override|package|status|usage
dotnext-cli log                   access|types
dotnext-cli mail                  deliverability|logs|queue|throttle
dotnext-cli mail-cert             disable|enable|reissue|show
dotnext-cli mail-group            add-member|create|delete|list|remove-member|
                                  set-external|set-members|set-resources|show
dotnext-cli mailbox               autoresponder|create|delete|forwarder|list|
                                  passwd|set-quota|shares
dotnext-cli malware               events|quarantine|scan|scan-status|settings|
                                  status|update-signatures|yara
dotnext-cli malware-purge         # Hard-delete terminated malware quarantine rows past retention
dotnext-cli migrate               force|imap|import|import-wp|pull-source|
                                  reap-secrets|refresh|restore|status|up
dotnext-cli migrate-from-dotnext  # Migrate an existing legacy DotNext Panel host layout to dotNext
dotnext-cli notification          broadcast|channels|dlq|events|inbox
dotnext-cli nspawn                build|list|prune
dotnext-cli package               create|delete|edit|list
dotnext-cli page-template         get|list|reset|set
dotnext-cli panel-cert            issue|new-key|status|toggle
dotnext-cli panel-primary         ensure
dotnext-cli pdns                  backfill|dnssec
dotnext-cli per-user-egress       approve|deny|flip-mature|get|requests|
                                  set-policy|summary
dotnext-cli php                   ext|pool|version
dotnext-cli php-defense           incidents|incidents-clear|mode|rule-toggle|
                                  rules|status
dotnext-cli python-app            create|delete|env|frameworks|list|logs|
                                  rebuild|restart|start|stop
dotnext-cli release               promote
dotnext-cli repair                # Detect and fix known deployment-host issues
dotnext-cli retention-sweep       # Prune expired rows from unbounded log/report tables
dotnext-cli secrets               rotate
dotnext-cli security              feed-check
dotnext-cli serve                 # Start the dotNext Panel HTTP(S) server
dotnext-cli service               action|status
dotnext-cli session               list|revoke|revoke-user
dotnext-cli settings              get|keys|mail-hostname|set
dotnext-cli shared-resource       create|grant|grants|list|remove|revoke|scanner
dotnext-cli ssh-key               add|delete|list
dotnext-cli ssl                   disable|enable|list|readiness|renew|retry|
                                  set-custom|shared
dotnext-cli sso                   prune-tokens|rotate-key
dotnext-cli sso-reap              # Sweep stranded dotnext-sso-<nonce>.php files
dotnext-cli system                diagnostic|info|process|resolver|restore|
                                  services|ssh-keys
dotnext-cli ufw                   default|disable|enable|migrate-ip-bans|rule|
                                  status
dotnext-cli update                apt-refresh
dotnext-cli user                  2fa-reset|create|delete|list|password|rename|
                                  reprovision|ssh-forwarding|suspend|unsuspend
dotnext-cli user-token            list|mint|revoke
dotnext-cli version               # Print dotNext version, commit, and runtime info
```
<!-- END generated -->

`dotnext domain create` takes `--web-enabled` / `--manage-dns` / `--mail`
(independent per-domain services), and `dotnext system restore --from-tar`
restores a whole server from a Full Server container archive. The complete,
generated flag reference lives in
[`docs/site/platform/cli-reference.md`](docs/site/platform/cli-reference.md).

See [`docs/CONVENTIONS.md`](docs/CONVENTIONS.md) for the full repo-wide
patterns (route families, SearchableTable, Drawer-for-CRUD, list envelope,
rate limits) and [`docs/adr/`](docs/adr/) for every load-bearing decision.

## Development

```
make build              # compile panel-api + panel-agent
make run                # run panel-api (dev, embedded SPA)
make test               # all Go tests, race detector on
make test-coverage      # coverage report (internal packages)
make test-integration   # needs DOTNEXT_TEST_DATABASE_URL + real MariaDB
make coverage-check     # fail if combined coverage < 80%
make lint               # golangci-lint v2
make fmt                # go fmt + vet
```

Frontend dev (from `panel-ui/`):

```
npm install
npm run dev             # Vite on http://localhost:5173
                        # proxies /api and /health to 127.0.0.1:8443
```

E2E (from `panel-ui/`):

```
npm run test:e2e        # Playwright against the dev server
```

See [`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md) for the full feature
development workflow (research → plan → TDD → review → ship).

### Versioning

The version string lives in `VERSION` (read at build time and exposed
via `/health`). When the installer clones the repo for a fresh install,
it reads `VERSION` to display the installed version. Always bump `VERSION`
in the same commit as the corresponding `install.sh` changes — drift
shows up as a mismatched footer and installer banner.

## Translations

<!-- The Weblate portal (translate.dotnext-panel.com) has no public URL yet:
     re-enable its status badge here once it is live. -->

The panel UI is translated. English is the source language, and every other
locale lives in `frontend/src/locales/` and ships with the repo.

Shipping locales: English (source), Arabic, Simplified Chinese, French, German,
Hebrew, Italian, Japanese, Brazilian Portuguese, Romanian, Russian, Spanish,
Turkish and Ukrainian. Arabic and Hebrew render right-to-left.

### For translators

English strings live in `frontend/src/locales/en/common.json`; other locales
are translations of that catalog. You do not need to code or set up a local
build to contribute a locale — send the translated JSON (or a pull request)
and it merges as-is. The short version:

- Do not translate product or protocol names (`nginx`, `MariaDB`, `Stalwart`,
  `PHP-FPM`, `DNS`, `SSL`, ...) or anything in `backticks`.
- Keep placeholders such as `{{count}}` exactly as they appear. You may move
  one to fit your word order, but never rename or drop it.
- Keep strings short. Most are buttons, menu items, table headers and form
  labels, and a long translation breaks the layout.
- Hebrew and Arabic render right-to-left. If a screen looks broken, report it
  rather than shortening the text to work around it.
- If an English string is ambiguous or wrong, open an issue instead of
  guessing — it gets fixed at the source, once, for every language.

Missing your language? Open an issue and we will add it.

### For developers

- **`panel-ui/src/locales/en/common.json` is the only catalog edited by hand.**
  Weblate owns every other locale file; hand edits there will be overwritten.
- User-facing text goes through `t("some.key")` (`react-i18next`), never a
  hardcoded literal. Add the key and its English text to `en/common.json` in the
  same commit that introduces it.
- `panel-ui/src/i18n.ts` holds the locale list plus the AntD and dayjs locale
  maps and the RTL set. Adding a locale means touching all three.
- A key with no translation yet falls back to the English source, so a partially
  translated locale never renders blank labels.

## Repository Layout

```
dotnext-panel/
├── panel-api/          # Go HTTP server (Gin) + reconciler + agent RPC client
│   ├── cmd/server/     # main entry
│   ├── internal/       # api/, auth/, repository/, reconciler/, config/, ...
│   └── migrations/     # golang-migrate SQL (000xxx_*.up/.down.sql)
├── panel-agent/        # Go binary running as root; typed NDJSON RPC handlers
│   ├── cmd/dotnext-agent/
│   └── internal/commands/
├── panel-ui/           # React SPA (AntD + TanStack Query)
│   ├── src/            # shells/, components/, theme/, pages/, ...
│   └── public/
├── agentwire/          # NDJSON RPC types shared by panel-api + panel-agent
├── internal/           # shared Go libs (cronvalidate, dbtuning, phpext, ...)
├── install/            # install.sh assets (nginx tmpl, stalwart plan,
│                       # letsencrypt deploy hooks, bulwark env, ...)
├── docs/               # CONVENTIONS, BLUEPRINT, adr/, runbooks/, KNOWN_ISSUES
├── plans/              # per-milestone implementation blueprints
├── .github/workflows/  # CI (Go + vitest + E2E)
├── install.sh          # single-supported install path (curl | sudo bash)
├── config.example.toml # reference config (copied to /etc/dotnext-panel/)
├── Makefile            # build / test / lint targets
└── go.mod              # Go workspace root
```

## License

MIT License (`MIT`) — see [LICENSE](LICENSE), [THIRD_PARTY_NOTICES](THIRD_PARTY_NOTICES), and [docs/LICENSING_AUDIT.md](docs/LICENSING_AUDIT.md). Component-level licence metadata (SPDX headers, UI `package.json`) matches: `MIT`. Companion WordPress plugins under `wp-plugins/` remain licensed under `GPL-2.0-or-later`.

## Mail Subdomain

Visiting `mail.<your-domain>` in a browser routes to webmail (Bulwark) via
the per-domain nginx vhost. The vhost installs an nginx `sub_filter` that
rewrites the panel hostname to the requested `$host` in Bulwark's
`/api/config` and Stalwart's `/.well-known/jmap` responses, so the SPA
stays same-origin on `mail.<tenant>` and Stalwart's JMAP Session URLs
never leak the panel hostname.

`autodiscover` / `autoconfig` paths are excluded so mail client
auto-discovery (Thunderbird, Outlook) keeps working.

## Documentation

Local documentation lives under [`docs/`](docs/):

- [Conventions](docs/CONVENTIONS.md) — repo-wide patterns (route families,
  SearchableTable, Drawer for create+edit, icon shim, list envelope, rate
  limits) + anti-patterns learnt the hard way
- [Blueprint](docs/BLUEPRINT.md) — full feature map + milestone roadmap
- [ADRs](docs/adr/) — every load-bearing architectural decision (110+)
- [User & admin guides](docs/site/) — per-feature docs (domains, mail,
  databases, backups, security, PHP, cron, files, SFTP/FTP, and more)
- [Plans](plans/) — per-milestone implementation blueprints
- [Runbooks](docs/runbooks/) — operational guides (SSL, mail, backups,
  secondary nameserver, DNS, applications)
- [Known Issues](docs/KNOWN_ISSUES.md) — caveats + workarounds
- [Contributing](docs/CONTRIBUTING.md) — feature development workflow
- Translations — Weblate portal (not published yet); English strings in `frontend/src/locales/en/`;
  is the source catalog (`panel-ui/src/locales/en/common.json`)
- [Environment](docs/ENV.md) — full env-var reference
