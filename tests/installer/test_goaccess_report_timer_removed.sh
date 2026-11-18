#!/usr/bin/env bash
# tests/installer/test_goaccess_report_timer_removed.sh — the dotnext-goaccess
# report timer is gone.
#
# It rendered an HTML report of every /var/log/nginx/*.access.log into
# /var/lib/dotnext-goaccess every 15 minutes, and nothing read the reports.
# install_goaccess must no longer create the timer, its service or its
# generator script, and must take them down on hosts that still have them.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
install_sh="$here/../../install.sh"

body="$(awk '/^install_goaccess\(\) \{/{f=1} f{print} f && /^}/{exit}' "$install_sh")"
[[ -n "$body" ]] || { echo "FAIL: install_goaccess not found"; exit 1; }

fail=0
for created in 'OnCalendar=' 'cat > "$timer_unit"' 'cat > "$service_unit"' 'cat > "$generator_script"' \
               'systemctl enable --quiet dotnext-goaccess.timer'; do
  if grep -qF -- "$created" <<<"$body"; then
    echo "FAIL: install_goaccess still sets up the report timer ($created)"
    fail=1
  fi
done
for teardown in 'systemctl disable --now dotnext-goaccess.timer' \
                '/etc/systemd/system/dotnext-goaccess.timer' \
                '/usr/local/bin/dotnext-goaccess-generator' \
                'rm -rf /var/lib/dotnext-goaccess' \
                'systemctl daemon-reload'; do
  if ! grep -qF -- "$teardown" <<<"$body"; then
    echo "FAIL: install_goaccess does not remove the old timer ($teardown missing)"
    fail=1
  fi
done
# The shared goaccess config stays: the bandwidth scan and the Logs tab's
# live report read it.
grep -qF '/etc/goaccess/goaccess.conf' <<<"$body" || { echo "FAIL: goaccess config no longer written"; fail=1; }

[[ $fail -eq 0 ]] && echo "PASS: dotnext-goaccess report timer removed, config kept"
exit $fail
