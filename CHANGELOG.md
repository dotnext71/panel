# Changelog

## dotNext

### Milestone 12 — AlmaLinux 8 / 9 / 10 (EL8 / EL9 / EL10) Support (2026-11-18)

- Added shared OS detection package `internal/platform/osinfo` covering Debian 12/13, Ubuntu 22.04/24.04, and AlmaLinux 8/9/10 (`FamilyDebian`/`FamilyRHEL`, `PackageManager`, `SSHServiceName`, `RedisServiceName` supporting `redis` on EL8/9 and `valkey` on EL10, `FirewallBackend`, `MACBackend`, `NologinShell`, PowerDNS paths, and Remi SCL PHP path/package resolution).
- Extended `internal/phpext` (`ResolvePackagesForOS`), `panel-agent` (`system.apt_check`, `system.apt_run`, `system.autoupdate_apply`, `service.list`, `security.apparmor`, `security.ufw` backed by `firewalld`, `php.*` Remi SCL version/extension/pool/CLI/Snuffleupagus handlers, `db.postgres.*`, `ftp.*`, `diagnostic`, `pdnsrecursor`, `dotnext-ssh-shell`), and `panel-api` (`eventsources`, `cmd/server`, `internal/api`) for EL8/9/10 compatibility while preserving Debian/Ubuntu and RPC contracts.
- Updated `install.sh`, `deploy/installer/lib/packages.sh`, `install/systemd/{fpm-pre-start,fpm-exec,pdns-local-address}`, and `install/snuffleupagus/build/build.sh` for EPEL + CRB/PowerTools + Remi SCL PHP + NodeSource 22, `dnf`, `firewalld`, SELinux configuration, `aide`, CrowdSec EL RPM repos, and `valkey` (EL10) / `redis` (EL8/9).

### Milestone 11 — Systemd Unit Metadata Rebrand & Final Validation (2026-11-18)

- Updated `Description=` metadata across 24 systemd service, timer, and slice units in `install/systemd/` (`dotNext ...`) while preserving unit names, sockets, and `install.sh` repair invariants.
- Executed full post-edit validation across Vitest (`259` files), TypeScript (`tsc -b`), Vite SPA bundle (`panel-ui/dist`), Go `go vet`, Go `-race` tests, `make demo-guard`, Linux `amd64` cross-compilation of all 7 `bin/dotnext*` binaries, and shell script syntax validation.

### Milestone 10 — Backend API Error Copy & Route Deeplink Alignment (2026-11-18)

- Migrated 38 Go backend files (`cmd/server`, `internal/api`, `eventsources`, `reconciler`, `ownershipops`, `drsync`, `middleware`, `notifications/senders`, `webui`) from the old `/dotnext-admin/*` notification/event deeplinks to canonical `/admin/*` and `/panel/*` routes (preserving legacy path acceptance in `notifications/links.go`).
- Rebranded user-facing API messages and embedded webmail SSO expiration page (`backup_restore_preflight`, `backup_restore_domain_check`, `dns`, `domain_email`, `uploaded_backups`, `webmail_sso`) to dotNext branding and dark-first obsidian/cyan styling.

### Milestone 9 — Deep UI Rebrand, Worker Telemetry Card & PageHeader System (2026-11-18)

- Completed UI copy & locale pack rebrand across `panel-ui/src/` (`QuickStartGuide`, `SupportChatCard`, `SystemUpdatesPage`, `UpdatesLogTab`, `FreeHostnameCard`, `BrandingCard`, `DomainCreate`, `DomainMailProviderSection`, `AdminSecurityAppArmor`, `AdminSecuritySnuffleupagus`, `SupportPage`, `DiagnosticReportModal`, `AdminSystemJobs`, `MCPSetupGuide`, `MigrateRemoteDrawer`, `channelKindConfig`) and `en`/`ru` locale bundles while retaining the legacy-panel migration-source options and compatibility keys.
- Added [PageHeader.tsx](file:///Users/khalequzzaman/Herd/dotnext/panel-ui/src/components/PageHeader.tsx) and [SupervisedWorkersCard.tsx](file:///Users/khalequzzaman/Herd/dotnext/panel-ui/src/shells/admin/server-status/SupervisedWorkersCard.tsx) (plus unit tests) wired into `ServerStatusPage` (`GET /api/v1/admin/system/workers`) and `UserList`.

### Milestone 8 — Documentation & Final Build Validation (2026-11-18)

- Added [docs/COMPATIBILITY_MATRIX.md](docs/COMPATIBILITY_MATRIX.md) capturing module architecture boundaries, backward-compatibility shims, verified test scopes, and Linux vs. macOS host notes.
- Performed final validation across embedded SPA production build (`panel-ui/dist`), TypeScript (`tsc -b`), full Vitest suite (`258` files), Go vet, Go test suites, `make demo-guard`, and Linux `amd64` cross-compilation of all 7 binaries (`bin/dotnext*`).

### Milestone 7 — Modular Installer Libraries, TUI Rebrand & Companion Integrations (2026-11-18)

- Added modular dotNext installer lifecycle helpers under `deploy/installer/lib/`:
  - `common.sh` (structured logging to `/var/log/dotnext-install.log`, root check, OS detection for Debian 12/13 and Ubuntu 22.04/24.04)
  - `packages.sh` (disk space, RAM, and baseline tooling preflight checks)
  - `services.sh` (systemd unit lifecycle, pre-upgrade binary/config snapshot, rollback helper, health verification)
  - `uninstall.sh` (`--uninstall` flow preserving `/home` tenant data unless `--purge-data` is specified)
- Rebranded `installer/` (`dotnext-installer`, self-contained dark truecolor dotNext `logo.ansi`, `"DOTNEXT PANEL"` banner, `DOTNEXT_HOSTNAME` / `DOTNEXT_MODULES` / `DOTNEXT_HOSTNAME_API` with legacy env fallback compatibility, `/etc/dotnext/hostname.env`) and `install.sh` (`DOTNEXT_*` env aliases, `[dotnext-install]` trap, `"DOTNEXT PANEL — installed"` summary banner).
- Finalised the bundled WordPress cache and migrator plugin headers while preserving WP-CLI/agent slug compatibility, and verified `hostedsvc` and `support-claim` Go suites.

### Milestone 6 — dotNext UI/UX Redesign, Dark-First System & `/admin` / `/panel` Routes (2026-11-18)

- Designed the new **dotNext** dark-first SaaS visual identity (`panel-ui/src/muiTheme.ts`, `panel-ui/src/global.css`, `panel-ui/index.html`, `panel-ui/public/manifest.webmanifest`): obsidian canvas (`#0B0F17`), layered slate cards (`#111827` / `#182235`), hairline borders (`#1E293B`), electric cyan accent (`#06B6D4` dark / `#0891B2` light), emerald/amber/rose status tokens, tabular numerals (`.dn-tabular`), and accessible focus-visible outlines.
- Added custom vector brand marks (`panel-ui/public/images/dotnext_logo.svg`, `dotnext_logo_dark.svg`) and redesigned shell chrome (`DotNextTitle`, `DotNextFooter`, frosted glass `DotNextHeader` with `/` command-search hint badge and Admin/Client shell badge), split-panel SaaS `Login` showcase, and telemetry-styled `Dashboard` / `UserDashboard` KPI cards.
- Switched primary SPA routes to `/admin/*` and `/panel/*`, updated impersonation header/session storage (`X-DotNext-Act-As` / `dotnext_act_as` with the legacy session-key fallback), rebranded tab titles (`<host> | dotNext`), support links, and locale strings, and added `LegacyShellRedirect` (`/dotnext-admin/*` → `/admin/*`, `/dotnext/*` → `/panel/*`) with unit tests.

### Milestones 4 & 5 — Modular Domain Core, CLI/Server Rebrand & Supervised Workers (2026-11-18)

- Added modular domain packages across Wave 1 & Wave 2 under `internal/modules/`:
  - `accounts`, `websites`, `dns`, `databases`, `system`, `mail`, `files`, `apps`, `backups`, `security` — with boundary validation, transaction helpers, and full unit tests (`internal/dbreserve` also extended to protect `dotnext`, `dotnext_*`, `dotnext-*`, and `dn_s_*` schemas/accounts).
- Wired `internal/platform/jobs.Supervisor` into `panel-api/cmd/server/serve.go` (wrapping 20 long-lived reconciler/ticker/sweeper/backup/event-source loops with panic recovery and coordinated shutdown drain) and mounted `GET /api/v1/admin/system/workers` via `internal/modules/system` in `panel-api/internal/app/app.go` + `panel-api/internal/api/openapi.yaml`.
- Rebranded the primary CLI and server metadata in `panel-api/cmd/server/{main,root,serve,version_cmd,listener,tenant_bootstrap}.go` and `panel-api/internal/api/{branding.go,openapi.yaml}` (`Use: "dotnext"`, `DOTNEXT_CONFIG`/`DOTNEXT_BOOTSTRAP_*`/`DOTNEXT_REDIS_PANEL_TOKEN` with legacy env fallbacks, `dotnext-sockets` with the legacy socket-group fallback, `BrandingDir = "/var/lib/dotnext-panel/branding"`), and regenerated `docs/site/platform/cli-reference.md` and `README.md` CLI references.
- Fixed remaining macOS test portability issues (`t.TempDir()` Unix socket path length in `panel-api/internal/agent/client_test.go`, `cmd/server/listener_test.go`, `cmd/server/update_gh727_test.go`, `internal/app/kratos_proxy_test.go`; `/var` symlink resolution in `webmailbranding_test.go` and `restore_cron_integration_test.go`).

### Milestone 3 — Agent Rebranding, Shared `pkg/peercred` & Binary Hygiene (2026-11-18)

- Extracted `pkg/peercred` (`peercred_linux.go` using `SO_PEERCRED`, `peercred_darwin.go` using `LOCAL_PEERCRED`, `peercred_other.go` fail-closed) to eliminate 5 duplicated Unix socket credential helpers across `panel-api` and `panel-agent`.
- Settled the binary names on `cmd/dotnext-*`, `cmd/dotnext-installer`, and `hostedsvc/cmd/dotnexthosted-svc`; updated `Makefile`, `tools/build-release.sh`, `panel-api/cmd/server/update.go`, `install.sh`, and `.gitignore`.
- Removed an accidentally tracked ELF binary from git.
- Updated `dotnext-agent` flags/env (`DOTNEXT_AGENT_*` with legacy env fallbacks), `pkg/paths.AgentSocket`, and log/service identifiers across `dotnext-agent`, `dotnext-mailhook`, and `dotnext-sendmail`.
- Added `panel-agent/internal/commands/catalog_test.go` verifying the 496-command RPC registry contract (`<domain>.<verb>`).

### Milestone 2 — Platform Core, Auth/Wire Rebrand & Cross-OS Middleware (2026-11-18)

- Added `internal/platform/` core packages:
  - `internal/platform/logger`: structured `slog` JSON/text factory and context propagation.
  - `internal/platform/database`: transaction-safe `WithTx` + context propagation (`FromContext`).
  - `internal/platform/httpserver`: standardized `APIError` envelope, `ListEnvelope[T]`, and bounded `ParsePagination`.
  - `internal/platform/jobs`: supervised background worker manager with panic recovery, restart tracking, and graceful shutdown.
- Rebranded config defaults (`pkg/paths`), env overrides (`DOTNEXT_*` with legacy env fallbacks), API tokens (`dnt_` primary + its legacy-prefix compat), Automation HMAC (`DotNext-HMAC` primary + the legacy HMAC header), Impersonation header (`X-DotNext-Act-As` + `X-Legacy-Act-As`), Webhook signature (`X-DotNext-Signature`), and Redis streams (`dotnext:audit:queue`, `dotnext:notifications:*`).
- Split `SO_PEERCRED` (`peercred_linux.go` / `peercred_other.go`) and cPanel `openat2` shadow reader (`mail_shadow_linux.go` / `mail_shadow_other.go`) so `panel-api` tests run natively with `-race` on macOS while keeping Linux kernel containment unchanged.

### Milestone 1 — Foundation, Module Rebrand & Cross-OS Hygiene (2026-11-18)

- First public release of the dotNext control panel (MIT licensed), consolidating the earlier private development milestones listed below.
- Go module path is `github.com/dotnext-panel/dotnext`, used across `go.mod`, `support-claim/go.mod`, Go imports, build scripts, and CI workflows.
- Added `pkg/paths` as the centralized source of truth for dotNext binaries, Unix sockets, config paths, systemd service names, databases, and auth/wire prefixes (alongside legacy constants for safe migration).
- Updated `internal/filesafe` (`DefaultAdminDeniedPrefixes`) to protect `/etc/dotnext`, `/run/dotnext*`, `/var/lib/dotnext*`, and `/usr/local/bin/dotnext*` alongside the legacy panel paths the upgrader rewrites.
- Added `//go:build linux` and `!linux` stubs in `internal/filesafe` and `internal/fsperm`, and shortened Unix socket test paths in `internal/kratosclient` and `internal/filesafe/filesafe_test.go`, enabling native macOS test execution without weakening Linux `openat2`/`RESOLVE_BENEATH` containment.
- Updated `Makefile` binary targets to `bin/dotnext`, `bin/dotnext-agent`, `bin/dotnext-installer`, `bin/dotnext-webdav`, `bin/dotnext-ssh-shell`, `bin/dotnext-mailhook`, and `bin/dotnext-sendmail`.

---

## Upstream History (DotNext)

Each entry is one move of the `stable` tag. Servers on the stable channel
install the commit `stable` points at, at their next automatic update (daily,
04:30). `stable` moves once a month, on the 1st, and in between only for a
hotfix ([runbook](docs/runbooks/hotfix-release.md)). Servers on the
development channel follow main and get each change as it merges.

An entry lists every change merged since the previous entry, oldest first,
and the migrations it adds. "Changing rows" names the migrations that update,
insert or delete existing data, not only the schema. The newest entry is
first.

Moves of `stable` before October 2026 did not follow the monthly schedule and
are not listed here. The commit history and the GitHub releases
(`release-<commit>`) have them.

## Adding an entry

Right after `stable` moves, generate the entry from the commit the previous
entry names, put it above the newest entry, and merge it to main through a
pull request:

```
git fetch origin --force main '+refs/tags/stable:refs/tags/stable'
tools/changelog-entry.sh <previous entry's commit> stable [--hotfix]
```

Add a line by hand when servers do something on their first update with it:
a data migration, a service restart, a new alert. The entry reaches stable
servers with the next move.

## Release log

Per-release entries (commit lists, migrations, "changing rows") are written
from here on by `tools/changelog-entry.sh` when `stable` moves. The milestones
above record the work that built the initial public release.
