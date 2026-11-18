# Operations Runbook

Day-2 operator reference. For day-1 install see [installation.md](./installation.md).

## Standard daily / weekly cadence

| When | Action |
|---|---|
| Daily | Review `/dotnext-admin/security` → CrowdSec → Decisions. Spot any allowlist additions you need to make. |
| Daily | Review `/dotnext-admin/notifications` (or the bell) for `cert.renew.fail`, `backup.fail` and `service.down`. |
| Weekly | `dotnext-cli update` — pulls latest fixes; auto-applies migrations. |
| Weekly | Inspect `/dotnext-admin/server-status` → trends. Spot disks filling up before quota alarms fire. |
| Monthly | Verify a restore round-trip from at least one backup destination. |

## Common operator tasks

### Lock out a panel user

```bash
dotnext-cli user suspend <id> --reason "…"
dotnext-cli user unsuspend <id>
```

(The same cascade as Suspend in the admin UI.)

### Force-reset a user's password / 2FA

```bash
dotnext-cli user password <email>            # generates new password, prints once
dotnext-cli user 2fa-reset <email>           # strip TOTP + recovery codes (CLI escape hatch)
```

### Get an admin back in

```bash
dotnext-cli user 2fa-reset <email|username|user-id>    # when the admin lost their 2FA device
dotnext-cli user password <email|username|user-id> --link
```

`--link` prints a one-click recovery URL, valid 24 hours.

### Rotate the DB root password

`/dotnext-admin/settings` → Database → Root Password → Rotate.

Or CLI:

```bash
dotnext-cli db root-password
```

Audited (success + failure), announced via notifications.

### Move a domain between users

A first-class **Change owner** action reassigns a domain to a different user;
the vhost, DNS, and mail follow the new owner on the next reconcile. The
domain's Python apps, Git repository settings and Composer project folder
setting belong to the old owner and are deleted. A Docker app of the old owner
is unlinked from the domain and keeps its data; a server-level one stays. The
admin can also **rename a tenant's username in place** — the panel
re-prefixes that user's MariaDB databases, DB users, and PostgreSQL shadow roles
to match.

### Add a new IPv4 to the pool

1. Bring the IP up on the host (`ip addr add <ip>/<prefix> dev <iface>`).
2. `/dotnext-admin/ips` → Add.
3. Assign per-domain via Domains → Edit → Listen IP.

### Verify the installer URL

The install command is served straight from the GitHub repository (there is no
short-link domain yet), so the check is against the raw file:

```bash
curl -fsSL https://raw.githubusercontent.com/dotnext71/panel/main/install.sh | head -5   # expect the bash shebang + installer header
```

`bootstrap.sh` (same directory) is the TUI path and pulls the signed release
tarball. `deploy/get-installer/` holds the DNS/redirect setup to wire up a short
`get.` URL later; until then, document both raw URLs everywhere.

### Add a new PHP version

1. `apt install php8.x-fpm php8.x-cli php8.x-mbstring …`
2. The panel detects it on the next reconciler tick; appears in `/dotnext-admin/php-pools`.
3. Domains can pick it under Edit → PHP Version.

## Common failures

| Symptom | First thing to check |
|---|---|
| Panel UI says "502 Bad Gateway" | `systemctl status dotnext-api`; `journalctl -u dotnext-api -f` |
| Cert never issues | `dotnext-cli ssl list` → see the row's last error; `journalctl -u dotnext-agent -f` while you toggle SSL again |
| `Dirty database version N` on boot | A failed migration. `dotnext-cli migrate up` to retry; if it's deterministic, file an issue with the journal output. |
| `nginx -t` fails after a domain edit | Missing FastCGI cache keyzone (if you enabled cache before `dotnext-cli update` landed the keyzone conf). Run `dotnext-cli update` once. |
| Mail accounts work in IMAP but not SMTP submission | Bulwark may be down. `systemctl status bulwark`; check `:587` listener. |
| User app suddenly can't reach external API | Per-user egress firewall may have dropped the destination. Users → Edit → Egress. |

## Repair

```bash
dotnext-cli repair --diagnose       # default: report only
dotnext-cli repair --auto           # safe auto-fixes (perms, drop-in regen, no destructive)
dotnext-cli repair --all --yes      # destructive (wipes nginx tmp, regen pool sockets) — only when --diagnose tells you to
```

`dotnext-cli repair` covers 7 detectors, grouped by destructiveness.

**Updates → Repair Center → Auto-Repair safe issues** runs `dotnext-cli repair --auto`. That includes `docroot-www-data-group`, which puts docroot files back in group www-data and sets setgid on their directories, so nginx can serve newly uploaded media. A broken row that Auto-Repair leaves says why. A destructive one shows its `dotnext-cli repair --<id> --yes` command. One with no fix (`dirty-migration`, `fpm-masters-down`) needs your decision. In `--diagnose` output these carry `[destructive]` and `[manual]` after the label.

## Migration cleanup

After a successful cPanel / DA / Hestia restore: `dotnext-cli domain prune-orphans` lists sites in nginx `sites-enabled` that have no panel row; `dotnext-cli domain prune-orphans --apply` deletes them.

`--apply` lists `sites-enabled` again after the deletes. Each site shows as `deleted`, or as `FAILED` with the reason when it is still there, and the command exits non-zero if any site is left. With `--json`, the result has `deleted` and `failed` (site and reason) lists. **Updates → Domain repair → Delete orphans** runs the same command and lists any site left, with its reason.

## Diag bundle

`/dotnext-admin/support` → encrypted bundle, `mailto:` send. See [support.md](./support.md).
