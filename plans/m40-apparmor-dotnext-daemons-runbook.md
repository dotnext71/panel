# M40 — AppArmor dotnext daemons runbook

Companion to plans/m40-apparmor-dotnext-daemons.md and ADR-0086.

## What it is, in one sentence

Path-based MAC profiles confine dotnext daemons (`dotnext-panel`,
`dotnext-agent`, `dotnext-bulwark`) plus critical system services
(`stalwart-mail`, `dotnext-kratos`) to a minimal file/cap/socket
surface. New profiles ship in **complain** mode for a 7-day burn-in
soak; operator flips per-profile to **enforce** after review.

## Files / units / commands

| Where | What |
|---|---|
| `/etc/apparmor.d/usr.local.bin.dotnext-*` | Loaded profile files. install.sh is the only writer; idempotent re-render on `dotnext update`. |
| `install/apparmor/usr.local.bin.dotnext-*` | Source profiles in the repo. |
| `/etc/dotnext/.apparmor-installed` | First-install marker (UTC timestamp). |
| `/etc/dotnext/.apparmor-grub-pending` | Sentinel — operator reboot required to activate AppArmor LSM. |
| `/etc/dotnext/.apparmor-disabled` | Sentinel — kernel doesn't expose AppArmor. |
| `aa-status` | Loaded profiles + per-profile mode. |
| `aa-status --json` | Machine-readable form (panel-agent uses this). |
| `aa-enforce <profile>` / `aa-complain <profile>` | Per-profile mode toggle. |
| `dotnext apparmor status` | dotnext profiles + modes (CLI). |
| `dotnext apparmor flip-mature [--profile X] [--dry-run]` | Bulk or per-profile flip. |

## Daily checks

```bash
aa-status | grep -E "dotnext-|stalwart-mail"
# Expected: 4+ profiles in complain or enforce

dotnext apparmor status
# Same data, dotnext-namespaced

journalctl -k --since today | grep 'apparmor="DENIED"' | grep -E "dotnext-|stalwart-mail"
# Expected: 0 lines on a stable host. Anything > 0 = profile drift
# (a daemon learned a new path between releases) — fix before next
# enforce flip.
```

## First-install workflow

Day 0 — fresh install:
- install.sh writes profiles, defaults all to complain.
- `aa-status | grep -c dotnext-` ≥ 3.
- `dotnext apparmor status` lists each profile in `complain` mode.

Days 1-7 — soak:
- Run a representative workload: create users, install WordPress,
  send mail, schedule a backup, scheduled malware scan, dotnext
  update, ssh in, etc.
- Tail `journalctl -k --since "1 hour ago" | grep apparmor` daily.
- Any `apparmor="ALLOWED"` line under complain mode is a would-deny
  in enforce mode — investigate the missing path/cap, edit the
  profile in `install/apparmor/`, commit + `dotnext update` (or copy
  to host + `apparmor_parser -r /etc/apparmor.d/<file>`).

Day 7 — flip:
```bash
dotnext apparmor flip-mature --dry-run
# Review the candidates list

dotnext apparmor flip-mature --profile dotnext-bulwark
# Start with smallest-surface daemon

# Day 8: dotnext-panel
dotnext apparmor flip-mature --profile dotnext-panel

# Day 9: dotnext-agent (widest cap set, flip last)
dotnext apparmor flip-mature --profile dotnext-agent

# Day 10: system services
dotnext apparmor flip-mature --profile stalwart-mail
dotnext apparmor flip-mature --profile dotnext-kratos

# Or, after independent review of all four:
dotnext apparmor flip-mature
```

## Adding a new path to a profile

When a daemon update needs a path not in the current profile:

1. Tail complain-mode AVC denials:
   ```bash
   journalctl -k --since "10 minutes ago" | grep "profile=\"dotnext-agent\""
   ```
2. Identify the path/cap from the denied entry.
3. Edit `install/apparmor/usr.local.bin.dotnext-agent` — add the path
   with the right mode (`r`, `rw`, `mr`, `ix` for exec children).
4. Apply locally for testing:
   ```bash
   apparmor_parser -r /etc/apparmor.d/usr.local.bin.dotnext-agent
   ```
5. Commit + `dotnext update` to ship.

## Troubleshooting "service won't start after enforce flip"

```bash
# Symptom: dotnext-panel in failed state after flip.
journalctl -u dotnext-panel --since "5 minutes ago" | tail -20
# Look for an early-boot file open / cap / socket failure.

# Revert the profile while you fix it:
aa-complain /etc/apparmor.d/usr.local.bin.dotnext-panel
systemctl restart dotnext-panel

# Tail complain-mode AVC to find the missing path:
journalctl -k --since "1 minute ago" | grep "dotnext-panel"

# Edit profile, reload, re-flip:
$EDITOR /etc/apparmor.d/usr.local.bin.dotnext-panel
apparmor_parser -r /etc/apparmor.d/usr.local.bin.dotnext-panel
aa-enforce /etc/apparmor.d/usr.local.bin.dotnext-panel

# Then propagate the fix back to install/apparmor/ and commit.
```

## Rollback levels

1. **Per-profile complain.** Single command:
   ```bash
   aa-complain /etc/apparmor.d/usr.local.bin.dotnext-<name>
   ```
   Profile stays loaded; would-deny logs continue but nothing is
   enforced.

2. **All dotnext profiles to complain.**
   ```bash
   for p in /etc/apparmor.d/usr.local.bin.dotnext-*; do aa-complain "$p"; done
   ```

3. **Unload all dotnext profiles.**
   ```bash
   for p in /etc/apparmor.d/usr.local.bin.dotnext-*; do
     apparmor_parser -R "$p"
     mv "$p" "$p.disabled"
   done
   ```
   `dotnext update` will re-render the files; rename `.disabled` again
   if you want the disable to persist.

4. **AppArmor disabled entirely** (last resort):
   ```bash
   systemctl disable apparmor
   # Edit /etc/default/grub: remove apparmor=1 security=apparmor
   update-grub && reboot
   ```

## What's NOT covered

- **php-fpm AppArmor profile** — operator FP intolerance (M9
  Snuffleupagus rejection, M33 Tetragon-default rejection); user
  PHP workload spans too many legitimate execs/file paths to
  enforce-confine without breaking sites.
- **mariadb / redis / pdns AppArmor profiles** — vendor packages
  ship their own; install.sh leaves them alone. Operator can
  enforce them manually.
- **Per-user profile per hosting account** — out of scope.
  Workload-level confinement is M34 nft (network) + M39 auditd
  (exec audit) + bubblewrap (M13, SSH shells).
- **Audit log shipping.** AppArmor AVC events go to `dmesg` /
  `journalctl -k`; dotnext doesn't archive them. Off-host shipping
  via syslog is operator-configured.

## Cross-references

- ADR-0086 — design decision + alternatives.
- ADR-0085 — narrow auditd (related forensic tier).
- ADR-0084 — per-user nft egress (related network tier).
- M13 (SSH shell sandbox) — bubblewrap for interactive shells.
- M33 / M39 (malware + exec audit) — file + process tiers.
