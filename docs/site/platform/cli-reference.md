# dotNext CLI reference

> Generated from the Cobra command tree — do not edit by hand.
> Regenerate with `go test ./internal/panelcmd -run TestCLIReferenceGolden -update`.

## `dotnext-cli`

dotNext Panel — web hosting control panel

```
dotnext-cli
```

**Flags:**

- `--config` — config file path (default: /etc/dotnext/config.toml)
- `--json` — output as JSON

### `dotnext-cli active-tasks`

Show running/queued tasks (updates, backups, malware scans)

```
dotnext-cli active-tasks
```

### `dotnext-cli admin`

Operator-only administrative subcommands

```
dotnext-cli admin
```

#### `dotnext-cli admin backfill-usernames`

Assign a unique username to every user missing one

```
dotnext-cli admin backfill-usernames [flags]
```

**Flags:**

- `--apply` — Write the derived usernames to the DB

#### `dotnext-cli admin purge-orphan-identities`

Delete Kratos identities with no matching panel user (frees stuck usernames)

```
dotnext-cli admin purge-orphan-identities [flags]
```

**Flags:**

- `--apply` — delete the orphan identities (default: dry-run)
- `--username` — only purge the orphan whose username or email matches this value

#### `dotnext-cli admin rebuild-kratos`

Recreate Kratos identities from panel users (DB-loss recovery)

```
dotnext-cli admin rebuild-kratos [flags]
```

**Flags:**

- `--dry-run` — List target users and exit without writing
- `--expires-in` — Kratos recovery-code TTL (e.g. 1h, 24h) — operators typically need ≥24h to distribute (default `24h`)
- `--output` — CSV file to emit temp passwords (default: /var/lib/dotnext/recovery/, root-only). Contains PLAINTEXT temporary passwords — delete after distribution.
- `--yes` — Skip interactive confirmation prompt

#### `dotnext-cli admin relabel-identifiers`

Re-key Kratos login identifiers email -> username

```
dotnext-cli admin relabel-identifiers [flags]
```

**Flags:**

- `--apply` — PATCH the identities (default: dry-run)

#### `dotnext-cli admin slice-cutover`

Migrate FPM from global master to per-user masters and mask distro units

```
dotnext-cli admin slice-cutover [flags]
```

**Flags:**

- `--dry-run` — run preflight only; do not mask global FPM or probe

### `dotnext-cli aide`

AIDE file integrity monitor operator commands

```
dotnext-cli aide
```

#### `dotnext-cli aide check`

Trigger an on-demand AIDE integrity check

```
dotnext-cli aide check
```

#### `dotnext-cli aide rebuild`

Re-baseline the AIDE database after a deliberate change

```
dotnext-cli aide rebuild [flags]
```

**Flags:**

- `--dry-run` — Preview plan without rebuilding (mutually exclusive with --full)
- `--full` — Confirm full DB re-init
- `--paths` — Partial re-baseline: only promote changes matching this regex (e.g. '^/usr/local/bin/dotnext-'); refuses if changes outside the regex are detected

#### `dotnext-cli aide status`

Print AIDE DB age + last check summary

```
dotnext-cli aide status
```

### `dotnext-cli app`

Manage one-click app installs (direct DB)

```
dotnext-cli app
```

#### `dotnext-cli app cache`

Enable/disable object + page cache for a WordPress app install

```
dotnext-cli app cache <install-id> [flags]
```

**Flags:**

- `--disable` — disable object + page cache
- `--enable` — enable object + page cache

#### `dotnext-cli app cache-doctor`

Detect (and optionally --repair) drift on cache-enabled WordPress installs

```
dotnext-cli app cache-doctor [flags]
```

**Flags:**

- `--json` — emit a JSON report
- `--migrate-acl` — re-provision every cache-enabled install to a per-install Redis ACL user, then reap orphaned legacy shared users (safe/idempotent)
- `--repair` — re-provision drifted installs (re-stamp constants, re-provision ACL, re-verify)

#### `dotnext-cli app cache-doctor-run`

Run a recorded fleet cache-doctor sweep (persists a run for the admin GUI history)

```
dotnext-cli app cache-doctor-run [flags]
```

**Flags:**

- `--kind` — sweep kind: doctor | repair | refresh (default `doctor`)

#### `dotnext-cli app clone`

Clone a WordPress app install onto another domain (files + DB)

```
dotnext-cli app clone <src-install-id> [flags]
```

**Flags:**

- `--to-domain` — destination domain (ID or name)

#### `dotnext-cli app delete`

Delete an installed app (direct DB + agent teardown)

```
dotnext-cli app delete <install-id|domain-name> [flags]
```

**Flags:**

- `--force` — Skip confirmation prompt

#### `dotnext-cli app e2e`

Install every app on a domain, report pass/fail, then delete

```
dotnext-cli app e2e [flags]
```

**Flags:**

- `--base-subdir` — Subdir prefix; each app installs under <prefix>_<app>_<rand> (default `e2e`)
- `--domain` — Domain name or ULID to install all apps under (required)
- `--keep` — Don't delete installs after the run (debug)
- `--only` — Only run these app_types (comma-separated) (default `[]`)
- `--skip` — Skip these app_types (comma-separated) (default `[]`)
- `--stop-on-fail` — Stop the sweep after the first failure
- `--wait-timeout` — Per-app install timeout in seconds (default `600`)

#### `dotnext-cli app get`

Show one installed app (direct DB)

```
dotnext-cli app get <install-id|domain-name>
```

#### `dotnext-cli app install`

Install an app on a domain (direct service)

```
dotnext-cli app install [flags]
```

**Flags:**

- `--app-type` — App descriptor name (see `dotnext app registry`)
- `--directory` — Subdirectory under docroot (empty = site root)
- `--domain` — Target domain name or ULID (e.g. example.com or 01KPR…)
- `--param` — Per-app param: --param key=value (value is JSON; repeat for multiple) (default `[]`)
- `--use-www` — Reachable at www.<domain> too
- `--user` — Owner user (email, username, or ULID; default: domain owner)
- `--wait` — Poll until status is ready or failed
- `--wait-timeout` — Seconds to wait when --wait is set (default `600`)

#### `dotnext-cli app list`

List installed apps (direct DB)

```
dotnext-cli app list
```

#### `dotnext-cli app magic-link`

Mint a one-time magic-link login for an installed app (sensitive, ~60s)

```
dotnext-cli app magic-link <install-id>
```

#### `dotnext-cli app refresh-cache-plugin`

Re-install the bundled dotnext-cache plugin on every cache-enabled WordPress site

```
dotnext-cli app refresh-cache-plugin
```

#### `dotnext-cli app registry`

List available app types and their parameter schemas (direct read — no DB)

```
dotnext-cli app registry
```

#### `dotnext-cli app scan`

Scan a user's homedir for unregistered apps and register them

```
dotnext-cli app scan [flags]
```

**Flags:**

- `--admin-email` — admin email for discovered installs (default: user's email)
- `--user` — user (id or username) whose homedir to scan (required)

### `dotnext-cli apparmor`

AppArmor profile management operator commands

```
dotnext-cli apparmor
```

#### `dotnext-cli apparmor flip-mature`

Flip soak-clean complain-mode profiles to enforce

```
dotnext-cli apparmor flip-mature [flags]
```

**Flags:**

- `--dry-run` — show what would change without invoking aa-enforce
- `--force` — flip even if recent violations exist (DANGEROUS: a profile forced to enforce despite violations can crash-loop its service)
- `--profile` — Flip a single profile only
- `--soak-days` — violation-scan window in days; a profile with any AppArmor DENIED or complain-mode ALLOWED event in this window is not flipped (default `7`)

#### `dotnext-cli apparmor status`

List dotnext AppArmor profiles and modes

```
dotnext-cli apparmor status
```

### `dotnext-cli appsec`

CrowdSec AppSec config operator subcommands

```
dotnext-cli appsec
```

#### `dotnext-cli appsec exclusion`

Manage operator CRS false-positive exclusions

```
dotnext-cli appsec exclusion
```

##### `dotnext-cli appsec exclusion add`

Add an exclusion (host + path + rule are all required)

```
dotnext-cli appsec exclusion add [flags]
```

**Flags:**

- `--host` — hostname the exclusion applies to (required)
- `--note` — why this exclusion exists
- `--rule` — CRS rule ID or AppSec rule name to exclude, e.g. 942100 or crowdsecurity/vpatch-CVE-2025-29927 (required)
- `--uri-prefix` — URI prefix the exclusion applies to, e.g. /wp-json/x/ (required)

##### `dotnext-cli appsec exclusion list`

List operator CRS exclusions

```
dotnext-cli appsec exclusion list
```

##### `dotnext-cli appsec exclusion rm`

Remove an exclusion by id

```
dotnext-cli appsec exclusion rm <id>
```

#### `dotnext-cli appsec explain`

Show which CRS rules recently blocked requests (AppSec FP triage)

```
dotnext-cli appsec explain [flags]
```

**Flags:**

- `--json` — emit raw JSON instead of the grouped report
- `--limit` — how many recent AppSec alerts to inspect (max 200) (default `25`)

#### `dotnext-cli appsec flarum-sync`

Add the scoped CRS 920450 exclusion for every existing Flarum forum

```
dotnext-cli appsec flarum-sync
```

#### `dotnext-cli appsec host-mode`

Manage per-host AppSec mode (detection-only)

```
dotnext-cli appsec host-mode
```

##### `dotnext-cli appsec host-mode clear`

Clear a host's mode, restoring full blocking

```
dotnext-cli appsec host-mode clear [flags]
```

**Flags:**

- `--host` — hostname to clear (required)

##### `dotnext-cli appsec host-mode list`

List per-host AppSec modes

```
dotnext-cli appsec host-mode list
```

##### `dotnext-cli appsec host-mode set`

Set a host's AppSec mode (host is required; mode defaults to detect)

```
dotnext-cli appsec host-mode set [flags]
```

**Flags:**

- `--host` — hostname to set the mode for (required)
- `--mode` — AppSec mode (only 'detect' is supported) (default `detect`)
- `--note` — why this host is in this mode

#### `dotnext-cli appsec render-config`

Write /etc/crowdsec/appsec-configs/dotnext-appsec.yaml from internal/appseccfg.Render

```
dotnext-cli appsec render-config [flags]
```

**Flags:**

- `--reconcile` — preserve operator dotnext-mode/dotnext-countries header from existing file (default true) (default `true`)
- `--reload` — best-effort `systemctl reload crowdsec` after a real diff (default false; install.sh owns its own deferred reload, so leave off inside an update)

### `dotnext-cli audit`

Query, verify, and prune the unified audit log

```
dotnext-cli audit
```

#### `dotnext-cli audit prune`

Delete audit rows older than --days (retention)

```
dotnext-cli audit prune [flags]
```

**Flags:**

- `--days` — retention window in days (delete older) (default `365`)

#### `dotnext-cli audit query`

List recent audit events (admin/forensics view)

```
dotnext-cli audit query [flags]
```

**Flags:**

- `--limit` — max rows (1-1000) (default `50`)
- `--q` — search (action/target/actor_kind/result)

#### `dotnext-cli audit verify`

Recompute the hash chain and report tamper-evidence integrity

```
dotnext-cli audit verify
```

### `dotnext-cli automation-token`

Manage Automation API tokens (headless provisioning)

```
dotnext-cli automation-token
```

#### `dotnext-cli automation-token list`

List automation tokens (never reveals secrets)

```
dotnext-cli automation-token list
```

#### `dotnext-cli automation-token mint`

Mint an automation token; reveals the secret once

```
dotnext-cli automation-token mint <name> [flags]
```

**Flags:**

- `--scope` — grant a scope (repeatable), e.g. --scope read:status (default `[]`)

#### `dotnext-cli automation-token revoke`

Revoke an automation token

```
dotnext-cli automation-token revoke <id-or-name>
```

### `dotnext-cli backup`

Backup & restore subcommands (restic-backed)

```
dotnext-cli backup
```

#### `dotnext-cli backup account-restore`

Restore one account's backup snapshot via the agent (bypass UI)

```
dotnext-cli backup account-restore [flags]
```

**Flags:**

- `--apply` — apply home+db onto live system (false = staging-only smoke test) (default `true`)
- `--destination` — destination name (e.g. 'test', 'b2-prod')
- `--force` — required — restore overwrites home tree + reloads databases
- `--interactive` — force interactive prompts even when flags are set
- `--snapshot` — manifest snapshot id (long form preferred)
- `--target-user-id` — disaster recovery: panel row gone; use this ULID directly (pair with --target-username)
- `--target-username` — disaster recovery: system account name to chown into (pair with --target-user-id)
- `--user` — username (e.g. alice) — looks up panel row
- `--user-id` — user ULID — looks up panel row (alternative to --user)

#### `dotnext-cli backup copy`

[RETIRED] superseded by per-destination model

```
dotnext-cli backup copy
```

#### `dotnext-cli backup destination`

Manage backup destinations (local, sftp, s3, b2, azure, gcs, rest)

```
dotnext-cli backup destination
```

##### `dotnext-cli backup destination create`

Create a backup destination

```
dotnext-cli backup destination create [flags]
```

**Flags:**

- `--disabled` — create in disabled state
- `--env` — credential env: KEY=VALUE (repeatable) (default `[]`)
- `--env-stdin` — read additional KEY=VALUE lines from stdin
- `--kind` — destination kind: local|sftp|s3|b2|azure|gcs|rest (required)
- `--name` — destination name (required, unique)
- `--url` — restic repo URL (required)

##### `dotnext-cli backup destination delete`

Delete a backup destination

```
dotnext-cli backup destination delete <id-or-name> [flags]
```

**Flags:**

- `--force` — skip confirmation

##### `dotnext-cli backup destination get`

Show a backup destination

```
dotnext-cli backup destination get <id-or-name>
```

##### `dotnext-cli backup destination list`

List backup destinations

```
dotnext-cli backup destination list
```

##### `dotnext-cli backup destination rotate-password`

Rotate a backup destination's restic password (revealed once)

```
dotnext-cli backup destination rotate-password <id-or-name>
```

##### `dotnext-cli backup destination test`

Test connectivity (auto-inits restic repo if missing)

```
dotnext-cli backup destination test <id-or-name>
```

##### `dotnext-cli backup destination update`

Update a backup destination

```
dotnext-cli backup destination update <id-or-name> [flags]
```

**Flags:**

- `--clear-creds` — delete stored credential env for this destination
- `--disable` — mark destination disabled
- `--enable` — mark destination enabled
- `--env` — rewrite credential env: KEY=VALUE (repeatable; s3/b2/azure/gcs/rest/sftp secrets) (default `[]`)
- `--env-stdin` — read additional KEY=VALUE credential lines from stdin
- `--name` — new name
- `--sftp-auth` — sftp auth: 'key' or 'password'
- `--sftp-host` — sftp host (sftp kind)
- `--sftp-key-path` — absolute path to private key (auth=key)
- `--sftp-password` — sftp password (auth=password; stored as SSHPASS)
- `--sftp-path` — sftp remote path
- `--sftp-port` — sftp port (default 22) (default `0`)
- `--sftp-user` — sftp user
- `--url` — new restic repo URL (validated against existing kind)

#### `dotnext-cli backup retention`

Manage restic retention (forget + prune per destination)

```
dotnext-cli backup retention
```

##### `dotnext-cli backup retention apply`

Run restic forget per (schedule, destination) + prune per destination

```
dotnext-cli backup retention apply [flags]
```

**Flags:**

- `--dry-run` — List the backup jobs that would be forgotten and pass --dry-run to prune; no destructive ops

##### `dotnext-cli backup retention verify`

Check that every backup's stage snapshots are still in its repository (read-only)

```
dotnext-cli backup retention verify [flags]
```

**Flags:**

- `--destination` — check only this destination (id or name)

#### `dotnext-cli backup schedule`

Manage backup schedules

```
dotnext-cli backup schedule
```

##### `dotnext-cli backup schedule create`

Create a backup schedule

```
dotnext-cli backup schedule create [flags]
```

**Flags:**

- `--cron` — 5-field cron expression (e.g. '0 3 * * *')
- `--destination` — destination id or name (repeatable) (default `[]`)
- `--disabled` — create in disabled state
- `--include-system` — for account_backup: also fire system_backup each tick
- `--keep-daily` — restic forget --keep-daily (default `0`)
- `--keep-monthly` — restic forget --keep-monthly (default `0`)
- `--keep-weekly` — restic forget --keep-weekly (default `0`)
- `--keep-within-days` — keep every backup of the last N days (an account's "Days to keep") (default `0`)
- `--kind` — schedule kind: account_backup|system_backup (required)
- `--preset` — preset: daily|weekly|monthly (mutually exclusive with --cron)
- `--user` — user id|email|username for account_backup fan-out (repeatable; empty=all non-admins) (default `[]`)

##### `dotnext-cli backup schedule delete`

Delete a backup schedule

```
dotnext-cli backup schedule delete <id> [flags]
```

**Flags:**

- `--force` — skip confirmation

##### `dotnext-cli backup schedule get`

Show a backup schedule with destinations + users

```
dotnext-cli backup schedule get <id>
```

##### `dotnext-cli backup schedule list`

List backup schedules

```
dotnext-cli backup schedule list
```

##### `dotnext-cli backup schedule run-now`

Trigger a schedule by advancing next_run_at to now (scheduler picks up within ≤60s)

```
dotnext-cli backup schedule run-now <id>
```

##### `dotnext-cli backup schedule update`

Update a backup schedule

```
dotnext-cli backup schedule update <id> [flags]
```

**Flags:**

- `--clear-destinations` — remove all destinations
- `--clear-users` — remove all users (= fan-out to all)
- `--cron` — new cron expression
- `--destination` — replace destinations (repeatable) (default `[]`)
- `--disable` — mark schedule disabled
- `--enable` — mark schedule enabled
- `--include-system` — true|false (account_backup only)
- `--keep-daily` —  (default `0`)
- `--keep-monthly` —  (default `0`)
- `--keep-weekly` —  (default `0`)
- `--keep-within-days` — keep every backup of the last N days (default `0`)
- `--preset` — preset: daily|weekly|monthly
- `--user` — replace users (repeatable) (default `[]`)

#### `dotnext-cli backup scheduler`

Backup scheduler ops (manual tick / debug)

```
dotnext-cli backup scheduler
```

##### `dotnext-cli backup scheduler tick`

Run one enqueue + dispatch pass of the backup scheduler synchronously

```
dotnext-cli backup scheduler tick
```

### `dotnext-cli cron`

Manage user cron jobs (systemd-user timers)

```
dotnext-cli cron
```

#### `dotnext-cli cron add`

Add a cron job (5-field cron, allowlisted commands only)

```
dotnext-cli cron add [flags]
```

**Flags:**

- `--command` — command to run (required, allowlisted)
- `--disabled` — create disabled
- `--name` — job name (required)
- `--schedule` — 5-field cron expression e.g. '*/15 * * * *' (required)
- `--user` — user (id|email|username) (required)

#### `dotnext-cli cron delete`

Delete a cron job (removes the systemd timer, then the row)

```
dotnext-cli cron delete <job-id> [flags]
```

**Flags:**

- `--force` — skip confirmation

#### `dotnext-cli cron http-trigger`

Fetch a self-domain URL with SSRF guard + rebind-safe IP pinning (internal cron exec helper)

```
dotnext-cli cron http-trigger <url>
```

#### `dotnext-cli cron list`

List cron jobs (filtered by user, or all)

```
dotnext-cli cron list [flags]
```

**Flags:**

- `--user` — filter by user (id|email|username); empty = all

#### `dotnext-cli cron run-now`

Run a cron job immediately; a run still going after 4 minutes keeps running

```
dotnext-cli cron run-now <job-id>
```

#### `dotnext-cli cron update`

Update a cron job

```
dotnext-cli cron update <job-id> [flags]
```

**Flags:**

- `--command` — 
- `--disable` — mark job disabled
- `--enable` — mark job enabled
- `--name` — 
- `--schedule` — 

### `dotnext-cli crowdsec`

CrowdSec + AppSec operations (decisions, allowlists, hub, alerts, geoblock, captcha)

```
dotnext-cli crowdsec
```

#### `dotnext-cli crowdsec alerts`

List / inspect alerts

```
dotnext-cli crowdsec alerts
```

##### `dotnext-cli crowdsec alerts inspect`

Inspect one alert

```
dotnext-cli crowdsec alerts inspect <id>
```

##### `dotnext-cli crowdsec alerts list`

List recent alerts

```
dotnext-cli crowdsec alerts list
```

#### `dotnext-cli crowdsec allowlists`

List / add / remove allowlist entries

```
dotnext-cli crowdsec allowlists
```

##### `dotnext-cli crowdsec allowlists add`

Add an allowlist entry (use --me to allowlist your current SSH IP)

```
dotnext-cli crowdsec allowlists add [flags]
```

**Flags:**

- `--me` — allowlist the IP you are SSH'd in from (roaming-admin recovery)
- `--reason` — reason label (default `manual (cli)`)
- `--value` — IP/CIDR to allowlist (required unless --me)

##### `dotnext-cli crowdsec allowlists list`

List allowlist entries

```
dotnext-cli crowdsec allowlists list
```

##### `dotnext-cli crowdsec allowlists remove`

Remove an allowlist entry

```
dotnext-cli crowdsec allowlists remove [flags]
```

**Flags:**

- `--value` — allowlist entry to remove (required)

#### `dotnext-cli crowdsec blocklists`

List / refresh blocklists

```
dotnext-cli crowdsec blocklists
```

##### `dotnext-cli crowdsec blocklists list`

List subscribed blocklists

```
dotnext-cli crowdsec blocklists list
```

##### `dotnext-cli crowdsec blocklists refresh`

Refresh blocklists now

```
dotnext-cli crowdsec blocklists refresh
```

#### `dotnext-cli crowdsec bouncers`

List CrowdSec bouncers

```
dotnext-cli crowdsec bouncers
```

#### `dotnext-cli crowdsec captcha`

Captcha remediation get / set

```
dotnext-cli crowdsec captcha
```

##### `dotnext-cli crowdsec captcha get`

Show captcha config (secret never shown)

```
dotnext-cli crowdsec captcha get
```

##### `dotnext-cli crowdsec captcha set`

Configure captcha remediation (merge: empty --secret-key keeps existing)

```
dotnext-cli crowdsec captcha set [flags]
```

**Flags:**

- `--enabled` — enable captcha remediation
- `--provider` — recaptcha|hcaptcha|turnstile
- `--secret-key` — captcha secret key (write-only; empty keeps existing)
- `--site-key` — captcha site key

#### `dotnext-cli crowdsec country-exempt`

Never block the selected countries (get / set / sync)

```
dotnext-cli crowdsec country-exempt
```

##### `dotnext-cli crowdsec country-exempt get`

Show exempt countries

```
dotnext-cli crowdsec country-exempt get
```

##### `dotnext-cli crowdsec country-exempt set`

Set exempt countries and/or extra CIDRs (both empty = feature off)

```
dotnext-cli crowdsec country-exempt set [flags]
```

**Flags:**

- `--country` — 2-letter ISO code (repeatable); omit all to disable (default `[]`)
- `--extra-cidr` — additional never-block IP or CIDR (repeatable); omit all to clear (default `[]`)

##### `dotnext-cli crowdsec country-exempt sync`

Force a CIDR re-sync (refetch zones, runs inline)

```
dotnext-cli crowdsec country-exempt sync
```

#### `dotnext-cli crowdsec decisions`

List / add / delete CrowdSec decisions

```
dotnext-cli crowdsec decisions
```

##### `dotnext-cli crowdsec decisions add`

Add a decision (ban)

```
dotnext-cli crowdsec decisions add [flags]
```

**Flags:**

- `--duration` — ban duration (e.g. 4h, 24h) (default `4h`)
- `--reason` — reason label (default `manual (cli)`)
- `--scope` — ip|range|country|as (default `ip`)
- `--value` — target value (required)

##### `dotnext-cli crowdsec decisions delete`

Delete a decision by id, or every decision for --ip

```
dotnext-cli crowdsec decisions delete [<id>] [flags]
```

**Flags:**

- `--ip` — delete all decisions targeting this IP/CIDR

##### `dotnext-cli crowdsec decisions list`

List active decisions

```
dotnext-cli crowdsec decisions list [flags]
```

**Flags:**

- `--limit` — max rows (1..1000) (default `0`)
- `--scope` — filter: ip|range|country|as

#### `dotnext-cli crowdsec geoblock`

AppSec geoblock get / set

```
dotnext-cli crowdsec geoblock
```

##### `dotnext-cli crowdsec geoblock get`

Show geoblock mode + countries

```
dotnext-cli crowdsec geoblock get
```

##### `dotnext-cli crowdsec geoblock set`

Set geoblock mode (off|allow|deny) and countries

```
dotnext-cli crowdsec geoblock set [flags]
```

**Flags:**

- `--country` — 2-letter ISO code (repeatable) (default `[]`)
- `--mode` — off|allow|deny (required)

#### `dotnext-cli crowdsec hub`

List / install / remove hub items

```
dotnext-cli crowdsec hub
```

##### `dotnext-cli crowdsec hub install`

Install a hub item

```
dotnext-cli crowdsec hub install [flags]
```

**Flags:**

- `--force` — force reinstall
- `--name` — hub item name (required)
- `--type` — collections|parsers|scenarios|appsec-rules (required)

##### `dotnext-cli crowdsec hub list`

List hub items (collections/parsers/scenarios/appsec-rules)

```
dotnext-cli crowdsec hub list [flags]
```

**Flags:**

- `--type` — filter by type

##### `dotnext-cli crowdsec hub remove`

Remove a hub item

```
dotnext-cli crowdsec hub remove [flags]
```

**Flags:**

- `--name` — item name (required)
- `--type` — item type (required)

#### `dotnext-cli crowdsec metrics`

Show CrowdSec metrics

```
dotnext-cli crowdsec metrics
```

#### `dotnext-cli crowdsec profiles`

Per-scenario remediation overrides get / set

```
dotnext-cli crowdsec profiles
```

##### `dotnext-cli crowdsec profiles get`

Show current profile overrides

```
dotnext-cli crowdsec profiles get
```

##### `dotnext-cli crowdsec profiles set`

Set per-scenario overrides as scenario:action (action=captcha|off)

```
dotnext-cli crowdsec profiles set [flags]
```

**Flags:**

- `--override` — scenario:action (repeatable; action=captcha|off) (default `[]`)

#### `dotnext-cli crowdsec status`

Show CrowdSec engine status

```
dotnext-cli crowdsec status
```

### `dotnext-cli db`

Manage user databases (mariadb / postgres)

```
dotnext-cli db
```

#### `dotnext-cli db backup`

Create a backup of a database (returns the dump path)

```
dotnext-cli db backup <db-id|db-name>
```

#### `dotnext-cli db chown`

Reassign a database to a different owner

```
dotnext-cli db chown <database> <new-owner> [flags]
```

**Flags:**

- `--yes` — skip the confirmation prompt

#### `dotnext-cli db config`

View / set database tuning (mariadb/postgres)

```
dotnext-cli db config
```

##### `dotnext-cli db config get`

Show tuning params + current values (allowlist)

```
dotnext-cli db config get [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)

##### `dotnext-cli db config set`

Set tuning params (allowlist-validated; applies + may restart the engine)

```
dotnext-cli db config set <param=value> [<param=value>...] [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)

#### `dotnext-cli db create`

Create a database for a user

```
dotnext-cli db create [flags]
```

**Flags:**

- `--as-admin` — Skip the username prefix (admin-only DB names)
- `--engine` — Engine: mariadb | postgres (default `mariadb`)
- `--name` — Database name (without user prefix) — required
- `--user` — User (email or username) — required

#### `dotnext-cli db delete`

Delete a database by ID

```
dotnext-cli db delete [flags]
```

**Flags:**

- `--id` — Database ID (ULID)

#### `dotnext-cli db kill`

Kill/terminate a database process

```
dotnext-cli db kill <process-id> [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)
- `--force` — confirm the kill

#### `dotnext-cli db list`

List databases (filtered by user, or all)

```
dotnext-cli db list [flags]
```

**Flags:**

- `--user` — Filter by user (email or username)

#### `dotnext-cli db maintenance`

Run optimize/analyze (mariadb) or vacuum/analyze (postgres)

```
dotnext-cli db maintenance [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)
- `--scope` — 'all' or a database name (default `all`)

#### `dotnext-cli db maintenance-status`

Show a maintenance job's status

```
dotnext-cli db maintenance-status <job-id>
```

#### `dotnext-cli db postgres`

PostgreSQL lifecycle: status / uninstall

```
dotnext-cli db postgres
```

##### `dotnext-cli db postgres status`

Show PostgreSQL installed/active/version

```
dotnext-cli db postgres status
```

##### `dotnext-cli db postgres uninstall`

DESTROY PostgreSQL: purge packages + delete all data dirs

```
dotnext-cli db postgres uninstall [flags]
```

**Flags:**

- `--force` — confirm permanent destruction of all PostgreSQL data

#### `dotnext-cli db processes`

List database processes/activity

```
dotnext-cli db processes [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)

#### `dotnext-cli db restore`

Restore a database from a .sql dump on the host

```
dotnext-cli db restore <db-id|db-name> [flags]
```

**Flags:**

- `--file` — path to a .sql dump on the host (required)
- `--force` — confirm the overwrite

#### `dotnext-cli db root-password`

Rotate the database root/superuser password (revealed once)

```
dotnext-cli db root-password [flags]
```

**Flags:**

- `--engine` — mariadb|postgres (default `mariadb`)

#### `dotnext-cli db sso`

Mint a single-use phpMyAdmin/Adminer SSO login URL for a database

```
dotnext-cli db sso --database <id> [--engine mariadb|postgres] [flags]
```

**Flags:**

- `--database` — database id (ULID)
- `--engine` — optional engine guard: mariadb|postgres (defaults to the database's engine)

#### `dotnext-cli db user`

Manage database users (mariadb / postgres)

```
dotnext-cli db user
```

##### `dotnext-cli db user create`

Create a database user (auto-generates password if --password omitted)

```
dotnext-cli db user create [flags]
```

**Flags:**

- `--as-admin` — Skip the panel-username prefix (admin-only DB user names)
- `--engine` — Engine: mariadb | postgres (default `mariadb`)
- `--name` — DB user name (without panel-username prefix) — required
- `--password` — Password (auto-generated ULID if omitted; revealed once)
- `--user` — Panel user (email or username) — required

##### `dotnext-cli db user delete`

Delete a database user by ID

```
dotnext-cli db user delete [flags]
```

**Flags:**

- `--id` — DB user ID (ULID)

##### `dotnext-cli db user grant`

Grant a db user privileges on a database

```
dotnext-cli db user grant [flags]
```

**Flags:**

- `--db-name` — Database name (with panel-prefix) — required
- `--db-user-id` — DB user ID (ULID) — required
- `--level` — Shortcut: 'rw' or 'ro' (alternative to --privileges)
- `--privileges` — MariaDB privilege list (e.g. SELECT,INSERT,UPDATE) (default `[]`)

###### `dotnext-cli db user grant revoke`

Revoke a single grant, keeping the database user

```
dotnext-cli db user grant revoke <grant-id>
```

###### `dotnext-cli db user grant update`

Change a grant's level (rw|ro)

```
dotnext-cli db user grant update <grant-id> [flags]
```

**Flags:**

- `--level` — new grant level: 'rw' or 'ro' (required)

##### `dotnext-cli db user list`

List database users (filtered by panel user, or all)

```
dotnext-cli db user list [flags]
```

**Flags:**

- `--user` — Filter by panel user (email or username)

##### `dotnext-cli db user rotate-password`

Rotate a database user's password (revealed once)

```
dotnext-cli db user rotate-password <db-user-id>
```

### `dotnext-cli disk-usage`

Inspect / refresh a tenant's disk-usage breakdown (files / email / databases)

```
dotnext-cli disk-usage
```

#### `dotnext-cli disk-usage refresh`

Recompute a tenant's disk usage live and store the snapshot

```
dotnext-cli disk-usage refresh <user-email|username|id>
```

#### `dotnext-cli disk-usage refresh-all`

Recompute the disk-usage snapshot for every tenant (serial; skips accounts refreshed within --max-age)

```
dotnext-cli disk-usage refresh-all [flags]
```

**Flags:**

- `--max-age` — skip a tenant whose snapshot is younger than this (0 = always refresh) (default `20h0m0s`)
- `--per-user-timeout` — per-tenant deadline for the live compute (default `1m30s`)
- `--sleep` — pause between tenants to keep the sweep gentle (default `2s`)

#### `dotnext-cli disk-usage show`

Show a tenant's last stored disk-usage snapshot

```
dotnext-cli disk-usage show <user-email|username|id>
```

### `dotnext-cli dns`

DNS zones and records (dns_zones / dns_records): list / add / update / delete

```
dotnext-cli dns
```

#### `dotnext-cli dns prune-orphan-records`

Remove PowerDNS backend rows left behind by a zone delete in an earlier release

```
dotnext-cli dns prune-orphan-records [flags]
```

**Flags:**

- `--apply` — actually delete (default is a dry run)

#### `dotnext-cli dns prune-service-records`

Remove imported cPanel service subdomains (cpanel/webmail/whm/…) from DNS

```
dotnext-cli dns prune-service-records [flags]
```

**Flags:**

- `--apply` — actually delete (default is a dry run)
- `--zone` — limit to one zone (default: every zone)

#### `dotnext-cli dns record`

DNS records: list / add / update / delete

```
dotnext-cli dns record
```

##### `dotnext-cli dns record add`

Add a record to a zone (same validation + conflict rules as the UI)

```
dotnext-cli dns record add <domain> --name <name> --type <TYPE> --content <content> [flags]
```

**Flags:**

- `--content` — record content (e.g. an IP, hostname, or TXT value)
- `--disabled` — create the record disabled (kept in DB, not served)
- `--name` — record name ('@' for apex)
- `--priority` — priority (MX/SRV) (default `0`)
- `--ttl` — TTL seconds (60-604800; 0 uses default 300) (default `0`)
- `--type` — record type: A, AAAA, CNAME, MX, TXT, NS, SRV, CAA

##### `dotnext-cli dns record delete`

Delete a record from a zone

```
dotnext-cli dns record delete <domain> <record-id> [flags]
```

**Flags:**

- `--force` — confirm deletion

##### `dotnext-cli dns record list`

List records in a zone

```
dotnext-cli dns record list <domain>
```

##### `dotnext-cli dns record update`

Update a record's content / ttl / priority / enabled flag

```
dotnext-cli dns record update <domain> <record-id> [flags]
```

**Flags:**

- `--content` — new content
- `--enabled` — set enabled flag (default `true`)
- `--priority` — new priority (MX/SRV) (default `0`)
- `--ttl` — new TTL seconds (60-604800) (default `0`)

#### `dotnext-cli dns zone`

DNS zones: list / show / delete / enable

```
dotnext-cli dns zone
```

##### `dotnext-cli dns zone delete`

Drop the DNS facet of a domain (host DNS elsewhere; keeps web + mail)

```
dotnext-cli dns zone delete <domain> [flags]
```

**Flags:**

- `--force` — confirm dropping the zone

##### `dotnext-cli dns zone enable`

Re-enable DNS management for a domain (host DNS here again)

```
dotnext-cli dns zone enable <domain>
```

##### `dotnext-cli dns zone list`

List all DNS zones

```
dotnext-cli dns zone list
```

##### `dotnext-cli dns zone show`

Show a zone and its records

```
dotnext-cli dns zone show <domain>
```

### `dotnext-cli docker`

Manage the docker engine + app-marketplace host

```
dotnext-cli docker
```

#### `dotnext-cli docker disable`

Disable the marketplace toggle (keeps docker installed)

```
dotnext-cli docker disable
```

#### `dotnext-cli docker enable`

Install docker engine + flip Server Settings toggle

```
dotnext-cli docker enable
```

#### `dotnext-cli docker enable-tenant`

Enable userns-remap + tenant docker apps on this host

```
dotnext-cli docker enable-tenant [flags]
```

**Flags:**

- `--yes` — proceed without the interactive confirmation

#### `dotnext-cli docker status`

Show docker engine status (active, marketplace toggle state)

```
dotnext-cli docker status
```

### `dotnext-cli docker-app`

Manage Docker app catalog installs (admin-only)

```
dotnext-cli docker-app
```

#### `dotnext-cli docker-app backup-create`

Create a manual backup of the app

```
dotnext-cli docker-app backup-create <id>
```

#### `dotnext-cli docker-app backups`

List restic backups taken for this install

```
dotnext-cli docker-app backups <id>
```

#### `dotnext-cli docker-app catalog`

List entries in the installed catalog

```
dotnext-cli docker-app catalog
```

#### `dotnext-cli docker-app delete`

Uninstall a docker app (stops the stack, removes its row)

```
dotnext-cli docker-app delete <id> [flags]
```

**Flags:**

- `--keep-volumes` — keep /var/lib/dotnext/docker-apps/<slug> data on disk

#### `dotnext-cli docker-app edit`

Edit a Docker app's domain and/or ports (container must be stopped)

```
dotnext-cli docker-app edit <id> [flags]
```

**Flags:**

- `--clear-domain` — detach the managed domain
- `--domain` — attach this domain (empty string detaches)
- `--port` — port override (repeatable): name=..,host=..,bind=..,proxy=..,enabled=.. (default `[]`)

#### `dotnext-cli docker-app env`

View, edit, or regenerate a Docker app's environment

```
dotnext-cli docker-app env
```

##### `dotnext-cli docker-app env get`

List the app's environment (reveals secrets)

```
dotnext-cli docker-app env get <id>
```

##### `dotnext-cli docker-app env regenerate`

Mint a fresh value for a generated secret, then recreate the container

```
dotnext-cli docker-app env regenerate <id> <KEY>
```

##### `dotnext-cli docker-app env set`

Edit env values, re-render the compose, and recreate the container

```
dotnext-cli docker-app env set <id> <KEY=VALUE> [<KEY=VALUE>...]
```

#### `dotnext-cli docker-app exec`

Run a one-off command in the app's container

```
dotnext-cli docker-app exec <id> -- <command...> [flags]
```

**Flags:**

- `--service` — compose service to exec into (default: primary)

#### `dotnext-cli docker-app install`

Install a catalog entry (creates the docker_apps row + dispatches the agent)

```
dotnext-cli docker-app install <slug> [flags]
```

**Flags:**

- `--cpu` — cgroup CPU limit (e.g. 1.0). Catalog default when omitted.
- `--domain` — domain the tenant app attaches to (required with --user; must be owned by that user or free)
- `--env` — KEY=VALUE override (repeatable) (default `[]`)
- `--memory` — memory limit (e.g. 512m). Catalog default when omitted.
- `--name` — install name (lowercase, ^[a-z0-9-]{1,32}$)
- `--pids` — pids cgroup limit. Catalog default when omitted. (default `0`)
- `--update-mode` — manual|auto (default `manual`)
- `--user` — install for this tenant (user id or username); enables the tenant-scoped install path

#### `dotnext-cli docker-app list`

List installed docker apps

```
dotnext-cli docker-app list
```

#### `dotnext-cli docker-app logs`

Tail container logs

```
dotnext-cli docker-app logs <id> [flags]
```

**Flags:**

- `--lines` — lines to tail (default `200`)
- `--service` — compose service name (default: first service)

#### `dotnext-cli docker-app maintenance`

Docker disk-usage report and prune

```
dotnext-cli docker-app maintenance
```

##### `dotnext-cli docker-app maintenance disk-usage`

Report Docker disk usage

```
dotnext-cli docker-app maintenance disk-usage
```

##### `dotnext-cli docker-app maintenance prune`

Prune unused Docker resources (reclaims disk)

```
dotnext-cli docker-app maintenance prune [flags]
```

**Flags:**

- `--all` — remove all unused images, not just dangling
- `--force` — confirm the prune
- `--volumes` — also prune unused volumes

#### `dotnext-cli docker-app patch`

Patch app settings (update mode, CPU/memory/PID limits)

```
dotnext-cli docker-app patch <id> [flags]
```

**Flags:**

- `--cpu` — CPU limit (e.g. 1.5; empty to clear)
- `--memory` — memory limit (e.g. 512m; empty to clear)
- `--pids` — PIDs limit (default `0`)
- `--update-mode` — update mode: manual|auto

#### `dotnext-cli docker-app rebuild`

Force-recreate (docker compose up --force-recreate)

```
dotnext-cli docker-app rebuild <id>
```

#### `dotnext-cli docker-app restart`

Restart an install

```
dotnext-cli docker-app restart <id>
```

#### `dotnext-cli docker-app restore`

Restore a selected backup snapshot (replaces the app's data)

```
dotnext-cli docker-app restore <id> <snapshot-id> [flags]
```

**Flags:**

- `--force` — confirm the destructive restore

#### `dotnext-cli docker-app start`

Start a stopped install

```
dotnext-cli docker-app start <id>
```

#### `dotnext-cli docker-app status`

Show full status of an installed app (DB row + agent status)

```
dotnext-cli docker-app status <id>
```

#### `dotnext-cli docker-app stop`

Stop a running install

```
dotnext-cli docker-app stop <id>
```

#### `dotnext-cli docker-app update`

Pull the latest image and re-create the stack

```
dotnext-cli docker-app update <id>
```

### `dotnext-cli domain`

Manage hosted domains

```
dotnext-cli domain
```

#### `dotnext-cli domain bandwidth`

Per-domain bandwidth report (bytes + requests, daily series)

```
dotnext-cli domain bandwidth <domain-name|domain-id> [flags]
```

**Flags:**

- `--days` — lookback window in days (1..365) (default `30`)

#### `dotnext-cli domain bot-exclusion`

Manage the paths of a domain that skip the bot-detection challenge

```
dotnext-cli domain bot-exclusion
```

##### `dotnext-cli domain bot-exclusion add`

Let requests to a path of a domain skip the bot-detection challenge

```
dotnext-cli domain bot-exclusion add <domain> <path-prefix> [flags]
```

**Flags:**

- `--method` — only these methods (GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS); repeat or comma-separate (default `[]`)
- `--note` — why the path skips the challenge

##### `dotnext-cli domain bot-exclusion delete`

Remove a path of a domain from the bot-detection exclusions

```
dotnext-cli domain bot-exclusion delete <domain> <id>
```

##### `dotnext-cli domain bot-exclusion list`

List the paths of a domain that skip the bot-detection challenge

```
dotnext-cli domain bot-exclusion list <domain>
```

#### `dotnext-cli domain browse`

List a path under the domain's docroot (as the domain owner)

```
dotnext-cli domain browse <domain-name|domain-id> [relative-path]
```

#### `dotnext-cli domain catchall`

Manage per-domain catch-all routing

```
dotnext-cli domain catchall
```

##### `dotnext-cli domain catchall clear`

Remove the catch-all rule for a domain

```
dotnext-cli domain catchall clear <domain-name-or-id>
```

##### `dotnext-cli domain catchall set`

Route mail to unknown@<domain> to --target

```
dotnext-cli domain catchall set <domain-name-or-id> [flags]
```

**Flags:**

- `--target` — Destination email address (required)

##### `dotnext-cli domain catchall show`

Print the current catch-all target

```
dotnext-cli domain catchall show <domain-name-or-id>
```

#### `dotnext-cli domain chown`

Reassign a domain to a different owner, moving its docroot

```
dotnext-cli domain chown <domain> <new-owner> [flags]
```

**Flags:**

- `--yes` — skip the confirmation prompt

#### `dotnext-cli domain create`

Create a new domain (direct DB; bypasses HTTP auth)

```
dotnext-cli domain create [flags]
```

**Flags:**

- `--dns-template` — Admin-defined custom DNS template ID: seeds the template's records into the fresh zone and sets an external mail posture. Requires --manage-dns; mutually exclusive with --mail
- `--doc-root` — Document root (optional, auto-generated if not provided)
- `--mail` — Mail provider: dotnext | none | m365 | google. 'none' for a web/DNS-only domain (default `dotnext`)
- `--manage-dns` — Host this domain's DNS zone on this server. --manage-dns=false when DNS lives elsewhere (external DNS) (default `true`)
- `--name` — Domain name (required)
- `--reverse-proxy` — Make this a reverse-proxy domain: allocate a loopback port and proxy '/' to it
- `--reverse-proxy-port` — Reverse-proxy to this specific loopback port; 0 = auto-assign from the pool (default `0`)
- `--user` — User email, username, or ULID (required)
- `--web-enabled` — Host a website for this domain (vhost + docroot). --web-enabled=false makes a docroot-less DNS-only zone or mail-only domain (default `true`)

#### `dotnext-cli domain delete`

Delete a domain (direct DB; reconciler tears down nginx)

```
dotnext-cli domain delete <domain-name|domain-id> [flags]
```

**Flags:**

- `--force` — Skip confirmation prompt

#### `dotnext-cli domain dir-privacy`

Manage password-protected directories (auth_basic)

```
dotnext-cli domain dir-privacy
```

##### `dotnext-cli domain dir-privacy add`

Protect a directory

```
dotnext-cli domain dir-privacy add <domain> --path /dir [--realm ...] [flags]
```

**Flags:**

- `--path` — directory path under the docroot (e.g. /admin)
- `--realm` — auth realm shown in the browser prompt (default Restricted)

##### `dotnext-cli domain dir-privacy cred-add`

Add a credential to a protected directory

```
dotnext-cli domain dir-privacy cred-add <domain> <rule-id> --user <name> (--password <p> | --password-stdin) [flags]
```

**Flags:**

- `--password` — basic-auth password (8-128 chars; prefer --password-stdin)
- `--password-stdin` — read the password from stdin (no argv leak)
- `--user` — basic-auth username

##### `dotnext-cli domain dir-privacy cred-list`

List credentials for a protected directory

```
dotnext-cli domain dir-privacy cred-list <domain> <rule-id>
```

##### `dotnext-cli domain dir-privacy cred-remove`

Remove a credential from a protected directory

```
dotnext-cli domain dir-privacy cred-remove <domain> <rule-id> <cred-id>
```

##### `dotnext-cli domain dir-privacy list`

List protected directories for a domain

```
dotnext-cli domain dir-privacy list <domain>
```

##### `dotnext-cli domain dir-privacy remove`

Remove a protected directory (and its credentials)

```
dotnext-cli domain dir-privacy remove <domain> <rule-id>
```

#### `dotnext-cli domain disable`

Disable a domain (direct DB)

```
dotnext-cli domain disable <domain-name|domain-id>
```

#### `dotnext-cli domain disclaimer`

Manage per-domain outbound disclaimer

```
dotnext-cli domain disclaimer
```

##### `dotnext-cli domain disclaimer clear`

Disable + remove the outbound disclaimer for a domain

```
dotnext-cli domain disclaimer clear <domain-name-or-id>
```

##### `dotnext-cli domain disclaimer set`

Set + enable the outbound disclaimer for a domain

```
dotnext-cli domain disclaimer set <domain-name-or-id> [flags]
```

**Flags:**

- `--file` — Read disclaimer from this file (overrides --text)
- `--text` — Disclaimer text (UTF-8 plain text; HTML mail shows it escaped)

##### `dotnext-cli domain disclaimer show`

Print the current disclaimer for a domain

```
dotnext-cli domain disclaimer show <domain-name-or-id>
```

#### `dotnext-cli domain email-disable`

Disable email for a domain (keeps the DKIM key)

```
dotnext-cli domain email-disable <domain-name-or-id>
```

#### `dotnext-cli domain email-dkim-rotate`

Rotate the domain's DKIM keys (operator-driven, not automatic)

```
dotnext-cli domain email-dkim-rotate <domain-name-or-id>
```

#### `dotnext-cli domain email-enable`

Enable email for a domain (generates DKIM + publishes DNS records)

```
dotnext-cli domain email-enable <domain-name-or-id>
```

#### `dotnext-cli domain enable`

Enable a domain (direct DB)

```
dotnext-cli domain enable <domain-name|domain-id>
```

#### `dotnext-cli domain fix-perms`

Repair docroot group/setgid (www-data) for all of a user's domains

```
dotnext-cli domain fix-perms <username>
```

#### `dotnext-cli domain htaccess-preview`

Preview the nginx translation of a .htaccess file (stateless; nothing applied)

```
dotnext-cli domain htaccess-preview <domain-name|domain-id> [flags]
```

**Flags:**

- `--file` — path to a .htaccess file to convert (required)

#### `dotnext-cli domain ip-acl`

Manage a domain's IP allow/deny ACL

```
dotnext-cli domain ip-acl
```

##### `dotnext-cli domain ip-acl add`

Add an IP ACL entry to a domain

```
dotnext-cli domain ip-acl add <domain> [flags]
```

**Flags:**

- `--action` — allow|deny (default `deny`)
- `--cidr` — IP or CIDR (required)
- `--comment` — optional comment
- `--priority` — rule priority (lower = first) (default `0`)

##### `dotnext-cli domain ip-acl delete`

Delete an IP ACL entry by id

```
dotnext-cli domain ip-acl delete <acl-id>
```

##### `dotnext-cli domain ip-acl list`

List a domain's IP ACL entries

```
dotnext-cli domain ip-acl list <domain>
```

#### `dotnext-cli domain list`

List domains (direct DB)

```
dotnext-cli domain list
```

#### `dotnext-cli domain mta-sts`

Enable/disable MTA-STS for a mail domain

```
dotnext-cli domain mta-sts <domain-name|domain-id> [flags]
```

**Flags:**

- `--disable` — disable MTA-STS
- `--enable` — enable MTA-STS

#### `dotnext-cli domain ownership`

Domain ownership proof: status, verify, approve, revoke, pending, policy

```
dotnext-cli domain ownership
```

##### `dotnext-cli domain ownership approve`

Approve a pending domain as an administrator (audited)

```
dotnext-cli domain ownership approve <domain-name|domain-id>
```

##### `dotnext-cli domain ownership pending`

List the domains waiting for ownership proof

```
dotnext-cli domain ownership pending
```

##### `dotnext-cli domain ownership policy`

Show or switch whether new domain names need ownership proof (audited)

```
dotnext-cli domain ownership policy [on|off] [flags]
```

**Flags:**

- `--yes` — confirm switching ownership proof off

##### `dotnext-cli domain ownership revoke`

Send a verified domain back to pending as an administrator (audited)

```
dotnext-cli domain ownership revoke <domain-name|domain-id>
```

##### `dotnext-cli domain ownership status`

Show a domain's ownership state and the TXT record that proves it

```
dotnext-cli domain ownership status <domain-name|domain-id>
```

##### `dotnext-cli domain ownership verify`

Check a pending domain's TXT record through the public resolvers now

```
dotnext-cli domain ownership verify <domain-name|domain-id>
```

#### `dotnext-cli domain php-settings`

Get/set a domain's php.ini directives

```
dotnext-cli domain php-settings
```

##### `dotnext-cli domain php-settings get`

Show a domain's php.ini directives

```
dotnext-cli domain php-settings get <domain-name-or-id>
```

##### `dotnext-cli domain php-settings set`

Set php.ini directives (only the flags you pass change; reconciler converges)

```
dotnext-cli domain php-settings set <domain-name-or-id> [flags]
```

**Flags:**

- `--max-execution-time` — max_execution_time seconds (1..86400) (default `0`)
- `--max-input-time` — max_input_time seconds (1..86400) (default `0`)
- `--max-input-vars` — max_input_vars (1..86400) (default `0`)
- `--memory-limit` — memory_limit (e.g. 256M)
- `--post-max-size` — post_max_size (e.g. 64M)
- `--upload-max-filesize` — upload_max_filesize (e.g. 64M)

#### `dotnext-cli domain php-version`

Manage a domain's PHP version (per-domain FPM pool)

```
dotnext-cli domain php-version
```

##### `dotnext-cli domain php-version set`

Bind a domain to a PHP version (find-or-create the pool; reconciler converges)

```
dotnext-cli domain php-version set <domain-name-or-id> --version X.Y [flags]
```

**Flags:**

- `--version` — PHP version X.Y (e.g. 8.2); required

##### `dotnext-cli domain php-version show`

Show a domain's bound PHP version

```
dotnext-cli domain php-version show <domain-name-or-id>
```

#### `dotnext-cli domain prune-orphans`

List sites in nginx sites-enabled that have no panel DB row (and optionally delete them)

```
dotnext-cli domain prune-orphans [flags]
```

**Flags:**

- `--apply` — Actually delete orphans (default: dry-run)

#### `dotnext-cli domain set`

Set advanced domain settings (redirect / nginx directives / index / ssl mode / cache)

```
dotnext-cli domain set <domain-name|domain-id> [flags]
```

**Flags:**

- `--cache` — on|off — nginx fastcgi page cache
- `--index-priority` — directory index priority (e.g. html_first, php_first)
- `--nginx-directives` — raw custom nginx directives for the server block
- `--redirect-all-to` — redirect the whole domain to this URL ('' clears)
- `--redirect-type` — 301|302|307|308 (permanent=301, temporary=302)
- `--ssl-mode` — le|self|none (custom = install a cert)

#### `dotnext-cli domain show`

Show a domain's full advanced-settings state (JSON)

```
dotnext-cli domain show <domain-name|domain-id>
```

#### `dotnext-cli domain skip-auto-san`

Opt a domain in/out of the panel-cert auto-SAN

```
dotnext-cli domain skip-auto-san <domain-name|domain-id> [flags]
```

**Flags:**

- `--disable` — include this domain in the panel cert SAN
- `--enable` — exclude this domain from the panel cert SAN

#### `dotnext-cli domain whois`

WHOIS lookup for a domain

```
dotnext-cli domain whois <domain-name|domain-id>
```

### `dotnext-cli dr`

Disaster-recovery standby: pair, status, promote

```
dotnext-cli dr
```

#### `dotnext-cli dr feed`

Primary side: ship system backups to the DR destination on a cadence

```
dotnext-cli dr feed [flags]
```

**Flags:**

- `--cron` — cron cadence for the DR feed (default hourly, "0 * * * *")
- `--destination` — DR backup destination ID to ship system + account backups to

#### `dotnext-cli dr pair`

Designate this box a DR standby of a primary

```
dotnext-cli dr pair [flags]
```

**Flags:**

- `--destination` — backup destination ID used as the read-only DR channel
- `--peer-label` — human label for the primary this box replicates (e.g. its hostname)

#### `dotnext-cli dr promote`

Promote this DR standby to the live primary

```
dotnext-cli dr promote [flags]
```

**Flags:**

- `--force` — promote even if the old primary still answers (operator asserts it is down)
- `--skip-restore` — skip the final account-inclusive restore (role flip only)
- `--yes` — skip the interactive confirmation (scriptable)

#### `dotnext-cli dr status`

Show this box's DR role + pairing

```
dotnext-cli dr status
```

#### `dotnext-cli dr unfeed`

Primary side: stop shipping system backups to a DR destination (inverse of `dr feed`)

```
dotnext-cli dr unfeed [flags]
```

**Flags:**

- `--destination` — the DR backup destination id to stop feeding

#### `dotnext-cli dr unpair`

Revert this box to a primary (clears standby role)

```
dotnext-cli dr unpair
```

### `dotnext-cli files`

Scoped tenant file manager (list/read/mkdir/move/chmod/archive/…) — same policy as the GUI

```
dotnext-cli files
```

#### `dotnext-cli files archive`

tar.gz one or more paths to a local file (--out)

```
dotnext-cli files archive <path...> [flags]
```

**Flags:**

- `--out` — local .tar.gz destination (required)
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files chmod`

Change a path's permissions (e.g. 0644)

```
dotnext-cli files chmod <path> <octal-mode> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files copy`

Copy into a destination directory (--to)

```
dotnext-cli files copy <path> [flags]
```

**Flags:**

- `--to` — destination directory (required)
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files delete`

Delete a file or directory

```
dotnext-cli files delete <path> [flags]
```

**Flags:**

- `--force` — confirm deletion
- `--recursive` — recurse into directories
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files download`

Download a file to a local path (--out)

```
dotnext-cli files download <path> [flags]
```

**Flags:**

- `--out` — local destination path (required)
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files extract`

Extract an archive inside the tenant tree (--dest)

```
dotnext-cli files extract <archive-path> [flags]
```

**Flags:**

- `--dest` — destination directory (default: archive's directory)
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files list`

List a directory

```
dotnext-cli files list <path> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files mkdir`

Create a directory

```
dotnext-cli files mkdir <path> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files move`

Move into a destination directory (--to)

```
dotnext-cli files move <path> [flags]
```

**Flags:**

- `--to` — destination directory (required)
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files read`

Print a file's contents

```
dotnext-cli files read <path> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files rename`

Rename within the same directory

```
dotnext-cli files rename <path> <new-name> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files stat`

Stat a path

```
dotnext-cli files stat <path> [flags]
```

**Flags:**

- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files upload`

Upload a local file into the tenant tree

```
dotnext-cli files upload <local-file> <dest-path> [flags]
```

**Flags:**

- `--overwrite` — overwrite if the destination exists
- `--user` — target tenant (email|username|id, required)

#### `dotnext-cli files write`

Write text to a file (--content or --from <local>)

```
dotnext-cli files write <path> [flags]
```

**Flags:**

- `--content` — inline text content
- `--from` — read content from this local file
- `--user` — target tenant (email|username|id, required)

### `dotnext-cli impersonation`

Manage admin act-as impersonation grants (list / create / end)

```
dotnext-cli impersonation
```

#### `dotnext-cli impersonation create`

Create a 60-minute impersonation grant for an admin to act as a target user

```
dotnext-cli impersonation create --admin <admin> --target <user> [flags]
```

**Flags:**

- `--admin` — admin who will act-as (email|username|id, required)
- `--target` — target user to impersonate (email|username|id, required)

#### `dotnext-cli impersonation end`

End an active impersonation grant

```
dotnext-cli impersonation end <grant-id>
```

#### `dotnext-cli impersonation list`

List active impersonation grants (all, or for one admin)

```
dotnext-cli impersonation list [flags]
```

**Flags:**

- `--admin` — scope to one admin (email|username|id); default = all admins

### `dotnext-cli ip`

Admin IP address pool (managed_ips): list / add / update / delete

```
dotnext-cli ip
```

#### `dotnext-cli ip add`

Add an IP to the managed pool

```
dotnext-cli ip add <address> [flags]
```

**Flags:**

- `--label` — human label
- `--user-selectable` — tenants may pick this IP

#### `dotnext-cli ip delete`

Delete an IP from the pool

```
dotnext-cli ip delete <id> [flags]
```

**Flags:**

- `--force` — confirm deletion

#### `dotnext-cli ip list`

List managed IP addresses

```
dotnext-cli ip list
```

#### `dotnext-cli ip update`

Update an IP's label / user-selectable / default flag

```
dotnext-cli ip update <id> [flags]
```

**Flags:**

- `--default` — promote to per-family default
- `--label` — new label
- `--user-selectable` — tenants may pick this IP

### `dotnext-cli limits`

Per-user resource limits (cgroups v2 + POSIX quota + nginx)

```
dotnext-cli limits
```

#### `dotnext-cli limits apply`

Re-apply effective limits for one user

```
dotnext-cli limits apply <username>
```

#### `dotnext-cli limits check`

Probe host for resource-limit prerequisites (cgroups v2, /home fs, nginx modules)

```
dotnext-cli limits check
```

#### `dotnext-cli limits override`

Per-user limit overrides (set / clear)

```
dotnext-cli limits override
```

##### `dotnext-cli limits override clear`

Remove a user's override (revert to package)

```
dotnext-cli limits override clear <username|user-id>
```

##### `dotnext-cli limits override set`

Set/replace a user's limit override

```
dotnext-cli limits override set <username|user-id> [flags]
```

**Flags:**

- `--cpu` — CPU quota percent (default `0`)
- `--disk-mb` — disk quota MB (default `0`)
- `--io-read-mbps` — IO read MB/s (default `0`)
- `--io-write-mbps` — IO write MB/s (default `0`)
- `--max-tasks` — max tasks (default `0`)
- `--memory-mb` — memory limit MB (default `0`)

#### `dotnext-cli limits package`

Bulk limits operations across every user of a package

```
dotnext-cli limits package
```

##### `dotnext-cli limits package apply`

Re-apply limits to every user of the given package

```
dotnext-cli limits package apply <package_id> [flags]
```

**Flags:**

- `--dry-run` — print what would be applied without making agent calls

#### `dotnext-cli limits status`

Show live resource usage for one user

```
dotnext-cli limits status <username>
```

#### `dotnext-cli limits usage`

Show a user's live limit usage (disk/cpu/mem/io/tasks)

```
dotnext-cli limits usage <username|user-id>
```

### `dotnext-cli log`

Inspect log types and mint/revoke log-stream access grants

```
dotnext-cli log
```

#### `dotnext-cli log access`

Mint or revoke log-stream access grants

```
dotnext-cli log access
```

##### `dotnext-cli log access create`

Mint a 15-minute log-stream access grant for a user

```
dotnext-cli log access create --user <user> --type <access|error|goaccess> [--domain <domain>] [flags]
```

**Flags:**

- `--domain` — scope to one domain (optional)
- `--type` — log type: access|error|goaccess (required)
- `--user` — grant owner (email|username|id, required)

##### `dotnext-cli log access revoke`

Revoke a log-stream access grant

```
dotnext-cli log access revoke <stream-key>
```

#### `dotnext-cli log types`

List the supported log types

```
dotnext-cli log types
```

### `dotnext-cli mail`

Admin mail queue + outbound throttle management

```
dotnext-cli mail
```

#### `dotnext-cli mail deliverability`

Show the mail deliverability score (server-wide or per-domain)

```
dotnext-cli mail deliverability [flags]
```

**Flags:**

- `--domain` — scope to one domain (default: server-wide)

#### `dotnext-cli mail logs`

Search / inspect mail delivery logs

```
dotnext-cli mail logs
```

##### `dotnext-cli mail logs detail`

Show a mail log entry's detail

```
dotnext-cli mail logs detail <id>
```

##### `dotnext-cli mail logs query`

Search mail logs

```
dotnext-cli mail logs query [flags]
```

**Flags:**

- `--from` — from date (YYYY-MM-DD)
- `--limit` — max entries (default `50`)
- `--offset` — offset (default `0`)
- `--recipient` — recipient prefix filter
- `--sender` — sender prefix filter
- `--to` — to date (YYYY-MM-DD)

#### `dotnext-cli mail queue`

Inspect / retry / delete queued mail

```
dotnext-cli mail queue
```

##### `dotnext-cli mail queue delete`

Delete a queued message

```
dotnext-cli mail queue delete <id> [flags]
```

**Flags:**

- `--force` — confirm deletion

##### `dotnext-cli mail queue list`

List queued mail

```
dotnext-cli mail queue list
```

##### `dotnext-cli mail queue retry`

Retry a queued message

```
dotnext-cli mail queue retry <id>
```

#### `dotnext-cli mail throttle`

Outbound mail throttle policies

```
dotnext-cli mail throttle
```

##### `dotnext-cli mail throttle create`

Create an outbound throttle policy

```
dotnext-cli mail throttle create [flags]
```

**Flags:**

- `--disabled` — create disabled
- `--max-per-day` — max messages/day (0 = unlimited) (default `0`)
- `--max-per-hour` — max messages/hour (0 = unlimited) (default `0`)
- `--scope` — user|domain|global (required)
- `--scope-ref` — email (user) or FQDN (domain); omit for global

##### `dotnext-cli mail throttle delete`

Delete a throttle policy

```
dotnext-cli mail throttle delete <id>
```

##### `dotnext-cli mail throttle list`

List outbound throttle policies

```
dotnext-cli mail throttle list
```

##### `dotnext-cli mail throttle update`

Update a throttle policy's limits / enabled flag

```
dotnext-cli mail throttle update <id> [flags]
```

**Flags:**

- `--disable` — disable the policy
- `--enable` — enable the policy
- `--max-per-day` — max messages/day (default `0`)
- `--max-per-hour` — max messages/hour (default `0`)

### `dotnext-cli mail-cert`

Manage a domain's mail TLS certificate (mail.<domain> SAN)

```
dotnext-cli mail-cert
```

#### `dotnext-cli mail-cert disable`

Disable the mail certificate for a domain

```
dotnext-cli mail-cert disable <domain-name|domain-id>
```

#### `dotnext-cli mail-cert enable`

Enable the mail certificate for a domain (issues on the next reconcile pass)

```
dotnext-cli mail-cert enable <domain-name|domain-id>
```

#### `dotnext-cli mail-cert reissue`

Clear backoff and re-issue a domain's mail certificate

```
dotnext-cli mail-cert reissue <domain-name|domain-id>
```

#### `dotnext-cli mail-cert show`

Show a domain's mail certificate state

```
dotnext-cli mail-cert show <domain-name|domain-id>
```

### `dotnext-cli mail-group`

Manage mail groups (distribution lists + shared-resource groups) and members

```
dotnext-cli mail-group
```

#### `dotnext-cli mail-group add-member`

Add one member to a mail group

```
dotnext-cli mail-group add-member <group-email|group-id> <member-email-or-id>
```

#### `dotnext-cli mail-group create`

Create a mail group on a domain

```
dotnext-cli mail-group create <domain-name|domain-id> <local-part> [flags]
```

**Flags:**

- `--description` — description
- `--display-name` — display name (the From name)
- `--internal-only` — reject mail from senders outside the group's domain
- `--kind` — group kind: resource | distribution (default `resource`)

#### `dotnext-cli mail-group delete`

Delete a mail group (strips the Stalwart principal first)

```
dotnext-cli mail-group delete <group-email|group-id> [flags]
```

**Flags:**

- `--force` — confirm deletion

#### `dotnext-cli mail-group list`

List mail groups (all, or scoped to one domain)

```
dotnext-cli mail-group list [domain-name|domain-id]
```

#### `dotnext-cli mail-group remove-member`

Remove one member from a mail group

```
dotnext-cli mail-group remove-member <group-email|group-id> <member-email-or-id>
```

#### `dotnext-cli mail-group set-external`

Replace the outside addresses of a distribution list (none clears them)

```
dotnext-cli mail-group set-external <group-email|group-id> [address ...]
```

#### `dotnext-cli mail-group set-members`

Replace the full member set of a mail group

```
dotnext-cli mail-group set-members <group-email|group-id> [member-email-or-id ...]
```

#### `dotnext-cli mail-group set-resources`

Toggle a resource group's shared collections (mailbox/calendar/addressbook/files)

```
dotnext-cli mail-group set-resources <group-email|group-id> [flags]
```

**Flags:**

- `--addressbook` — shared addressbook (true|false)
- `--calendar` — shared calendar (true|false)
- `--files` — shared files (true|false)
- `--mailbox` — shared mailbox (true|false)

#### `dotnext-cli mail-group show`

Show a mail group and its members

```
dotnext-cli mail-group show <group-email|group-id>
```

### `dotnext-cli mailbox`

Manage mailboxes (email via Stalwart)

```
dotnext-cli mailbox
```

#### `dotnext-cli mailbox autoresponder`

Manage per-mailbox vacation responders

```
dotnext-cli mailbox autoresponder
```

##### `dotnext-cli mailbox autoresponder clear`

Disable + delete the autoresponder for a mailbox

```
dotnext-cli mailbox autoresponder clear <email>
```

##### `dotnext-cli mailbox autoresponder set`

Enable an autoresponder for a mailbox

```
dotnext-cli mailbox autoresponder set <email> [flags]
```

**Flags:**

- `--body` — Plain text body (optional if --html-body set)
- `--from` — Start date (RFC3339, e.g. 2026-05-01T00:00:00Z)
- `--html-body` — HTML body (optional)
- `--subject` — Subject line (required)
- `--to` — End date (RFC3339)

##### `dotnext-cli mailbox autoresponder show`

Print the current autoresponder for a mailbox

```
dotnext-cli mailbox autoresponder show <email>
```

#### `dotnext-cli mailbox create`

Create a mailbox (password shown once if auto-generated)

```
dotnext-cli mailbox create [flags]
```

**Flags:**

- `--display-name` — Human-readable display name
- `--domain` — Domain name or ID (required)
- `--local` — Local part, e.g. "alice" (required)
- `--password` — Explicit password (omit to auto-generate)
- `--quota-mb` — Disk quota in MiB (default 1024) (default `0`)
- `--send-only` — SMTP-submission only — authenticates to send but never receives (the only kind a domain whose mail is hosted elsewhere can have)

#### `dotnext-cli mailbox delete`

Delete a mailbox (agent destroys Stalwart account first)

```
dotnext-cli mailbox delete <email> [flags]
```

**Flags:**

- `--force` — Skip confirmation prompt

#### `dotnext-cli mailbox forwarder`

Manage per-mailbox aliases + external forwards

```
dotnext-cli mailbox forwarder
```

##### `dotnext-cli mailbox forwarder add`

Add an alias or external forwarder to a mailbox

```
dotnext-cli mailbox forwarder add <email> [flags]
```

**Flags:**

- `--keep-copy` — type=external: keep a copy in the mailbox (Sieve redirect :copy)
- `--local` — Alias local part (required for type=alias)
- `--target` — External destination email (required for type=external)
- `--type` — alias | external (required)

##### `dotnext-cli mailbox forwarder list`

List forwarders attached to a mailbox

```
dotnext-cli mailbox forwarder list <email>
```

##### `dotnext-cli mailbox forwarder remove`

Delete a forwarder by ID (find via 'forwarder list')

```
dotnext-cli mailbox forwarder remove <forwarder-id>
```

#### `dotnext-cli mailbox list`

List mailboxes in a domain

```
dotnext-cli mailbox list [flags]
```

**Flags:**

- `--domain` — Domain name or ID (required)

#### `dotnext-cli mailbox passwd`

Rotate a mailbox password (auto-generated and shown once if --password omitted)

```
dotnext-cli mailbox passwd <email> [flags]
```

**Flags:**

- `--password` — Explicit new password (omit to auto-generate)

#### `dotnext-cli mailbox set-quota`

Update a mailbox disk quota (in MiB)

```
dotnext-cli mailbox set-quota <email> <mb>
```

#### `dotnext-cli mailbox shares`

Manage shared mailbox folders

```
dotnext-cli mailbox shares
```

##### `dotnext-cli mailbox shares add`

Grant a target mailbox shared access to the owner's mailbox

```
dotnext-cli mailbox shares add [flags]
```

**Flags:**

- `--owner` — Owner mailbox email (required)
- `--rights` — Preset: ro | rw | admin (default rw) (default `rw`)
- `--shared-with` — Mailbox to grant share to (required)

##### `dotnext-cli mailbox shares list`

List shares for a given owner email

```
dotnext-cli mailbox shares list [flags]
```

**Flags:**

- `--owner` — Owner mailbox email (required)

##### `dotnext-cli mailbox shares remove`

Revoke a share by ID

```
dotnext-cli mailbox shares remove [flags]
```

**Flags:**

- `--id` — Share ID (ULID, from `dotnext mailbox shares list`)

### `dotnext-cli malware`

Malware scan / quarantine / YARA / settings (incident response)

```
dotnext-cli malware
```

#### `dotnext-cli malware events`

List malware detection events

```
dotnext-cli malware events [flags]
```

**Flags:**

- `--severity` — filter: info|warn|critical
- `--source` — filter: maldet|clamav|yara
- `--user` — filter by user id

#### `dotnext-cli malware quarantine`

List / restore / delete quarantined files

```
dotnext-cli malware quarantine
```

##### `dotnext-cli malware quarantine delete`

Permanently delete a quarantined file

```
dotnext-cli malware quarantine delete <id> [flags]
```

**Flags:**

- `--force` — confirm permanent deletion

##### `dotnext-cli malware quarantine list`

List active quarantine entries

```
dotnext-cli malware quarantine list [flags]
```

**Flags:**

- `--limit` — max rows (default `100`)

##### `dotnext-cli malware quarantine restore`

Restore a quarantined file to its original path

```
dotnext-cli malware quarantine restore <id> [flags]
```

**Flags:**

- `--reason` — why the file is being restored (required, min 5 chars)

#### `dotnext-cli malware scan`

Start a malware scan of a path or a user's home

```
dotnext-cli malware scan [flags]
```

**Flags:**

- `--path` — absolute path to scan
- `--user` — username whose home to scan

#### `dotnext-cli malware scan-status`

Poll a scan session's status

```
dotnext-cli malware scan-status <session-id>
```

#### `dotnext-cli malware settings`

Get / set malware settings

```
dotnext-cli malware settings
```

##### `dotnext-cli malware settings get`

Show malware settings

```
dotnext-cli malware settings get
```

##### `dotnext-cli malware settings set`

Update malware settings (same validation as the UI)

```
dotnext-cli malware settings set [flags]
```

**Flags:**

- `--max-scan-mb` — max scan file size MB (1..10240) (default `0`)
- `--notify-threshold` — info|warn|critical
- `--realtime` — on|off (real-time monitor)
- `--retain-days` — quarantine retention days (1..365) (default `0`)

#### `dotnext-cli malware status`

Show malware engine status

```
dotnext-cli malware status
```

#### `dotnext-cli malware update-signatures`

Update malware signatures (maldet + the signature-base YARA pack)

```
dotnext-cli malware update-signatures
```

#### `dotnext-cli malware yara`

Manage custom YARA rules

```
dotnext-cli malware yara
```

##### `dotnext-cli malware yara delete`

Delete a custom YARA rule

```
dotnext-cli malware yara delete <filename>
```

##### `dotnext-cli malware yara list`

List custom YARA rules

```
dotnext-cli malware yara list
```

##### `dotnext-cli malware yara toggle`

Enable or disable a custom YARA rule

```
dotnext-cli malware yara toggle <filename> [flags]
```

**Flags:**

- `--disable` — disable the rule
- `--enable` — enable the rule

##### `dotnext-cli malware yara upload`

Upload a custom YARA rule file

```
dotnext-cli malware yara upload <file.yar>
```

### `dotnext-cli malware-purge`

Hard-delete terminated malware quarantine rows past retention

```
dotnext-cli malware-purge
```

### `dotnext-cli migrate`

Database migration commands

```
dotnext-cli migrate
```

#### `dotnext-cli migrate force`

Clear the dirty flag by asserting the schema is at <version> (runs no SQL)

```
dotnext-cli migrate force <version>
```

#### `dotnext-cli migrate imap`

Migrate a remote IMAP mailbox into a dotnext mailbox

```
dotnext-cli migrate imap [flags]
```

**Flags:**

- `--allow-private` — permit RFC1918/loopback targets (migrating from a LAN server)
- `--csv` — batch: CSV with header host,user,password,to[,port,starttls]
- `--host` — remote IMAP host (e.g. imap.gmail.com)
- `--password-file` — read the app-password from a file (single account)
- `--password-stdin` — read the app-password from stdin (single account)
- `--port` — remote IMAP port (0 = 993 implicit-TLS, or 143 with --starttls) (default `0`)
- `--probe` — list the remote folder layout and exit without migrating
- `--starttls` — use STARTTLS on 143 instead of implicit TLS on 993
- `--to` — destination dotnext mailbox (must already exist)
- `--user` — remote IMAP login (full address)

#### `dotnext-cli migrate import`

Run (or resume) a migration job through the four-stage pipeline

```
dotnext-cli migrate import [flags]
```

**Flags:**

- `--allow-degraded` — exit 0 even if the restore ends degraded (a core area — DB/mail/health — failed)
- `--job-id` — migration_jobs.id (ULID) — required
- `--keep-staging` — do NOT delete /var/lib/dotnext-migrations/<job-id>/ after run (debug aid)
- `--preserve-source-state` — keep imported source state ACTIVE (mail forwarders/catchalls/filters/autoresponders, and — as later PRs land — credentials/SSL/SSH). Default OFF: imported artifacts land inert and the operator re-activates trusted ones in the panel. Only use for trusted same-owner migrations.
- `--target-email` — destination user email (only used when auto-creating)
- `--target-package-id` — hosting package ULID (only used when auto-creating)
- `--target-password` — destination user password (only used when auto-creating; ≥10 chars)
- `--target-user` — destination dotnext username — auto-created if --target-email + --target-password supplied

#### `dotnext-cli migrate import-wp`

Import a staged wordpress_ssh migration into a destination

```
dotnext-cli migrate import-wp [flags]
```

**Flags:**

- `--dest-domain` — destination domain (docroot = /home/<user>/domains/<domain>/public_html)
- `--dest-user` — destination OS username (owns the docroot)
- `--job-id` — wordpress_ssh migration_jobs.id (staged)

#### `dotnext-cli migrate pull-source`

Connect to source via SSH, run kind-appropriate backup, pull + extract tarball

```
dotnext-cli migrate pull-source [flags]
```

**Flags:**

- `--job-id` — migration_jobs.id (ULID) — required
- `--ssh-user` — SSH login on the source (default 'root') (default `root`)

#### `dotnext-cli migrate reap-secrets`

Wipe per-job migration-secrets env files + stale tarball/extracted dirs

```
dotnext-cli migrate reap-secrets [flags]
```

**Flags:**

- `--dry-run` — List would-delete paths without removing them
- `--staging-max-age` — Reap /var/lib/dotnext-migrations/<id>/ only when the job has been terminal at least this long (default 168h = 7d; pass 0 to wipe immediately) (default `168h0m0s`)

#### `dotnext-cli migrate refresh`

Force re-pull (refresh) an already-migrated account from a staged source

```
dotnext-cli migrate refresh [flags]
```

**Flags:**

- `--db` — Dest DB name (identity unchanged)
- `--docroot` — Dest docroot (overwritten)
- `--domain` — Dest domain (for cache purge)
- `--force` — Required — refresh overwrites a live account
- `--new-url` — New site URL (search-replace)
- `--old-url` — Old site URL (search-replace)
- `--os-user` — Dest Linux user
- `--source-docroot` — Staged source docroot
- `--source-sql` — Staged source SQL dump

#### `dotnext-cli migrate restore`

One-shot offline restore from a cpmove tarball (create job + stage + import)

```
dotnext-cli migrate restore [flags]
```

**Flags:**

- `--cpanel` — source is a cPanel cpmove / WHM pkgacct tarball
- `--directadmin` — source is a DirectAdmin user backup (user.<reseller>.<user>.tar.zst / .tar.gz / .tar)
- `--file` — path to the cpmove tarball (cpmove-<user>.tar.gz) — required
- `--fresh` — alias of --retry-from-scratch
- `--hestiacp` — source is a HestiaCP v-backup-user tarball (<user>.<YYYY-MM-DD_HH-MM-SS>.tar[.gz])
- `--keep-staging` — keep /var/lib/dotnext-migrations/<job-id>/ after the run (debug)
- `--preserve-source-state` — keep imported source state ACTIVE + carry source credentials where safe: preserves the mailbox password (Stalwart-verifiable bcrypt only), keeps mail forwarders/catchalls/filters/autoresponders active, restores DB user creds + source SSL. Default OFF (secure): mail gets a fresh password (tenant must reset), routing artifacts land inert. Parity with `dotnext migrate import`. Only use for a trusted same-owner migration.
- `--restore-file` — alias of --file
- `--retry-from-scratch` — reuse the source + options but wipe the existing job's stages and re-run the whole pipeline from analyze (recreates the target user, replaces stale manifest); default is a gentle resume
- `--source-host` — informational source host (offline restore leaves this empty)
- `--source-user` — cPanel account (default: derived from the cpmove filename)
- `--target-email` — destination email (only used when auto-creating the user)
- `--target-package-id` — hosting package ULID (auto-create only)
- `--target-password` — destination password (auto-create only; ≥10 chars)
- `--target-user` — destination dotnext username (default: the source account)

#### `dotnext-cli migrate status`

Show the schema version and whether it is dirty

```
dotnext-cli migrate status
```

#### `dotnext-cli migrate up`

Run pending database migrations

```
dotnext-cli migrate up
```

### `dotnext-cli migrate-from-dotnext`

Migrate an existing legacy DotNext Panel host layout to dotNext

```
dotnext-cli migrate-from-dotnext [flags]
```

**Flags:**

- `--apply` — perform the in-place migration to dotNext paths and units
- `--dry-run` — preview planned migrations without mutating files (default unless --apply is passed)

### `dotnext-cli notification`

Inspect notification channels and toggle event notifications

```
dotnext-cli notification
```

#### `dotnext-cli notification broadcast`

Broadcast a notification to every enabled channel

```
dotnext-cli notification broadcast --title <t> [--body <b>] [--severity info|warning|error|critical] [flags]
```

**Flags:**

- `--body` — notification body (<=2000)
- `--deeplink` — optional in-panel deeplink
- `--severity` — info|warning|error|critical (default info)
- `--title` — notification title (required, <=200)

#### `dotnext-cli notification channels`

Notification channels

```
dotnext-cli notification channels
```

##### `dotnext-cli notification channels create`

Create a notification channel

```
dotnext-cli notification channels create --name <n> --kind <email|slack|discord|ntfy|webhook|webpush|sms> --config <json> [flags]
```

**Flags:**

- `--config` — per-kind config JSON (e.g. '{"url":"https://…"}')
- `--disabled` — create in the disabled state
- `--kind` — email|slack|discord|ntfy|webhook|webpush|sms (required)
- `--name` — channel name (required)
- `--user` — owner user id (tenant-owned channel; omit for server-wide)

##### `dotnext-cli notification channels delete`

Delete a notification channel

```
dotnext-cli notification channels delete <id> [flags]
```

**Flags:**

- `--force` — confirm deletion

##### `dotnext-cli notification channels list`

List configured notification channels

```
dotnext-cli notification channels list
```

##### `dotnext-cli notification channels seal-secrets`

Encrypt any plaintext channel secrets at rest (one-time backfill)

```
dotnext-cli notification channels seal-secrets
```

##### `dotnext-cli notification channels test`

Send a synthetic test notification to one channel

```
dotnext-cli notification channels test <id>
```

#### `dotnext-cli notification dlq`

Inspect and operate the notification dead-letter queue

```
dotnext-cli notification dlq
```

##### `dotnext-cli notification dlq clear`

Delete ALL DLQ entries (does not replay)

```
dotnext-cli notification dlq clear [flags]
```

**Flags:**

- `--force` — confirm clearing the entire DLQ

##### `dotnext-cli notification dlq drop`

Delete one DLQ entry without replaying

```
dotnext-cli notification dlq drop <id>
```

##### `dotnext-cli notification dlq list`

List dead-lettered notifications (newest first)

```
dotnext-cli notification dlq list [flags]
```

**Flags:**

- `--limit` — max entries to show (1-500) (default `100`)

##### `dotnext-cli notification dlq replay`

Re-publish one DLQ entry to the main queue

```
dotnext-cli notification dlq replay <id>
```

##### `dotnext-cli notification dlq replay-all`

Re-publish every replayable DLQ entry to the main queue

```
dotnext-cli notification dlq replay-all
```

#### `dotnext-cli notification events`

Per-event notification toggles

```
dotnext-cli notification events
```

##### `dotnext-cli notification events list`

List notification event kinds and whether each is enabled

```
dotnext-cli notification events list
```

##### `dotnext-cli notification events set`

Enable or disable notifications for an event kind

```
dotnext-cli notification events set <event-kind> --enabled <true|false> [flags]
```

**Flags:**

- `--enabled` — true|false (required)

#### `dotnext-cli notification inbox`

Manage a user's notification bell inbox

```
dotnext-cli notification inbox
```

##### `dotnext-cli notification inbox clear`

Delete all of a user's notifications

```
dotnext-cli notification inbox clear --user <user> [flags]
```

**Flags:**

- `--force` — confirm clearing
- `--include-broadcast` — also remove broadcast rows (admin)
- `--user` — target user (required)

##### `dotnext-cli notification inbox list`

List a user's recent notifications

```
dotnext-cli notification inbox list --user <user> [flags]
```

**Flags:**

- `--limit` — max rows (1-200) (default `50`)
- `--user` — target user (email|username|id, required)

##### `dotnext-cli notification inbox read`

Mark one notification read

```
dotnext-cli notification inbox read <notification-id>
```

##### `dotnext-cli notification inbox read-all`

Mark all of a user's notifications read

```
dotnext-cli notification inbox read-all --user <user> [flags]
```

**Flags:**

- `--user` — target user (required)

### `dotnext-cli nspawn`

Manage SSH sandbox nspawn images

```
dotnext-cli nspawn
```

#### `dotnext-cli nspawn build`

debootstrap a deterministic, immutable nspawn rootfs

```
dotnext-cli nspawn build [flags]
```

**Flags:**

- `--codename` — image family (e.g. debian-13) (default `debian-13`)
- `--includes` — comma-separated debootstrap --include list (default `bash,coreutils,procps,findutils,grep,sed,gawk,less,nano,ca-certificates,git,curl,wget,vim-tiny,php-cli,php-mysql,php-curl,php-xml,php-mbstring,php-zip,php-gd,unzip,rsync,mariadb-client`)
- `--snapshot` — snapshot.debian.org timestamp YYYYMMDDTHHMMSSZ (mandatory)
- `--suite` — debootstrap suite (trixie, bookworm, ...) (default `trixie`)
- `--version` — image version label (e.g. v1, v2)

#### `dotnext-cli nspawn list`

List sealed nspawn images

```
dotnext-cli nspawn list
```

#### `dotnext-cli nspawn prune`

Remove sealed images that no user is pinned to

```
dotnext-cli nspawn prune [flags]
```

**Flags:**

- `--dry-run` — explicit dry-run (default; mutually exclusive with --yes)
- `--yes` — actually delete (default: dry-run)

### `dotnext-cli package`

Manage hosting packages

```
dotnext-cli package
```

#### `dotnext-cli package create`

Create a hosting package (direct DB)

```
dotnext-cli package create [flags]
```

**Flags:**

- `--backup-kinds` — CSV of allowed backup destination kinds (empty=none)
- `--backup-retention` — retention policy at the cap: reject or prune (empty=reject)
- `--bw-mb` — bandwidth quota in MB (0=unlimited) (default `0`)
- `--cgi` — enable CGI
- `--composer` — let users on this package run Composer for their domains; their composer.json scripts can run commands in their SSH sandbox, so grant it as you would SSH
- `--cpu` — CPU quota percent across all cores (0=unlimited) (default `0`)
- `--databases` — max databases (0=unlimited) (default `0`)
- `--disk-mb` — disk quota in MB (0=unlimited) (default `0`)
- `--docker-app-slugs` — CSV allowlist of catalog slugs tenants may install (empty=server default)
- `--domains` — max domains (0=unlimited) (default `0`)
- `--egress-icmp` — allow outbound ICMP echo-request (ping) for enforced tenants on this package
- `--egress-ssh-out` — allow outbound SSH (:22) for enforced tenants on this package
- `--egress-ssh-out-cidrs` — JSON array of CIDRs scoping outbound SSH (empty=anywhere), e.g. '["140.82.112.0/20"]'
- `--emails` — max email accounts (0=unlimited) (default `0`)
- `--fpm-advanced` — unlock tenant advanced FPM knobs (implies --fpm-user-can-edit)
- `--fpm-max-children` — FPM pm.max_children cap (0=default 20) (default `0`)
- `--fpm-user-can-edit` — let tenants pick an FPM performance mode
- `--fpm-version-defaults` — JSON map of per-version FPM defaults (empty={})
- `--fpm-worker-mem-mb` — FPM advisory per-worker memory budget in MB (0=default 64) (default `0`)
- `--git-deploy` — let users on this package deploy their domains from a Git repository; a deployment runs their own commands in their SSH sandbox, so grant it as you would SSH
- `--io-read-mbps` — disk read bandwidth limit in MB/s (0=unlimited) (default `0`)
- `--io-write-mbps` — disk write bandwidth limit in MB/s (0=unlimited) (default `0`)
- `--max-backup-schedules` — max tenant scheduled backups (0=default 1) (default `0`)
- `--max-backups` — tenant backup retention cap (0=backups not included) (default `0`)
- `--max-db-users` — max database users (0=unlimited) (default `0`)
- `--max-docker-apps` — max tenant docker apps (0=docker apps not included) (default `0`)
- `--max-ftp-accounts` — max tenant FTP/SFTP subaccounts (0=feature not included) (default `0`)
- `--max-python-apps` — max tenant python apps (0=python apps not included) (default `0`)
- `--max-tasks` — max processes/threads per user slice (0=unlimited) (default `0`)
- `--memory-mb` — memory limit in MB (0=unlimited) (default `0`)
- `--name` — package name (required)
- `--nspawn-image` — pin users to an nspawn rootfs version (empty=server default)
- `--php-disabled-functions` — PHP functions disabled for this package's sites, comma-separated. "default" = the command-exec lockdown (exec,passthru,shell_exec,system,proc_open,popen,pcntl_exec,pcntl_fork,proc_nice,dl); "none" = nothing disabled. Allowing system/exec/shell_exec/passthru/popen/proc_open/pcntl_exec also lifts the PHP Defense ban on it for these sites. Wins over --php-exec.
- `--php-exec` — opt out of the PHP command-exec lockdown (exec/proc_open work)
- `--php-settings-policy` — JSON object of php.ini directive to level (admin_only / tenant_allowed; tenant_privileged for security-sensitive directives) saying who may set it on the per-domain PHP Settings page. Empty = defaults (tenants may set every directive they can today). The CLI and admins are not limited by it. e.g. '{"memory_limit":"admin_only"}'
- `--redis-keys-per-call` — most keys one Redis key browser API call may list or delete for users on this package (100 to 5000, default 500) (default `500`)
- `--scheduled-backups` — allow tenant scheduled backups
- `--ssh` — enable SSH access
- `--webmail` — enable webmail (Bulwark UI) for this package (default `true`)
- `--webmail-branding` — let users on this package brand the webmail of their domains (mail.<domain>): its name, logos, favicon and links
- `--website-mail` — let the sites on this package send email with PHP mail()

#### `dotnext-cli package delete`

Delete a hosting package (direct DB)

```
dotnext-cli package delete <package-id> [flags]
```

**Flags:**

- `--force` — skip confirmation

#### `dotnext-cli package edit`

Edit a hosting package (direct DB)

```
dotnext-cli package edit <package-id> [flags]
```

**Flags:**

- `--backup-kinds` — CSV of allowed backup destination kinds
- `--backup-retention` — retention policy at the cap: reject or prune
- `--bw-mb` — bandwidth MB (default `0`)
- `--cgi` — CGI access (true/false)
- `--composer` — let users on this package run Composer for their domains; their composer.json scripts can run commands in their SSH sandbox, so grant it as you would SSH (true/false)
- `--cpu` — CPU quota percent (default `0`)
- `--databases` — max databases (default `0`)
- `--disk-mb` — disk quota MB (default `0`)
- `--docker-app-slugs` — CSV allowlist of catalog slugs tenants may install
- `--domains` — max domains (default `0`)
- `--egress-icmp` — allow outbound ICMP ping for enforced tenants (true/false)
- `--egress-ssh-out` — allow outbound SSH :22 for enforced tenants (true/false)
- `--egress-ssh-out-cidrs` — JSON array of CIDRs scoping outbound SSH (empty=anywhere)
- `--emails` — max emails (default `0`)
- `--fpm-advanced` — tenant advanced FPM knobs (true/false, true implies fpm-user-can-edit)
- `--fpm-max-children` — FPM pm.max_children cap (default `0`)
- `--fpm-user-can-edit` — tenant FPM performance mode (true/false)
- `--fpm-version-defaults` — JSON map of per-version FPM defaults
- `--fpm-worker-mem-mb` — FPM advisory per-worker memory budget MB (default `0`)
- `--git-deploy` — let users on this package deploy their domains from a Git repository; a deployment runs their own commands in their SSH sandbox, so grant it as you would SSH (true/false)
- `--io-read-mbps` — io read MB/s (default `0`)
- `--io-write-mbps` — io write MB/s (default `0`)
- `--max-backup-schedules` — max tenant scheduled backups (default `0`)
- `--max-backups` — tenant backup retention cap (default `0`)
- `--max-db-users` — max database users (default `0`)
- `--max-docker-apps` — max tenant docker apps (default `0`)
- `--max-ftp-accounts` — max tenant FTP/SFTP subaccounts (default `0`)
- `--max-python-apps` — max tenant python apps (default `0`)
- `--max-tasks` — max processes/threads (default `0`)
- `--memory-mb` — memory limit MB (default `0`)
- `--name` — package name
- `--nspawn-image` — pin users to an nspawn rootfs version (empty clears the pin)
- `--php-disabled-functions` — PHP functions disabled for this package's sites, comma-separated. "default" = the command-exec lockdown (exec,passthru,shell_exec,system,proc_open,popen,pcntl_exec,pcntl_fork,proc_nice,dl); "none" = nothing disabled. Allowing system/exec/shell_exec/passthru/popen/proc_open/pcntl_exec also lifts the PHP Defense ban on it for these sites. Wins over --php-exec.
- `--php-exec` — PHP command-exec opt-out (true/false)
- `--php-settings-policy` — JSON object of php.ini directive to level (admin_only / tenant_allowed; tenant_privileged for security-sensitive directives) saying who may set it on the per-domain PHP Settings page. Empty = defaults (tenants may set every directive they can today). The CLI and admins are not limited by it. e.g. '{"memory_limit":"admin_only"}'
- `--redis-keys-per-call` — most keys one Redis key browser API call may list or delete for users on this package (100 to 5000); unchanged unless given (default `0`)
- `--scheduled-backups` — tenant scheduled backups (true/false)
- `--ssh` — SSH access (true/false)
- `--webmail` — webmail Bulwark UI (true/false)
- `--webmail-branding` — let users on this package brand the webmail of their domains (mail.<domain>): its name, logos, favicon and links (true/false)
- `--website-mail` — sites may send email with PHP mail() (true/false)

#### `dotnext-cli package list`

List hosting packages (direct DB)

```
dotnext-cli package list
```

### `dotnext-cli page-template`

Manage error/index page templates

```
dotnext-cli page-template
```

#### `dotnext-cli page-template get`

Print a page template's content

```
dotnext-cli page-template get <key>
```

#### `dotnext-cli page-template list`

List page templates

```
dotnext-cli page-template list
```

#### `dotnext-cli page-template reset`

Reset a page template to the built-in default

```
dotnext-cli page-template reset <key>
```

#### `dotnext-cli page-template set`

Set a page template's content from a file

```
dotnext-cli page-template set <key> [flags]
```

**Flags:**

- `--file` — path to the template content (required)

### `dotnext-cli panel-cert`

Manage the panel's own TLS certificates (hostname + mail)

```
dotnext-cli panel-cert
```

#### `dotnext-cli panel-cert issue`

Request (re)issuance of a panel certificate on the next reconcile pass

```
dotnext-cli panel-cert issue [hostname|mail]
```

#### `dotnext-cli panel-cert new-key`

Replace a self-signed panel certificate with a new one, with a new key

```
dotnext-cli panel-cert new-key [hostname|mail]
```

#### `dotnext-cli panel-cert status`

Show both panel certificates (hostname + mail)

```
dotnext-cli panel-cert status
```

#### `dotnext-cli panel-cert toggle`

Enable/disable Let's Encrypt (or staging) for a panel certificate

```
dotnext-cli panel-cert toggle [hostname|mail] [flags]
```

**Flags:**

- `--le` — enable Let's Encrypt issuance (true|false)
- `--staging` — use the LE staging environment (true|false)

### `dotnext-cli panel-primary`

Manage the panel's primary mail domain row

```
dotnext-cli panel-primary
```

#### `dotnext-cli panel-primary ensure`

Ensure a panel-primary domain row exists for the given hostname

```
dotnext-cli panel-primary ensure [flags]
```

**Flags:**

- `--hostname` — panel hostname (e.g. dotnext-panel.local)

### `dotnext-cli pdns`

PowerDNS helpers (recursor forwarders, etc.)

```
dotnext-cli pdns
```

#### `dotnext-cli pdns backfill`

Converge /etc/powerdns/recursor.forwards with the panel database

```
dotnext-cli pdns backfill [flags]
```

**Flags:**

- `--dry-run` — explicit dry-run (default; mutually exclusive with --yes)
- `--no-confirm` — skip the y/N confirmation when --yes is used (for scripted runs; otherwise set DOTNEXT_PDNS_BACKFILL_NONINTERACTIVE=1)
- `--verbose` — print per-zone detail
- `--yes` — apply the plan (default is dry-run)

#### `dotnext-cli pdns dnssec`

Per-domain DNSSEC management

```
dotnext-cli pdns dnssec
```

##### `dotnext-cli pdns dnssec disable`

Disable DNSSEC (removes keys + rectifies)

```
dotnext-cli pdns dnssec disable <domain> [flags]
```

**Flags:**

- `--force` — skip confirmation

##### `dotnext-cli pdns dnssec ds`

Print DS records to publish at the parent registrar

```
dotnext-cli pdns dnssec ds <domain>
```

##### `dotnext-cli pdns dnssec enable`

Enable DNSSEC for a zone (creates KSK+ZSK, rectifies, persists keys)

```
dotnext-cli pdns dnssec enable <domain>
```

##### `dotnext-cli pdns dnssec status`

Show cached DNSSEC keys for a domain

```
dotnext-cli pdns dnssec status <domain>
```

### `dotnext-cli per-user-egress`

Per-user PHP-FPM egress firewall operator commands

```
dotnext-cli per-user-egress
```

#### `dotnext-cli per-user-egress approve`

Approve a pending egress request

```
dotnext-cli per-user-egress approve <request-id>
```

#### `dotnext-cli per-user-egress deny`

Deny a pending egress request

```
dotnext-cli per-user-egress deny <request-id>
```

#### `dotnext-cli per-user-egress flip-mature`

Flip mature LEARNING policies to ENFORCED

```
dotnext-cli per-user-egress flip-mature [flags]
```

**Flags:**

- `--dry-run` — show what would change without writing to DB
- `--soak-days` — minimum LEARNING age before auto-flip to ENFORCED (default `7`)

#### `dotnext-cli per-user-egress get`

Show a user's egress policy

```
dotnext-cli per-user-egress get <email-or-id>
```

#### `dotnext-cli per-user-egress requests`

List pending egress requests (the admin queue)

```
dotnext-cli per-user-egress requests
```

#### `dotnext-cli per-user-egress set-policy`

Set a user's egress state + allowed destinations (replaces the list)

```
dotnext-cli per-user-egress set-policy <email-or-id> [flags]
```

**Flags:**

- `--allow` — allowed destination CIDR[,PORT[,PROTO]] (repeatable; replaces the list) (default `[]`)
- `--state` — egress state: off|learning|enforced (required)

#### `dotnext-cli per-user-egress summary`

Show egress policy state counts + pending queue depth

```
dotnext-cli per-user-egress summary
```

### `dotnext-cli php`

PHP version + extension + per-user pool management

```
dotnext-cli php
```

#### `dotnext-cli php ext`

Manage PHP extensions (server-wide per PHP version)

```
dotnext-cli php ext
```

##### `dotnext-cli php ext disable`

Disable an installed extension via phpdismod

```
dotnext-cli php ext disable <ext> [flags]
```

**Flags:**

- `--version` — PHP version (required)

##### `dotnext-cli php ext enable`

Enable an installed extension via phpenmod

```
dotnext-cli php ext enable <ext> [flags]
```

**Flags:**

- `--version` — PHP version (required)

##### `dotnext-cli php ext install`

Install (apt) an extension package

```
dotnext-cli php ext install <ext> [flags]
```

**Flags:**

- `--version` — PHP version (required)

##### `dotnext-cli php ext list`

List PHP extensions and their installed/enabled state for a version

```
dotnext-cli php ext list [flags]
```

**Flags:**

- `--version` — PHP version (e.g. 8.4) (required)

##### `dotnext-cli php ext remove`

Remove (apt) an extension package

```
dotnext-cli php ext remove <ext> [flags]
```

**Flags:**

- `--version` — PHP version (required)

#### `dotnext-cli php pool`

Per-user PHP-FPM pool

```
dotnext-cli php pool
```

##### `dotnext-cli php pool create`

Create a PHP-FPM pool for a user

```
dotnext-cli php pool create [flags]
```

**Flags:**

- `--idle-timeout` — process idle timeout (seconds) (default `60`)
- `--php-version` — PHP version e.g. 8.3 (required)
- `--pm-max-children` — pm.max_children (default `20`)
- `--pm-mode` — ondemand|dynamic|static (default `ondemand`)
- `--user` — owner (id or username) (required)

##### `dotnext-cli php pool delete`

Delete a PHP-FPM pool and its ini overrides

```
dotnext-cli php pool delete <pool-id> [flags]
```

**Flags:**

- `--force` — confirm deletion

##### `dotnext-cli php pool get`

Show a user's PHP pool

```
dotnext-cli php pool get <user>
```

##### `dotnext-cli php pool ini`

Manage a pool's php.ini overrides

```
dotnext-cli php pool ini
```

###### `dotnext-cli php pool ini add`

Add an ini override to a pool

```
dotnext-cli php pool ini add <pool-id> [flags]
```

**Flags:**

- `--directive` — php.ini directive (required)
- `--kind` — value|flag (default `value`)
- `--value` — value (for kind=value) or on/off (for kind=flag)

###### `dotnext-cli php pool ini list`

List a pool's ini overrides

```
dotnext-cli php pool ini list <pool-id>
```

###### `dotnext-cli php pool ini remove`

Remove an ini override

```
dotnext-cli php pool ini remove <override-id>
```

###### `dotnext-cli php pool ini update`

Update an ini override's value

```
dotnext-cli php pool ini update <override-id> [flags]
```

**Flags:**

- `--value` — new value

##### `dotnext-cli php pool list`

List all PHP-FPM pools

```
dotnext-cli php pool list
```

##### `dotnext-cli php pool reapply-all`

Mark all active pools pending so the reconciler re-renders them from the current template

```
dotnext-cli php pool reapply-all
```

##### `dotnext-cli php pool set`

Set a user's PHP version (reconciler regenerates pool conf)

```
dotnext-cli php pool set [flags]
```

**Flags:**

- `--user` — user (id|email|username) (required)
- `--version` — PHP version e.g. 8.4 (required)

#### `dotnext-cli php version`

Manage installed PHP versions

```
dotnext-cli php version
```

##### `dotnext-cli php version install`

Install a PHP version (e.g. 8.4) — installs base + required extensions, starts php<v>-fpm

```
dotnext-cli php version install <version>
```

##### `dotnext-cli php version list`

List installed PHP versions

```
dotnext-cli php version list
```

##### `dotnext-cli php version reload`

Reload php<v>-fpm.service (zero-downtime SIGUSR2)

```
dotnext-cli php version reload <version>
```

### `dotnext-cli php-defense`

PHP-defense (Snuffleupagus) status / mode / rules / incidents

```
dotnext-cli php-defense
```

#### `dotnext-cli php-defense incidents`

List Snuffleupagus incidents

```
dotnext-cli php-defense incidents [flags]
```

**Flags:**

- `--limit` — max rows (default `50`)
- `--rule` — filter by rule name

#### `dotnext-cli php-defense incidents-clear`

Delete every recorded Snuffleupagus incident

```
dotnext-cli php-defense incidents-clear [flags]
```

**Flags:**

- `--yes` — confirm deleting the whole incident history

#### `dotnext-cli php-defense mode`

Set the Snuffleupagus mode (renders + applies active.rules)

```
dotnext-cli php-defense mode <off|simulation|enforce>
```

#### `dotnext-cli php-defense rule-toggle`

Enable or disable a Snuffleupagus rule (renders + applies)

```
dotnext-cli php-defense rule-toggle <rule-name> [flags]
```

**Flags:**

- `--disable` — disable the rule
- `--enable` — enable the rule
- `--reason` — optional reason for the override

#### `dotnext-cli php-defense rules`

List rule overrides + current mode (bundle catalog is GUI-only)

```
dotnext-cli php-defense rules
```

#### `dotnext-cli php-defense status`

Show Snuffleupagus status

```
dotnext-cli php-defense status
```

### `dotnext-cli python-app`

Manage Python and Node.js apps (admin-only)

```
dotnext-cli python-app
```

#### `dotnext-cli python-app create`

Create a Python or Node.js app (same validation + port/proxy allocation as the UI)

```
dotnext-cli python-app create [flags]
```

**Flags:**

- `--app-root` — app root under the owner's home (required)
- `--app-type` — wsgi|asgi; for --runtime node: node (default) or npm (default `wsgi`)
- `--base-uri` — mount path on the domain (default `/`)
- `--domain` — domain name to mount on (resolved to its id)
- `--domain-id` — domain id to mount on (or use --domain)
- `--entrypoint` — module:callable, e.g. myapp.wsgi:application (required unless --framework); for --runtime node: a script like server.js, or with --app-type npm a package.json script (default start)
- `--env` — KEY=VALUE (repeatable) (default `[]`)
- `--framework` — marketplace slug (e.g. django, flask, fastapi); derives app-type/entrypoint and scaffolds the starter
- `--name` — app name (required)
- `--node-version` — Node.js major for --runtime node (only 22) (default `22`)
- `--python-version` — Python version, e.g. 3.12 (required for --runtime python)
- `--runtime` — python|node (default `python`)

#### `dotnext-cli python-app delete`

Stop + remove an app (app files are kept)

```
dotnext-cli python-app delete <app-id>
```

#### `dotnext-cli python-app env`

View or update a Python app's environment variables

```
dotnext-cli python-app env
```

##### `dotnext-cli python-app env get`

List the app's environment variables (values are not masked — python_app_env has no secret flag)

```
dotnext-cli python-app env get <app-id>
```

##### `dotnext-cli python-app env set`

Set env vars (merge by default; --replace swaps the whole set)

```
dotnext-cli python-app env set <app-id> <KEY=VALUE> [<KEY=VALUE>...] [flags]
```

**Flags:**

- `--replace` — replace the entire env set instead of merging

#### `dotnext-cli python-app frameworks`

List the Python frameworks installable via --framework

```
dotnext-cli python-app frameworks
```

#### `dotnext-cli python-app list`

List Python and Node.js apps

```
dotnext-cli python-app list
```

#### `dotnext-cli python-app logs`

Show an app's recent logs

```
dotnext-cli python-app logs <app-id> [flags]
```

**Flags:**

- `--lines` — number of log lines (default `200`)

#### `dotnext-cli python-app rebuild`

Build an app again on the next pass (npm or pip) and restart it when done

```
dotnext-cli python-app rebuild <app-id>
```

#### `dotnext-cli python-app restart`

Restart an app

```
dotnext-cli python-app restart <app-id>
```

#### `dotnext-cli python-app start`

Start an app

```
dotnext-cli python-app start <app-id>
```

#### `dotnext-cli python-app stop`

Stop an app

```
dotnext-cli python-app stop <app-id>
```

### `dotnext-cli release`

Release-channel management (stable/development)

```
dotnext-cli release
```

#### `dotnext-cli release promote`

Promote a reviewed build to the stable channel (move the `stable` tag)

```
dotnext-cli release promote [<commit-ish>] [flags]
```

**Flags:**

- `--allow-source-build` — promote even when no release tarball is published for the commit (stable hosts then build it from source)

### `dotnext-cli repair`

Detect and fix known deployment-host issues

```
dotnext-cli repair [flags]
```

**Flags:**

- `--all` — Fix every issue including destructive ones
- `--apparmor-profiles-disabled` — Fix only: dotnext AppArmor profiles exist but are disabled
- `--apparmor-profiles-missing` — Fix only: dotnext AppArmor profiles absent from /etc/apparmor.d/
- `--auto` — Fix every non-destructive (safe) issue
- `--automation-443-include` — Fix only: panel-hostname :443 vhost missing the automation-API include (admin can't opt into API-on-443 for firewalled billing hosts)
- `--bulwark-jwt-secret` — Fix only: Bulwark webmail-SSO secret poisoned / out of sync with bulwark.env (mail impersonation 'Invalid signature')
- `--crowdsec-bouncer-key` — Fix only: crowdsec-firewall-bouncer crash-loops with stale LAPI key
- `--daemon-reload` — Fix only: systemd has unloaded unit-file changes on disk
- `--diagnose` — Report broken conditions without fixing
- `--dirty-ftp-264` — Fix only: schema half-applied migration 264 (FTP row size) — panel-api cannot start
- `--dirty-mailhostname-301` — Fix only: schema dirty at migration 301 (mail_hostname, interrupted update) — panel-api cannot start
- `--dirty-migration` — Fix only: database schema is dirty — panel-api cannot start
- `--docroot-www-data-group` — Fix only: web docroot files not group www-data / dirs not setgid (nginx 403 on newly uploaded media)
- `--etc-dotnext-perms` — Fix only: /etc/dotnext not traversable by hosting users (SSH/SFTP locked out — sandbox-mode unreadable)
- `--fpm-masters-down` — Fix only: a per-user PHP-FPM master is crash-looping (sites 500/502 with missing fpm.sock)
- `--git-ownership` — Fix only: /opt/dotnext/.git owned by wrong user
- `--git-pointer` — Fix only: /opt/dotnext/.git is a corrupted worktree pointer
- `--git-stale-worktrees` — Fix only: /opt/dotnext/.git/worktrees has stale entries
- `--nginx-config-invalid` — Fix only: dotnext-default/dotnext-api.conf has `http2 on;` on nginx<1.25.1 (nginx -t fails, reloads rejected)
- `--nginx-missing-includes` — Fix only: panel :8443 vhost includes a missing optional snippet (phpMyAdmin/Adminer) — nginx -t fails, nothing on 8443
- `--node-modules` — Fix only: panel-ui/node_modules partial (missing .bin/tsc)
- `--ondrej-nginx-ppa` — Fix only: stale ondrej/nginx PPA in apt sources (404 on noble)
- `--orphan-migration-staging` — Fix only: /var/lib/dotnext-migrations/* dirs for jobs already terminal in DB
- `--orphan-slices` — Fix only: dotnext-user-*.slice units exist for deleted unix users
- `--panel-csp-snippet` — Fix only: panel :8443 vhost includes a missing CSP snippet (dotnext-panel-csp.conf) — nginx -t fails, nothing on 8443
- `--uploads-dir` — Fix only: /var/lib/dotnext-uploads missing or wrong perms
- `--yes` — Skip interactive confirmation for destructive repairs

### `dotnext-cli retention-sweep`

Prune expired rows from unbounded log/report tables

```
dotnext-cli retention-sweep [flags]
```

**Flags:**

- `--dry-run` — count would-delete rows without removing them

### `dotnext-cli secrets`

Rotate panel secrets after remediation (operator ceremony, run as root)

```
dotnext-cli secrets
```

#### `dotnext-cli secrets rotate`

Rotate an exposed secret (see docs/secret-rotation.md)

```
dotnext-cli secrets rotate
```

##### `dotnext-cli secrets rotate all`

Rotate every built panel secret in a lockout-safe order

```
dotnext-cli secrets rotate all [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

##### `dotnext-cli secrets rotate db-app-user`

Rotate the panel DB app-user (dotnext_panel_app) password + DATABASE_URL

```
dotnext-cli secrets rotate db-app-user [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

##### `dotnext-cli secrets rotate jwt`

Rotate JWT_SECRET in panel.env (vestigial since the move to Kratos; safe near-noop)

```
dotnext-cli secrets rotate jwt [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

##### `dotnext-cli secrets rotate kratos`

Rotate the Kratos DB password and cookie/default secrets (signs every user out)

```
dotnext-cli secrets rotate kratos [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

##### `dotnext-cli secrets rotate pdns`

Rotate the PowerDNS DB user password (pdns.env + gmysql backend conf)

```
dotnext-cli secrets rotate pdns [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

##### `dotnext-cli secrets rotate redis-panel-token`

Rotate DOTNEXT_REDIS_PANEL_TOKEN (panel.env + redis aclfile, live ACL SETUSER)

```
dotnext-cli secrets rotate redis-panel-token [flags]
```

**Flags:**

- `--dry-run` — print the plan and touch nothing

### `dotnext-cli security`

Security advisories for this server

```
dotnext-cli security
```

#### `dotnext-cli security feed-check`

Act on new security advisories now (dotnext-security-feed.timer runs this hourly)

```
dotnext-cli security feed-check
```

### `dotnext-cli serve`

Start the dotNext Panel HTTP(S) server

```
dotnext-cli serve
```

### `dotnext-cli service`

Admin service control (start/stop/restart/reload/enable/disable)

```
dotnext-cli service
```

#### `dotnext-cli service action`

Run a service action on an allowlisted unit

```
dotnext-cli service action <name> <restart|start|stop|reload|enable|disable> [flags]
```

**Flags:**

- `--force` — confirm a disruptive action (stop/restart/disable)

#### `dotnext-cli service status`

Show service status details

```
dotnext-cli service status
```

### `dotnext-cli session`

Audit and revoke active login sessions (Kratos)

```
dotnext-cli session
```

#### `dotnext-cli session list`

List active login sessions (optionally filtered by --user email)

```
dotnext-cli session list [flags]
```

**Flags:**

- `--user` — only show sessions for this user email

#### `dotnext-cli session revoke`

Revoke a single login session

```
dotnext-cli session revoke <session-id>
```

#### `dotnext-cli session revoke-user`

Revoke ALL login sessions for a user (sign out everywhere)

```
dotnext-cli session revoke-user <email-or-id>
```

### `dotnext-cli settings`

Inspect and patch server settings (headless equivalent of /admin/settings)

```
dotnext-cli settings
```

#### `dotnext-cli settings get`

Print current server settings (JSON or key=value)

```
dotnext-cli settings get
```

#### `dotnext-cli settings keys`

List the settable keys for `settings set`

```
dotnext-cli settings keys
```

#### `dotnext-cli settings mail-hostname`

Show the panel mail hostname, or request, cancel or follow a change

```
dotnext-cli settings mail-hostname [flags]
```

**Flags:**

- `--applied` — print only a custom applied mail hostname; print nothing when mail.<hostname> is in effect
- `--cancel` — withdraw a pending or failed mail hostname change
- `--set` — request a change of the panel mail hostname to this name
- `--status` — print the progress of a requested mail hostname change

#### `dotnext-cli settings set`

Patch one or more server settings (same validation + side effects as the UI)

```
dotnext-cli settings set <key>=<value> [<key>=<value>...]
```

### `dotnext-cli shared-resource`

Manage shared mail resources — calendars, contacts, files

```
dotnext-cli shared-resource
```

#### `dotnext-cli shared-resource create`

Create a shared resource

```
dotnext-cli shared-resource create [flags]
```

**Flags:**

- `--display-name` — display name
- `--domain` — domain ID (ULID)
- `--kind` — calendar|addressbook|files (for a shared mailbox, create a mail group)
- `--name` — host address local part

#### `dotnext-cli shared-resource grant`

Add or update a grant (upsert into the grant set)

```
dotnext-cli shared-resource grant [flags]
```

**Flags:**

- `--grantee` — grantee mailbox/group ID
- `--grantee-kind` — mailbox|group (default `mailbox`)
- `--resource` — shared resource ID
- `--rights` — read|readwrite|admin, or readdelete (files only: read and delete) (default `read`)

#### `dotnext-cli shared-resource grants`

List grants on a shared resource

```
dotnext-cli shared-resource grants [flags]
```

**Flags:**

- `--resource` — shared resource ID

#### `dotnext-cli shared-resource list`

List shared resources in a domain

```
dotnext-cli shared-resource list [flags]
```

**Flags:**

- `--domain` — domain ID (ULID)

#### `dotnext-cli shared-resource remove`

Delete a shared resource (reconciler tears down the host principal)

```
dotnext-cli shared-resource remove [flags]
```

**Flags:**

- `--resource` — shared resource ID

#### `dotnext-cli shared-resource revoke`

Remove a grantee's grant from a shared resource

```
dotnext-cli shared-resource revoke [flags]
```

**Flags:**

- `--grantee` — grantee ID to revoke
- `--resource` — shared resource ID

#### `dotnext-cli shared-resource scanner`

A Files resource's scanner: an FTP/SFTP drop folder moved into the shared folder

```
dotnext-cli shared-resource scanner
```

##### `dotnext-cli shared-resource scanner disable`

Turn a Files resource's scanner off: remove its login (files left in the drop folder stay)

```
dotnext-cli shared-resource scanner disable [flags]
```

**Flags:**

- `--resource` — shared resource ID (a Files resource)

##### `dotnext-cli shared-resource scanner enable`

Turn a Files resource's scanner on: make its drop folder and FTP/SFTP login

```
dotnext-cli shared-resource scanner enable [flags]
```

**Flags:**

- `--login` — login name; the login is <owner>_<login> (default `scanner`)
- `--password` — login password (generated and shown once if omitted)
- `--resource` — shared resource ID (a Files resource)
- `--space-mb` — space the scans may keep, in MB (the login's disk quota where logins are isolated) (default `1024`)

##### `dotnext-cli shared-resource scanner status`

Show a Files resource's scanner

```
dotnext-cli shared-resource scanner status [flags]
```

**Flags:**

- `--resource` — shared resource ID (a Files resource)

### `dotnext-cli ssh-key`

Manage user SSH authorized keys

```
dotnext-cli ssh-key
```

#### `dotnext-cli ssh-key add`

Add an SSH public key for a user

```
dotnext-cli ssh-key add [flags]
```

**Flags:**

- `--name` — key label (required)
- `--pub-key` — raw public key (e.g. 'ssh-ed25519 AAAA... user@host')
- `--pub-key-file` — path to public key file
- `--pub-key-stdin` — read public key from stdin
- `--user` — user (id|email|username) (required)

#### `dotnext-cli ssh-key delete`

Delete an SSH key

```
dotnext-cli ssh-key delete <key-id> [flags]
```

**Flags:**

- `--force` — skip confirmation

#### `dotnext-cli ssh-key list`

List SSH keys (filtered by --user, or all without filter)

```
dotnext-cli ssh-key list [flags]
```

**Flags:**

- `--all` — list every user's keys (default when --user is omitted)
- `--user` — filter by user (id|email|username); empty = all

### `dotnext-cli ssl`

Manage Let's Encrypt SSL certificates

```
dotnext-cli ssl
```

#### `dotnext-cli ssl disable`

Disable SSL for a domain (reconciler will revoke cert)

```
dotnext-cli ssl disable <domain>
```

#### `dotnext-cli ssl enable`

Enable SSL for a domain and wait until the vhost serves the real certificate

```
dotnext-cli ssl enable <domain> [flags]
```

**Flags:**

- `--nginx-dir` — directory holding the enabled nginx vhosts (default `/etc/nginx/sites-enabled`)
- `--no-wait` — return as soon as the domain is marked, without waiting for the certificate
- `--wait-timeout` — how long to wait for the vhost to serve the real certificate (default `3m0s`)

#### `dotnext-cli ssl list`

List SSL certificates (optionally filtered by user)

```
dotnext-cli ssl list [flags]
```

**Flags:**

- `--user` — filter by user (id|email|username)

#### `dotnext-cli ssl readiness`

Report domains whose origin certificate would fail Cloudflare Full (strict)

```
dotnext-cli ssl readiness [flags]
```

**Flags:**

- `--all` — list every domain, not only the ones that would fail
- `--nginx-dir` — directory holding the enabled nginx vhosts (default `/etc/nginx/sites-enabled`)

#### `dotnext-cli ssl renew`

Renew SSL cert via certbot (synchronous, calls agent)

```
dotnext-cli ssl renew <domain> [flags]
```

**Flags:**

- `--force` — force renewal even if cert is not due

#### `dotnext-cli ssl retry`

Reset a stuck cert (failed / pending_acme_retry) and re-attempt ACME issuance now

```
dotnext-cli ssl retry <domain>
```

#### `dotnext-cli ssl set-custom`

Install an operator-supplied SSL cert + key

```
dotnext-cli ssl set-custom [flags]
```

**Flags:**

- `--cert` — path to the certificate PEM file (leaf + chain) (required)
- `--domain` — domain name or id (required)
- `--key` — path to the private key PEM file (required)

#### `dotnext-cli ssl shared`

Manage shared wildcard/multi-SAN certificates (upload once, serve many domains)

```
dotnext-cli ssl shared
```

##### `dotnext-cli ssl shared attach`

Attach a domain to a shared certificate (must cover the domain)

```
dotnext-cli ssl shared attach --domain <name> --cert-id <id> [flags]
```

**Flags:**

- `--cert-id` — shared certificate id
- `--domain` — domain name

##### `dotnext-cli ssl shared delete`

Delete a shared certificate (fails if domains are still attached)

```
dotnext-cli ssl shared delete --id <id> [flags]
```

**Flags:**

- `--id` — shared certificate id

##### `dotnext-cli ssl shared detach`

Detach a domain from its shared certificate (reverts to LE)

```
dotnext-cli ssl shared detach --domain <name> [flags]
```

**Flags:**

- `--domain` — domain name

##### `dotnext-cli ssl shared list`

List shared certificates + their attached-domain counts

```
dotnext-cli ssl shared list
```

##### `dotnext-cli ssl shared upload`

Upload a server-wide shared certificate (agent validates + writes it)

```
dotnext-cli ssl shared upload --name <name> --cert <fullchain.pem> --key <privkey.pem> [flags]
```

**Flags:**

- `--cert` — path to fullchain PEM
- `--key` — path to private key PEM
- `--name` — human-readable name

### `dotnext-cli sso`

SSO (Single Sign-On) management commands

```
dotnext-cli sso
```

#### `dotnext-cli sso prune-tokens`

Manually purge expired SSO tokens

```
dotnext-cli sso prune-tokens
```

#### `dotnext-cli sso rotate-key`

Rotate the SSO encryption key

```
dotnext-cli sso rotate-key [flags]
```

**Flags:**

- `--current-key` — path to current encryption key (default `/etc/dotnext/sso_key.txt`)
- `--new-key` — path to new encryption key (required)

### `dotnext-cli sso-reap`

Sweep stranded dotnext-sso-<nonce>.php files

```
dotnext-cli sso-reap
```

### `dotnext-cli system`

System information and services

```
dotnext-cli system
```

#### `dotnext-cli system diagnostic`

Generate a redacted, encrypted diagnostic bundle for support handoff

```
dotnext-cli system diagnostic [flags]
```

**Flags:**

- `--json` — print the raw JSON result (for scripting)

#### `dotnext-cli system info`

Show system info (hostname, uptime, CPU, memory, disk)

```
dotnext-cli system info
```

#### `dotnext-cli system process`

List / kill host processes (admin)

```
dotnext-cli system process
```

##### `dotnext-cli system process kill`

Kill a process (TERM, or KILL with --force)

```
dotnext-cli system process kill <pid> [flags]
```

**Flags:**

- `--force` — confirm the kill (and use SIGKILL)

##### `dotnext-cli system process list`

List host processes

```
dotnext-cli system process list
```

#### `dotnext-cli system resolver`

Show / set system DNS resolvers

```
dotnext-cli system resolver
```

##### `dotnext-cli system resolver get`

Show current resolver source + addresses

```
dotnext-cli system resolver get
```

##### `dotnext-cli system resolver set`

Set resolver addresses (validated; same rules as the UI)

```
dotnext-cli system resolver set [flags]
```

**Flags:**

- `--addr` — resolver IP (repeatable, IPv4 or IPv6) (default `[]`)
- `--search-domain` — optional search domain

#### `dotnext-cli system restore`

Restore a system backup onto this host (CLI)

```
dotnext-cli system restore [flags]
```

**Flags:**

- `--apply` — after staging, apply selected stages onto live host (default true) (default `true`)
- `--apply-stage` — restrict apply to named stages (repeatable). Empty = panel_db + panel_config + tls (the safe defaults) (default `[]`)
- `--credentials-ref` — absolute path to env file with backend creds (root:root 0600)
- `--force` — required — restore overwrites the running panel
- `--from-tar` — restore the SYSTEM leg from a downloaded Full Server container (.tar) instead of a restic repo. Applies panel_db + panel_config + tls by default
- `--include-accounts` — also restore each linked account
- `--interactive` — force interactive prompts even when --remote-url is set
- `--password` — restic password (literal; overrides --password-file). Avoid in shell history; prefer --interactive
- `--password-file` — restic password file (default: /etc/dotnext/restic-repo.password)
- `--remote-url` — restic repo URL or local path (e.g. sftp:user@host:/path)
- `--sftp-auth` — SFTP auth mode: key | password (default: ssh config)
- `--sftp-key` — absolute path to the SSH private key for SFTP
- `--sftp-port` — SFTP port when not 22 (default `0`)
- `--snapshot` — system_manifest snapshot ID, or 'latest' to auto-pick newest

#### `dotnext-cli system services`

Show systemd service status

```
dotnext-cli system services
```

#### `dotnext-cli system ssh-keys`

System SSH keys for backup destinations (list/generate)

```
dotnext-cli system ssh-keys
```

##### `dotnext-cli system ssh-keys generate`

Generate a system SSH keypair (for sftp backup destinations)

```
dotnext-cli system ssh-keys generate [flags]
```

**Flags:**

- `--name` — key name (required)
- `--type` — key type (ed25519|rsa) (default `ed25519`)

##### `dotnext-cli system ssh-keys list`

List system SSH keys

```
dotnext-cli system ssh-keys list
```

### `dotnext-cli ufw`

UFW utilities (port baseline only; IP decisions live in CrowdSec)

```
dotnext-cli ufw
```

#### `dotnext-cli ufw default`

Set the default policy for a chain

```
dotnext-cli ufw default <incoming|outgoing> <allow|deny|reject> [flags]
```

**Flags:**

- `--force` — confirm a lockout-prone default policy

#### `dotnext-cli ufw disable`

Disable UFW (lockout risk)

```
dotnext-cli ufw disable [flags]
```

**Flags:**

- `--force` — confirm this destructive action

#### `dotnext-cli ufw enable`

Enable UFW

```
dotnext-cli ufw enable [flags]
```

**Flags:**

- `--force` — confirm this destructive action

#### `dotnext-cli ufw migrate-ip-bans`

Migrate UFW `from <IP>` deny rules to CrowdSec decisions

```
dotnext-cli ufw migrate-ip-bans [flags]
```

**Flags:**

- `--dry-run` — Show what would migrate; make no changes
- `--no-cdn` — Confirm panel is not behind a CDN; bypasses the trusted_ips hard guard
- `--revert` — Restore UFW rules from snapshot and remove matching CrowdSec decisions
- `--yes` — Required for any destructive operation (migrate or revert)

#### `dotnext-cli ufw rule`

Add or delete UFW rules

```
dotnext-cli ufw rule
```

##### `dotnext-cli ufw rule add`

Add a UFW rule

```
dotnext-cli ufw rule add [flags]
```

**Flags:**

- `--action` — allow|deny|reject|limit (default `allow`)
- `--from` — source address/CIDR (optional)
- `--port` — port or app name (e.g. 443, OpenSSH) (required)
- `--proto` — tcp|udp (optional)

##### `dotnext-cli ufw rule delete`

Delete a UFW rule by its numbered position (see `ufw status`)

```
dotnext-cli ufw rule delete <num>
```

#### `dotnext-cli ufw status`

Show UFW status and numbered rules

```
dotnext-cli ufw status
```

### `dotnext-cli update`

Pull latest code, rebuild, migrate, and restart services

```
dotnext-cli update [flags]
```

**Flags:**

- `--force` — Run the full rebuild/restart cycle even when git pull found no new commits
- `--from-source` — Build binaries on this host instead of downloading the release tarball from Gitea Releases. Default is to download the tarball (90s update vs 5-10min source build). Use --from-source when offline, on a private fork, or to test uncommitted changes.

#### `dotnext-cli update apt-refresh`

Refresh the Updates Center OS-patch state now (used after an auto-apply)

```
dotnext-cli update apt-refresh
```

### `dotnext-cli user`

Manage panel users

```
dotnext-cli user
```

#### `dotnext-cli user 2fa-reset`

Strip TOTP + recovery codes from a user (CLI escape hatch when locked out)

```
dotnext-cli user 2fa-reset <email|username|user-id>
```

#### `dotnext-cli user create`

Create a new user (direct DB + Kratos; bypasses HTTP auth)

```
dotnext-cli user create [username] [flags]
```

**Flags:**

- `--admin` — grant admin role
- `--email` — user email (optional; a placeholder is synthesized from the username if omitted)
- `--name-first` — first name
- `--name-last` — last name
- `--password` — user password (required, min 10 chars)
- `--password-stdin` — read password from stdin (no prompt, no echo)
- `--username` — login username — the identifier users sign in with (required unless --email is given to derive one)

#### `dotnext-cli user delete`

Delete a user — destructive: domains, databases, mailboxes, OS account, /home, all related rows.

```
dotnext-cli user delete <email|username|user-id> [flags]
```

**Flags:**

- `--force` — skip confirmation prompt

#### `dotnext-cli user list`

List all users (direct DB)

```
dotnext-cli user list
```

#### `dotnext-cli user password`

Reset a user's password (auto-generates one if --password is omitted)

```
dotnext-cli user password <email|username|user-id> [flags]
```

**Flags:**

- `--expires-in` — TTL for recovery link (only with --link) (default `24h`)
- `--link` — emit a one-click recovery URL instead of setting the password directly
- `--password` — explicit new password (omit to auto-generate)
- `--password-stdin` — read new password from stdin (no prompt, no echo)

#### `dotnext-cli user rename`

Rename a tenant's Linux/login username in place

```
dotnext-cli user rename <email|username|user-id> <new-username> [flags]
```

**Flags:**

- `--yes` — skip the confirmation prompt

#### `dotnext-cli user reprovision`

Re-run OS provisioning for a user + sync the panel password

```
dotnext-cli user reprovision <email|username|user-id> [flags]
```

**Flags:**

- `--password` — new password (min 10 chars; omit to auto-generate)
- `--password-stdin` — read password from stdin (no prompt, no echo)

#### `dotnext-cli user ssh-forwarding`

Opt a hosting user into/out of SSH TCP forwarding for VS Code Remote-SSH

```
dotnext-cli user ssh-forwarding <email|username|user-id> <on|off>
```

#### `dotnext-cli user suspend`

Suspend a user (same cascade as the admin GUI/API)

```
dotnext-cli user suspend <id> [flags]
```

**Flags:**

- `--reason` — operator-facing suspend reason (shown in the admin user list)

#### `dotnext-cli user unsuspend`

Unsuspend a user (reverse the cascade)

```
dotnext-cli user unsuspend <id>
```

### `dotnext-cli user-token`

Manage a user's API tokens (mint / list / revoke)

```
dotnext-cli user-token
```

#### `dotnext-cli user-token list`

List a user's API tokens (active + revoked)

```
dotnext-cli user-token list <user-email|username|id>
```

#### `dotnext-cli user-token mint`

Mint a new API token for a user (secret shown once)

```
dotnext-cli user-token mint <user-email|username|id> --name <name> [--scope read:dns ...] [flags]
```

**Flags:**

- `--expires-in` — token lifetime (e.g. 720h); 0 = no expiry; max 8760h (default `0s`)
- `--name` — token name (required)
- `--scope` — scope to grant (repeatable; empty = full owner access) (default `[]`)

#### `dotnext-cli user-token revoke`

Revoke a user's API token

```
dotnext-cli user-token revoke <token-id>
```

### `dotnext-cli version`

Print dotNext version, commit, and runtime info

```
dotnext-cli version
```

