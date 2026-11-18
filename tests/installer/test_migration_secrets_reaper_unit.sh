#!/usr/bin/env bash
# tests/installer/test_migration_secrets_reaper_unit.sh — the migration-secrets
# reaper can actually do its job, on fresh installs AND on updated boxes.
#
# `dotnext migrate reap-secrets` deletes per-job secrets under
# /etc/dotnext/migration-secrets AND stale staging trees under
# /var/lib/dotnext-migrations. The unit runs with ProtectSystem=strict, so
# every path it deletes from must be in ReadWritePaths — otherwise each
# staging removal fails with "read-only file system" (seen in the test box's
# daily run). And the unit is only fixed on existing boxes if `dotnext update`
# (provision_new_software) re-installs it; before this test it was installed
# by main() alone.
#
# Asserts, source-level:
#   1. ReadWritePaths covers both /etc/dotnext/migration-secrets and
#      /var/lib/dotnext-migrations.
#   2. provision_new_software calls install_migration_secrets_reaper.
#
# Run from repo root:
#     bash tests/installer/test_migration_secrets_reaper_unit.sh
set -euo pipefail

cd "$(dirname "$0")/../.."

unit="infrastructure/systemd/dotnext-migration-secrets-reap.service"
fail=0

rw=$(grep -E '^ReadWritePaths=' "$unit" | cut -d= -f2- | tr ' ' '\n' || true)
for need in /etc/dotnext/migration-secrets /var/lib/dotnext-migrations; do
  if ! grep -qxF "$need" <<<"$rw"; then
    echo "FAIL: $unit ReadWritePaths lacks $need — the reaper cannot delete there under ProtectSystem=strict"
    fail=1
  fi
done

# Body of provision_new_software (up to its closing brace at column 0).
body=$(awk '/^provision_new_software\(\) *\{/ {inside = 1; next} inside && /^\}/ {exit} inside' install.sh)
if [[ -z "$body" ]]; then
  echo "FAIL: provision_new_software not found in install.sh"
  fail=1
elif ! grep -qE '^[[:space:]]*install_migration_secrets_reaper([[:space:]]|$)' <<<"$body"; then
  echo "FAIL: provision_new_software does not call install_migration_secrets_reaper — updated boxes keep the old reaper unit"
  fail=1
fi

if [[ "$fail" -eq 0 ]]; then
  echo "OK: migration-secrets reaper unit can reach both of its paths and is re-installed on update"
else
  exit 1
fi
