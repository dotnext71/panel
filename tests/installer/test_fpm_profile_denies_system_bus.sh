#!/usr/bin/env bash
# tests/installer/test_fpm_profile_denies_system_bus.sh — the dotnext-fpm-app
# AppArmor profile denies tenant PHP the system D-Bus socket.
#
# abstractions/nameservice includes dbus-strict, which allows
# /run/dbus/system_bus_socket, and tenant PHP has no use for the system bus.
# The profile ships, and often stays, in complain mode, where only an explicit
# `deny` rule is enforced; an allow-list omission would only be logged.
#
# Asserts:
#   1. The dotnext-fpm-app profile block carries
#      `deny /run/dbus/system_bus_socket rw,` (and the /var/run alias).
#   2. When apparmor_parser and the system abstractions are present: the
#      shipped profile compiles.
#
# Run from repo root:
#     bash tests/installer/test_fpm_profile_denies_system_bus.sh
#
# Exit 0 = pass.
set -euo pipefail

cd "$(dirname "$0")/../.."

profile=infrastructure/apparmor/usr.local.libexec.dotnext.fpm-exec
fail=0

# --- 1. The deny rules sit inside the dotnext-fpm-app profile block. ---
block=$(awk '/^profile dotnext-fpm-app /{on=1} on{print} on && /^}/{exit}' "$profile")
if [[ -z "$block" ]]; then
  echo "FAIL: profile dotnext-fpm-app not found in $profile"
  exit 1
fi
for path in /run/dbus/system_bus_socket /var/run/dbus/system_bus_socket; do
  if ! grep -qE "^[[:space:]]*deny[[:space:]]+${path//\//\\/}[[:space:]]+rw,[[:space:]]*$" <<<"$block"; then
    echo "FAIL: dotnext-fpm-app does not deny $path"
    fail=1
  fi
done

# --- 2. The profile compiles (where AppArmor's tools are installed). ---
if command -v apparmor_parser >/dev/null 2>&1 && [[ -f /etc/apparmor.d/tunables/global ]]; then
  if ! out=$(apparmor_parser --skip-kernel-load --skip-cache -I /etc/apparmor.d "$profile" 2>&1); then
    echo "FAIL: $profile does not compile:"
    echo "$out"
    fail=1
  fi
else
  echo "SKIP: apparmor_parser or /etc/apparmor.d/tunables/global missing; compile check skipped"
fi

if [[ $fail -ne 0 ]]; then
  exit 1
fi
echo "PASS: dotnext-fpm-app denies the system D-Bus socket"
