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
- **Ubuntu 22.04 LTS / 24.04 LTS / 26.04 LTS**: Supported. Ubuntu is a Debian derivative, so it shares the Debian code path (`apt`, Sury-layout PHP in `/etc/php/<x.y>`, `ufw`, AppArmor, `redis-server`, `ssh.service`).
  - **PHP source**: the `ondrej/php` Launchpad PPA is preferred (no Fastly in front of it). The PPA only publishes suites for releases Ondřej Surý has built for — **26.04 `resolute` has no PPA suite yet**, so the installer probes `dists/<codename>/Release` and falls back to `packages.sury.org`, which does publish `resolute`. Without that probe every `apt-get update` on 26.04 fails with `404 … resolute Release`.
  - **PHP versions**: Ubuntu 24.04 ships 8.3 and has no archive 8.4; 26.04 ships 8.5. The panel's default stays **PHP 8.4** because phpMyAdmin 5.2.3 requires `>=7.2,<8.4` — on 26.04 that comes from Sury/PPA, not the archive.
  - **AppArmor**: 24.04+ enforces `kernel.apparmor_restrict_unprivileged_userns=1`, which blocks the bwrap SSH sandbox. `install.sh` loads a scoped `flags=(unconfined)` profile for `dotnext-ssh-shell` when that sysctl is on; it applies on 26.04 (AppArmor 5.0, `abi/4.0` is forward-compatible).
  - **systemd**: 24.04 = 255, 26.04 = 259, so every unit directive used here is supported.
  - **pdns-recursor**: 4.9 on 24.04, 5.3.5 on 26.04. The `--enable-old-settings` ExecStart drop-in is added from the *installed binary's* version (5.2+), not the distro, so both work unchanged.
  - **Third-party suites**: Docker (`download.docker.com/linux/ubuntu/<codename>`) publishes `resolute`; NodeSource uses the distro-agnostic `nodistro` suite and CrowdSec uses packagecloud `any`, so neither is release-coupled.
  - **Interim releases are not supported**: 25.10 (questing) reached end of life in July 2026 and neither the PPA nor Sury publishes a `questing` suite, so PHP cannot be sourced. `install.sh` refuses it with that reason instead of failing later in the apt run.

- **AlmaLinux 8 / 9 / 10** and **Rocky Linux 8 / 9 / 10**: Supported across `install.sh`, `deploy/installer/lib/packages.sh`, `internal/platform/osinfo`, `internal/phpext`, `agent/`, and `internal/` (`dnf`/`rpm`, `dnf-automatic`, Remi SCL multi-version PHP under `/etc/opt/remi/php<n>`, `firewalld`, SELinux booleans + port contexts, `nftables` (per-tenant egress firewall), `sshd.service`, `/etc/pdns` + `/etc/pdns-recursor`, `/etc/my.cnf.d`, `redis` on EL8/9 and `valkey` on EL10).
  - **EL8**: `redis:6` module is enabled where available, because per-tenant cache isolation uses `ACL SETUSER` (Redis 6+). On an EL8.5-or-older image without that module the panel still installs, but tenant-scoped Redis ACLs are unavailable.
  - **EL10**: `pdns`/`pdns-recursor` come from EPEL10 (5.0.x / 5.4.x). The `repo.powerdns.com` fallback is pinned to the `auth-50`/`rec-53` series there — repo.powerdns.com does not publish `auth-49`/`rec-51` for EL10.
  - **AppArmor profiles are Debian/Ubuntu only**; on EL the same hardening is SELinux, and the Security tab reports AppArmor as disabled with the reason.
- **Bash Version**: `install.sh` requires Bash 4+ (Debian/Ubuntu/AlmaLinux ship Bash 4.4/5.x). macOS ships GNU Bash 3.2.57, which rejects Bash 4+ syntax on both upstream and dotNext `install.sh`; modular helpers in `deploy/installer/lib/*.sh` and `infrastructure/systemd/*` pass syntax checking under Bash 3.2+.
