# CLI

The `dotnext` binary is one Cobra root command with subcommands. Run it as root on the panel host. Most commands read the panel database and call the agent over `/run/dotnext/agent.sock`.

This page groups the commands by area and shows the common tasks. The [CLI reference](./cli-reference.md) lists every command and every flag. It is generated from the command tree, and CI fails when it drifts. `dotnext-cli <command> --help` gives the same detail on a box.

## Global flags

- `--json`: print JSON instead of text, where the command supports it.
- `--config <path>`: the config file (default `/etc/dotnext/config.toml`).

## Command groups

### Panel and updates

| Command | What it does |
|---|---|
| `dotnext-cli serve` | Start the dotNext HTTP(S) server |
| `dotnext-cli update` | Pull latest code, rebuild, migrate, and restart services |
| `dotnext-cli repair` | Detect and fix known deployment-host issues |
| `dotnext-cli version` | Print DotNext version, commit, and runtime info |
| `dotnext-cli release` | Release-channel management (stable/development) |
| `dotnext-cli migrate` | Database migration commands |
| `dotnext-cli retention-sweep` | Prune expired rows from unbounded log/report tables |
| `dotnext-cli active-tasks` | Show running/queued tasks (updates, backups, malware scans) |
| `dotnext-cli secrets` | Rotate panel secrets after remediation (operator ceremony, run as root) |
| `dotnext-cli settings` | Inspect and patch server settings (headless equivalent of /admin/settings) |

### Users and access

| Command | What it does |
|---|---|
| `dotnext-cli user` | Manage panel users |
| `dotnext-cli user-token` | Manage a user's API tokens (mint / list / revoke) |
| `dotnext-cli session` | Audit and revoke active login sessions (Kratos) |
| `dotnext-cli impersonation` | Manage admin act-as impersonation grants (list / create / end) |
| `dotnext-cli automation-token` | Manage Automation API tokens (headless provisioning) |
| `dotnext-cli sso` | SSO (Single Sign-On) management commands |
| `dotnext-cli sso-reap` | Sweep stranded dotnext-sso-<nonce>.php files |
| `dotnext-cli admin` | Operator-only administrative subcommands |

### Domains and web

| Command | What it does |
|---|---|
| `dotnext-cli domain` | Manage hosted domains |
| `dotnext-cli ssl` | Manage Let's Encrypt SSL certificates |
| `dotnext-cli panel-cert` | Manage the panel's own TLS certificates (hostname + mail) |
| `dotnext-cli page-template` | Manage error/index page templates |
| `dotnext-cli php` | PHP version + extension + per-user pool management |
| `dotnext-cli php-defense` | PHP-defense (Snuffleupagus) status / mode / rules / incidents |
| `dotnext-cli package` | Manage hosting packages |
| `dotnext-cli limits` | Per-user resource limits (cgroups v2 + POSIX quota + nginx) |
| `dotnext-cli files` | Scoped tenant file manager (list/read/mkdir/move/chmod/archive/…) — same policy as the GUI |
| `dotnext-cli disk-usage` | Inspect / refresh a tenant's disk-usage breakdown (files / email / databases) |
| `dotnext-cli log` | Inspect log types and mint/revoke log-stream access grants |

### DNS

| Command | What it does |
|---|---|
| `dotnext-cli dns` | DNS zones and records (dns_zones / dns_records): list / add / update / delete |
| `dotnext-cli pdns` | PowerDNS helpers (recursor forwarders, etc.) |

### Mail

| Command | What it does |
|---|---|
| `dotnext-cli mailbox` | Manage mailboxes (email via Stalwart) |
| `dotnext-cli mail` | Admin mail queue + outbound throttle management |
| `dotnext-cli mail-cert` | Manage a domain's mail TLS certificate (mail.<domain> SAN) |
| `dotnext-cli mail-group` | Manage mail groups (distribution lists + shared-resource groups) and members |
| `dotnext-cli shared-resource` | Manage shared mail resources — calendars, contacts, files |
| `dotnext-cli panel-primary` | Manage the panel's primary mail domain row |

### Databases, cron and SSH

| Command | What it does |
|---|---|
| `dotnext-cli db` | Manage user databases (mariadb / postgres) |
| `dotnext-cli cron` | Manage user cron jobs (systemd-user timers) |
| `dotnext-cli ssh-key` | Manage user SSH authorized keys |
| `dotnext-cli nspawn` | Manage SSH sandbox nspawn images |

### Apps

| Command | What it does |
|---|---|
| `dotnext-cli app` | Manage one-click app installs (direct DB) |
| `dotnext-cli docker` | Manage the docker engine + app-marketplace host |
| `dotnext-cli docker-app` | Manage Docker app catalog installs (admin-only) |
| `dotnext-cli python-app` | Manage Python apps (admin-only) |

### Backups and recovery

| Command | What it does |
|---|---|
| `dotnext-cli backup` | Backup & restore subcommands (restic-backed) |
| `dotnext-cli dr` | Disaster-recovery standby: pair, status, promote |

### Security

| Command | What it does |
|---|---|
| `dotnext-cli crowdsec` | CrowdSec + AppSec operations (decisions, allowlists, hub, alerts, geoblock, captcha) |
| `dotnext-cli appsec` | CrowdSec AppSec config operator subcommands |
| `dotnext-cli apparmor` | AppArmor profile management operator commands |
| `dotnext-cli aide` | AIDE file integrity monitor operator commands |
| `dotnext-cli malware` | Malware scan / quarantine / YARA / settings (incident response) |
| `dotnext-cli malware-purge` | Hard-delete terminated malware quarantine rows past retention |
| `dotnext-cli ufw` | UFW utilities (port baseline only; IP decisions live in CrowdSec) |
| `dotnext-cli per-user-egress` | Per-user PHP-FPM egress firewall operator commands |
| `dotnext-cli audit` | Query, verify, and prune the unified audit log |

### Server

| Command | What it does |
|---|---|
| `dotnext-cli system` | System information and services |
| `dotnext-cli service` | Admin service control (start/stop/restart/reload/enable/disable) |
| `dotnext-cli ip` | Admin IP address pool (managed_ips): list / add / update / delete |
| `dotnext-cli notification` | Inspect notification channels and toggle event notifications |

## Common tasks

### Get back into the panel

The installer prints the admin username and password at the end. If you lost them, set a new password, or print a one-click recovery link:

```bash
dotnext-cli user password <email|username|user-id>          # sets and prints a new password
dotnext-cli user password <email|username|user-id> --link   # prints a recovery URL (valid 24 h)
dotnext-cli user 2fa-reset <email|username|user-id>         # removes TOTP and recovery codes
```

### Users and domains

```bash
dotnext-cli user list
dotnext-cli user create --username alice --email alice@example.com --password-stdin
dotnext-cli user suspend <id> --reason "unpaid invoice"
dotnext-cli user unsuspend <id>
dotnext-cli domain list
dotnext-cli domain create --user alice --name example.com
dotnext-cli ssl enable example.com
dotnext-cli ssl renew example.com --force
```

### Mail and databases

```bash
dotnext-cli mailbox create --domain example.com --local info    # prints the password if it generated one
dotnext-cli mailbox passwd info@example.com
dotnext-cli domain email-dkim-rotate example.com
dotnext-cli db create --user alice --name shop
dotnext-cli db user create --user alice --name shop
dotnext-cli db root-password                                     # rotates the MariaDB root password
```

### Backups

```bash
dotnext-cli backup destination list
dotnext-cli backup destination test <id-or-name>                 # creates the restic repo if it is missing
dotnext-cli backup schedule list
dotnext-cli backup schedule run-now <id>
dotnext-cli backup account-restore --user alice --snapshot <id> --destination <name> --force
```

### Updates, repair and support

```bash
dotnext-cli update                   # download the release, migrate and restart
dotnext-cli repair --diagnose        # list known host problems, change nothing
dotnext-cli repair --auto            # fix every non-destructive problem
dotnext-cli service status           # the units the panel watches
dotnext-cli system diagnostic        # upload an encrypted support bundle; prints a link and a password
dotnext-cli audit query --limit 100
dotnext-cli audit verify             # check the audit log's hash chain
```

## Exit status

A command exits `0` on success. On any error it prints the error on stderr and exits `1`.
