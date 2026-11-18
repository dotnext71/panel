# Updates

`/dotnext-admin/updates`.

## What `dotnext-cli update` does

1. `git fetch origin main`
2. `git reset --hard origin/main` (the VM is treated as a deploy target; local edits get overwritten — see "VM local mods byte-identical to origin" scar for the rationale).
3. Rebuild `dotnext-api-api` and `dotnext-agent` (Go).
4. `npm ci && npm run build` for `frontend` (Vite).
5. Run pending DB migrations.
6. Re-apply install-time host configs that the agent owns (nginx drop-ins, php-fpm pool dirs, AppSec rules, FastCGI cache keyzone). This is the "install.sh is truth, not runbook" rule — anything required for a fresh host to work right is re-applied on every update.
7. Reload nginx → restart panel-api → restart agent → restart Bulwark, in that order.
8. Print success or the first failure line.

## How the UI runs it

Updates from the panel UI fire as **transient systemd units** (`systemd-run --transient --unit=dotnext-update-<timestamp>`). This means:

- The update survives the panel restarting itself mid-update (step 7).
- You can watch the live log from a fresh page load even if the browser disconnected.
- `journalctl -u dotnext-update-<timestamp>` is the audit trail.

The transient-unit survival was live-verified on 192.168.100.150.

## CLI

```bash
dotnext-cli update                 # download the release, migrate and restart
dotnext-cli update --force         # run the full cycle even when there is no new release
dotnext-cli update --from-source   # build on this host instead of downloading the release
```

If the update fails, `dotnext-cli update` prints a hint pointing at `dotnext-cli repair --diagnose` (after a failed update, repair is usually the next step).

## OS security auto-updates

Separate from `dotnext-cli update` (which updates the panel), **OS security
auto-updates are on by default** via `unattended-upgrades`, so security patches
for the base system land without operator action. The Update Center reports the
OS-update state (pending / applied) alongside the panel version.

## Stuck update tasks are reaped

An update run is tracked as a DB row and sealed success/failed when it finishes.
A run whose status can't be read — the agent is down or too old to answer, or a
transient unit hangs — is now **backstopped** so it can't spin "running" forever
in the Tasks indicator: any run older than 2 h is reaped as failed,
and a pollable run whose status call errors is reaped once clearly stale.

## What it does *not* update

- The kernel, system packages outside DotNext's drop-ins (use `apt update && apt full-upgrade`).
- PHP itself (use `apt`).
- MariaDB / PostgreSQL major version (manual; both data and config require operator decisions).

## Frequency

Panel self-update ships **off**. Turn it on under **Updates → DotNext panel
self-update** (`/dotnext-admin/updates`) and pick a time (default 04:30). The
panel then keeps `dotnext-autoupdate.timer` in step with that setting, and the
box runs `dotnext-cli update -f` daily from its release channel. Until you turn it on,
run `dotnext-cli update` yourself. (OS *security* patches auto-apply as above.) See
[Server Updates](./admin/server-updates.md) for why the stable channel suits
unattended updates.
