#!/usr/bin/env bash
# tests/installer/test_snuffleupagus_rules_traversal_heal.sh — tenant PHP must
# keep loading the Snuffleupagus rules.
#
# Tenant PHP-FPM masters run as the tenant and load
# /etc/dotnext/snuffleupagus/active.rules. A test box had /etc/dotnext at 0750
# (owned by the panel user): every tenant master failed to load the rules on
# its next reload, and every tenant site answered 500/502 until /etc/dotnext
# was put back to root:root 0755 and the masters restarted. Nothing on the
# update path checked that directory (ensure_snuffleupagus_loadable covers the
# library dir and the ini, tools/test-snuffleupagus-loadable.sh).
#
# 1. Functional: ensure_snuffleupagus_rules_readable against a temporary tree
#    (paths rewritten), with chown, owner lookups and systemctl stubbed.
# 2. Static: `dotnext update` (provision_new_software) calls it, and install.sh
#    never makes /etc/dotnext anything but 0755.
#
# Run from repo root:
#     bash tests/installer/test_snuffleupagus_rules_traversal_heal.sh
set -euo pipefail

cd "$(dirname "$0")/../.."

fail=0
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

body=$(awk '/^ensure_snuffleupagus_rules_readable\(\)/{f=1} f{print} f&&/^\}/{exit}' install.sh)
if [[ -z "$body" ]]; then
  echo "FAIL: ensure_snuffleupagus_rules_readable() not found in install.sh"
  exit 1
fi
body=${body//\/etc\/dotnext/$tmp\/etc\/dotnext}

# run_case <etc owner> <etc mode> <sp owner> <sp mode> <rules mode>; prints
# the recorded chown/systemctl calls and the resulting modes.
run_case() {
  (
    rm -rf "$tmp/etc"
    mkdir -p "$tmp/etc/dotnext/snuffleupagus"
    : >"$tmp/etc/dotnext/snuffleupagus/active.rules"
    command chmod "$5" "$tmp/etc/dotnext/snuffleupagus/active.rules"
    command chmod "$4" "$tmp/etc/dotnext/snuffleupagus"
    command chmod "$2" "$tmp/etc/dotnext"
    E_OWN=$1 S_OWN=$3
    _log() { :; }
    _ok() { :; }
    _warn() { :; }
    stat() {
      if [[ "$2" == "%U:%G" ]]; then
        case "$3" in
          "$tmp/etc/dotnext") echo "$E_OWN" ;;
          "$tmp/etc/dotnext/snuffleupagus") echo "$S_OWN" ;;
          *) echo root:root ;;
        esac
        return 0
      fi
      command stat "$@"
    }
    chown() { echo "chown $*" >&3; }
    systemctl() { echo "systemctl $*" >&3; }
    eval "$body"
    ensure_snuffleupagus_rules_readable >/dev/null 2>&1
    echo "modes $(command stat -c '%a' "$tmp/etc/dotnext" "$tmp/etc/dotnext/snuffleupagus" "$tmp/etc/dotnext/snuffleupagus/active.rules" | tr '\n' ' ')" >&3
  ) 3>&1
}

check() { # <name> <output> <want substring>...
  local name=$1 out=$2
  shift 2
  for w in "$@"; do
    if [[ "$out" != *"$w"* ]]; then
      echo "FAIL: $name: missing [$w] in:"
      echo "$out" | sed 's/^/    /'
      fail=1
    fi
  done
}

out=$(run_case dotnext:dotnext 750 root:root 755 644)
check "panel-owned 0750 /etc/dotnext" "$out" "chown root:root $tmp/etc/dotnext" "modes 755 755 644" "systemctl try-restart dotnext-fpm@*.service"
if [[ "$out" == *"chown root:root $tmp/etc/dotnext/snuffleupagus"* ]]; then
  echo "FAIL: changed the owner of a root-owned snuffleupagus dir"
  fail=1
fi

out=$(run_case root:root 700 root:root 700 600)
check "closed dirs and rules" "$out" "modes 755 755 644" "systemctl try-restart dotnext-fpm@*.service"

out=$(run_case root:root 755 root:root 755 644)
if [[ "$out" == *chown* || "$out" == *systemctl* ]]; then
  echo "FAIL: a healthy host was changed: $out"
  fail=1
fi
check "healthy host" "$out" "modes 755 755 644"

# A symlinked /etc/dotnext is not followed.
(
  rm -rf "$tmp/etc" "$tmp/elsewhere"
  mkdir -p "$tmp/etc" "$tmp/elsewhere/snuffleupagus"
  command chmod 700 "$tmp/elsewhere"
  ln -s "$tmp/elsewhere" "$tmp/etc/dotnext"
  _log() { :; }
  chown() { echo "chown $*" >>"$tmp/calls"; }
  systemctl() { :; }
  stat() {
    if [[ "$2" == "%U:%G" ]]; then
      echo nobody:nogroup
      return 0
    fi
    command stat "$@"
  }
  eval "$body"
  ensure_snuffleupagus_rules_readable >/dev/null 2>&1
)
if [[ -s "$tmp/calls" || "$(stat -c '%a' "$tmp/elsewhere")" != 700 ]]; then
  echo "FAIL: followed a symlinked /etc/dotnext"
  fail=1
fi
chmod 755 "$tmp/elsewhere"

# Static: the update path calls it; install.sh only ever makes /etc/dotnext 0755.
prov=$(awk '/^provision_new_software\(\)/{f=1} f{print} f&&/^\}/{exit}' install.sh)
if ! grep -qE '^\s*ensure_snuffleupagus_rules_readable\s*$' <<<"$prov"; then
  echo "FAIL: provision_new_software does not call ensure_snuffleupagus_rules_readable"
  fail=1
fi
if bad=$(grep -nE 'install -d\b.*/etc/dotnext/?\s*$' install.sh | grep -vE -- '-m 0?755\b'); then
  echo "FAIL: install.sh makes /etc/dotnext other than 0755:"
  echo "$bad" | sed 's/^/    /'
  fail=1
fi

if [[ "$fail" == 0 ]]; then
  echo "PASS: test_snuffleupagus_rules_traversal_heal"
fi
exit "$fail"
