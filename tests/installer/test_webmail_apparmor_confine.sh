#!/usr/bin/env bash
# tests/installer/test_webmail_apparmor_confine.sh — regression coverage for
# JAB-375. The dotnext-webmail (Bulwark, Next.js) service is confined by the
# AppArmor profile dotnext-bulwark. Two hard constraints, both learned the hard
# way on the .86 box:
#
#   A. The unit attaches the profile via AppArmorProfile= and MUST NOT set any
#      systemd private-mount-namespace directive (ProtectSystem= / PrivateTmp= /
#      ProtectHome=). A mount namespace + this profile's label transition makes
#      Node self-abort at node::Start (status=6/ABRT, no message) — deterministic,
#      any one of the three is enough. So confinement is via AppArmor file MAC,
#      not the mount-ns directives, and TMPDIR+`deny /tmp/**` replace PrivateTmp.
#
#   B. The profile must grant the REAL runtime paths. The app runs as
#      /usr/bin/node /opt/dotnext-webmail/server-unix.js, config in
#      /etc/dotnext/bulwark.env, state in /var/lib/dotnext-webmail. The M40
#      profile granted the old /usr/local/bin + /var/lib/dotnext-bulwark layout
#      and confined nothing real.
#
# Source-level asserts (no host mutation). Run from repo root:
#     bash tests/installer/test_webmail_apparmor_confine.sh
set -euo pipefail

cd "$(dirname "$0")/../.."

fail=0
unit=infrastructure/systemd/dotnext-webmail.service
prof=infrastructure/apparmor/usr.local.bin.dotnext-bulwark

# --- A. unit: AppArmor attach, NO mount-ns directives, dedicated tmp ---
if ! grep -qE '^AppArmorProfile=-?dotnext-bulwark' "$unit"; then
  echo "FAIL: $unit must attach the profile (AppArmorProfile=-dotnext-bulwark) (JAB-375)"
  fail=1
fi
for d in ProtectSystem PrivateTmp ProtectHome; do
  if grep -qE "^${d}=" "$unit"; then
    echo "FAIL: $unit sets ${d}= — a private mount ns + the AppArmor label aborts Node at node::Start (JAB-375). Remove it; the profile provides the file MAC."
    fail=1
  fi
done
# ReadWritePaths is only meaningful with ProtectSystem; it must go too.
if grep -qE '^ReadWritePaths=' "$unit"; then
  echo "FAIL: $unit sets ReadWritePaths= — meaningless without ProtectSystem and part of the mount-ns coupling (JAB-375)"
  fail=1
fi
# dedicated tmp subtree replaces PrivateTmp isolation
if ! grep -qE '^Environment=TMPDIR=/var/lib/dotnext-webmail/tmp' "$unit"; then
  echo "FAIL: $unit must set TMPDIR to the dedicated per-service subtree (replaces PrivateTmp) (JAB-375)"
  fail=1
fi
if ! grep -qE '^ExecStartPre=.*mkdir.*-p /var/lib/dotnext-webmail/tmp' "$unit"; then
  echo "FAIL: $unit must pre-create the TMPDIR subtree (JAB-375)"
  fail=1
fi
# NoNewPrivileges is not the trigger and is kept.
if ! grep -qE '^NoNewPrivileges=true' "$unit"; then
  echo "FAIL: $unit should keep NoNewPrivileges=true (it is not the abort trigger; good hardening)"
  fail=1
fi

# --- B. profile: real paths, dedicated-tmp isolation, abi pin ---
if ! grep -qE '^\s*/opt/dotnext-webmail/\*\*\s+r,' "$prof"; then
  echo "FAIL: $prof must grant the real app tree /opt/dotnext-webmail/** r (JAB-375)"
  fail=1
fi
if ! grep -qE '^\s*/etc/dotnext/bulwark\.env\s+r,' "$prof"; then
  echo "FAIL: $prof must grant the real config /etc/dotnext/bulwark.env r (JAB-375)"
  fail=1
fi
if ! grep -qE '^\s*/var/lib/dotnext-webmail/\*\*\s+rwk,' "$prof"; then
  echo "FAIL: $prof must grant the real state dir /var/lib/dotnext-webmail/** rwk (JAB-375)"
  fail=1
fi
# shared /tmp denied — the dedicated-tmp isolation half
if ! grep -qE '^\s*deny /tmp/\*\* ' "$prof"; then
  echo "FAIL: $prof must deny the shared /tmp (tmp is redirected to the dedicated subtree) (JAB-375)"
  fail=1
fi
# panel secrets still denied
if ! grep -qE '^\s*deny /etc/dotnext/\*\* ' "$prof"; then
  echo "FAIL: $prof must still deny panel secrets under /etc/dotnext/** (JAB-351/357)"
  fail=1
fi
# abi 3.0 pin retained (unix bind on broken-mediation kernels)
if ! grep -qE '^abi <abi/3\.0>,' "$prof"; then
  echo "FAIL: $prof must keep the abi <abi/3.0> pin so unix bind works on Debian13/Ubuntu24.04 (GH #705)"
  fail=1
fi
# the profile must NOT grant broad write to the app tree beyond .next/cache
# (version-check writes are silently denied to keep the tree read-only)
if ! grep -qE '^\s*deny /opt/dotnext-webmail/data/\*\* ' "$prof"; then
  echo "FAIL: $prof should silently deny writes under the read-only app tree's data dir (soak-clean) (JAB-375)"
  fail=1
fi

if [[ "$fail" -eq 0 ]]; then
  echo "OK: webmail AppArmor confinement guards hold (JAB-375)"
else
  exit 1
fi
