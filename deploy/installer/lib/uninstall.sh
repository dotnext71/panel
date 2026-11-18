#!/usr/bin/env bash
# deploy/installer/lib/uninstall.sh — safe dotNext uninstall mechanism.
# Preserves /home tenant files & databases unless --purge-data is passed.
set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=deploy/installer/lib/common.sh
source "${SCRIPT_DIR}/common.sh"

PURGE_DATA=0
ASSUME_YES=0

for arg in "$@"; do
  case "$arg" in
    --purge-data) PURGE_DATA=1 ;;
    -y|--yes)     ASSUME_YES=1 ;;
  esac
done

dn_require_root
dn_log_init

if [[ "$ASSUME_YES" -ne 1 && -t 0 ]]; then
  read -r -p "Stop dotNext services and remove panel binaries? [y/N] " ans
  if [[ ! "$ans" =~ ^[Yy]$ ]]; then
    dn_info "uninstall aborted."
    exit 0
  fi
fi

for unit in dotnext-panel dotnext-agent dotnext-kratos dotnext-mailrelay dotnext-api dotnext-agent dotnext-kratos; do
  if systemctl list-unit-files "${unit}.service" >/dev/null 2>&1; then
    systemctl stop "${unit}.service" 2>/dev/null || true
    systemctl disable "${unit}.service" 2>/dev/null || true
    rm -f "/etc/systemd/system/${unit}.service"
  fi
done
systemctl daemon-reload || true

rm -f /usr/local/bin/dotnext /usr/local/bin/dotnext-agent /usr/local/bin/dotnext-installer \
      /usr/local/bin/dotnext-webdav /usr/local/bin/dotnext-ssh-shell \
      /usr/local/bin/dotnext-mailhook /usr/local/bin/dotnext-sendmail \
      /usr/local/bin/dotnext-cli /usr/local/bin/dotnext-agent /usr/local/bin/dotnext-api

if [[ "$PURGE_DATA" -eq 1 ]]; then
  dn_warn "purging dotNext configuration and panel state (--purge-data)"
  rm -rf /etc/dotnext /var/lib/dotnext /var/lib/dotnext-panel /opt/dotnext
else
  dn_info "preserved /etc/dotnext, /var/lib/dotnext*, and /home tenant data (pass --purge-data to remove panel state)"
fi

dn_info "dotNext uninstall complete."
