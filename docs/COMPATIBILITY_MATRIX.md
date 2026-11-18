# dotNext Feature Compatibility & Migration Matrix

**dotNext** is an independently structured, dark-first multi-tenant web hosting control panel licensed under **MIT** (see [LICENSE](../LICENSE), [THIRD_PARTY_NOTICES](../THIRD_PARTY_NOTICES), and [LICENSING_AUDIT.md](LICENSING_AUDIT.md)), with companion WordPress plugins in `wp-plugins/` retained under `GPL-2.0-or-later`.

## 1. Architecture & Package Organization

| Layer | Location | Responsibility |
|---|---|---|
| OS abstraction | `internal/platform/osinfo/` | Unified Linux distro detection (`debian`/`ubuntu`/`almalinux`), package manager (`apt`/`dnf`), PHP layout (Debian/Sury `/etc/php` vs. Remi SCL `/etc/opt/remi`), Redis/Valkey, SSH, PowerDNS, firewall (`ufw`/`firewalld`), and MAC (`apparmor`/`selinux`) helpers |
| Canonical identifiers | `pkg/paths/` | Single source of truth for dotNext binaries, sockets, configs, systemd units, DB names, headers, and legacy compatibility constants |
| Cross-OS peer auth | `pkg/peercred/` | Unified `SO_PEERCRED` (Linux), `LOCAL_PEERCRED` (macOS), fail-closed stub (`!linux && !darwin`) replacing 5 duplicate helpers |
| Shared platform runtime | `internal/platform/{logger,database,httpserver,jobs}` | Structured `slog` logging, `WithTx` transaction helper, `APIError`/`ListEnvelope`/`ParsePagination`, supervised background worker manager (`Supervisor`) |
| Domain boundary modules | `internal/modules/{accounts,websites,dns,databases,system,mail,files,apps,backups,security}` | Normalized input validation, reserved-name guards, docroot/admin `filesafe` confinement, and worker health endpoint (`GET /api/v1/admin/system/workers`) |
| Control-plane API & CLI | `panel-api/` | Gin REST API (`/api/v1`), Kratos cookie + `dnt_`/`jat_` bearer + `DotNext-HMAC` auth, Cobra CLI (`dotnext`) |
| Privileged host agent | `panel-agent/` | Root-owned Unix-socket NDJSON-RPC daemon (`dotnext-agent`) with 496 cataloged commands, strict `openat2` confinement, WebDAV, SSH sandbox, mail shims |
| Web UI SPA | `frontend/` | React 19 + Ant Design 6 + TanStack Query dark-first SaaS interface served at `/admin/*` and `/panel/*` (with `/dotnext-admin/*` & `/dotnext-api/*` redirects) |
| Installer & TUI | `installer/`, `install.sh`, `deploy/installer/lib/` | Bubble Tea installer (`dotnext-installer`) + modular preflight/services/rollback/uninstall helpers |
| Integrations | `hostedsvc/`, `support-claim/`, `wp-plugins/` | Rebranded service daemons (`dotnexthosted-svc`) and WordPress plugins (`dotNext Cache`, `dotNext Migrator`) |

## 2. Backward Compatibility & Wire Contracts

All security boundaries and wire protocols remain preserved during incremental migration:

- **API Bearer tokens**: `dnt_` primary prefix (`jat_` legacy tokens still accepted and validated identically).
- **Automation HMAC**: `DotNext-HMAC` primary header (`DotNext-HMAC` accepted as fallback).
- **Impersonation**: `X-DotNext-Act-As` primary header & `dotnext_act_as` session key (`X-DotNext-Act-As` / `dotnext_act_as` fallback).
- **Webhook signatures**: `X-DotNext-Signature` (paired with `X-DotNext-Signature` verification compatibility).
- **UI URLs**: `/admin/*` and `/panel/*` primary routes; `/dotnext-admin/*` and `/dotnext-api/*` automatically redirect preserving path, query, and hash.
- **Configuration & Env**: `DOTNEXT_CONFIG`, `DOTNEXT_BOOTSTRAP_*`, `DOTNEXT_REDIS_PANEL_TOKEN`, `DOTNEXT_MODULES`, `DOTNEXT_HOSTNAME` with `DOTNEXT_*` fallbacks and `/etc/dotnext/*` → `/etc/dotnext*` path fallbacks.

## 3. Feature Preservation Matrix

| Capability | Backend Implementation | Verification Scope (Observed) |
|---|---|---|
| Multi-tenant accounts & packages | `internal/userops`, `packageops`, `internal/modules/accounts`, `internal/limits` | Verified by Go unit/integration tests on host OS |
| Websites, Nginx & PHP-FPM | `internal/api` + `reconciler` + `agent/commands` (`496` RPC commands) | API/reconciler Go tests pass; Linux agent `commands` & `server` test binaries cross-compile cleanly |
| DNS & DNSSEC | `internal/api`, `dnsverify`, `internal/modules/dns`, `agent/pdns` | Go `-race` tests pass (`pdns`, `dnsverify`, `modules/dns`) |
| SSL & Let's Encrypt | `internal/panelcert`, `reconciler/ssl_*`, `agent/certbot` | Go `-race` tests pass (`panelcert`, `certbot`, `api`) |
| MariaDB & PostgreSQL | `internal/dbreserve`, `internal/dbshadow`, `internal/modules/databases`, `dbconsoleops` | Go `-race` tests pass |
| Email, Stalwart & Webmail | `internal/mailboxops`, `stalwartadmin`, `mailscan`, `internal/modules/mail` | Go `-race` tests pass (`mailhook`, `mailrelay`, `sendmailshim`, `modules/mail`) |
| File Manager & WebDAV | `internal/filesafe`, `internal/fsperm`, `internal/modules/files`, `dotnext-webdav` | Linux `openat2`/`RESOLVE_BENEATH` preserved; macOS tests + Linux build pass |
| WordPress, Docker & Python Apps | `internal/apps`, `dockerapp`, `pyframeworks`, `internal/modules/apps`, `wp-plugins/` | Go `-race` tests pass |
| Backups, Schedules & DR Sync | `internal/backup`, `backupscheduler`, `backupfinalizer`, `drsync`, `internal/modules/backups` | Go `-race` tests pass |
| Migrations (cPanel, Plesk, etc.) | `internal/migrate/*` | Go `-race` tests pass across all 7 migration importers |
| Security, Firewall & Audit | `internal/appseccfg`, `internal/audit`, `notifications`, `internal/modules/security` | Go `-race` tests + `make demo-guard` pass |
| Supervised Workers & Telemetry | `internal/platform/jobs`, `internal/modules/system` (`GET /api/v1/admin/system/workers`), `openapi.yaml` | Go `-race` tests + OpenAPI coverage parity tests pass |

## 4. Known Environment Boundaries

- **Live VM / Systemd Execution**: Not available on this macOS (`darwin/arm64`) development host. Linux agent privileged execution, systemd unit activation, and end-to-end `install.sh` VM runs require a target Debian 12/13, Ubuntu 22.04/24.04, or AlmaLinux 8/9/10 Linux host.
- **AlmaLinux 8 / 9 / 10**: Supported across `install.sh`, `deploy/installer/lib/packages.sh`, `internal/platform/osinfo`, `internal/phpext`, `panel-agent`, and `panel-api` (`dnf`/`rpm`, `dnf-automatic`, Remi SCL multi-version PHP, `firewalld`, SELinux configuration, `sshd.service`, `/etc/pdns` + `/etc/pdns-recursor`, `/etc/my.cnf.d`, `redis` on EL8/9 and `valkey` on EL10).
- **Bash Version**: `install.sh` requires Bash 4+ (Debian/Ubuntu/AlmaLinux ship Bash 4.4/5.x). macOS ships GNU Bash 3.2.57, which rejects Bash 4+ syntax on both upstream and dotNext `install.sh`; modular helpers in `deploy/installer/lib/*.sh` and `infrastructure/systemd/*` pass syntax checking under Bash 3.2+.
