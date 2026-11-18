#!/usr/bin/env bash
# deploy/installer/lib/services.sh — systemd unit lifecycle, snapshot, and rollback helpers.

DOTNEXT_ROLLBACK_DIR="${DOTNEXT_ROLLBACK_DIR:-/var/lib/dotnext-panel/rollback}"

dn_units() {
  echo "dotnext-agent.service dotnext-kratos.service dotnext-panel.service"
}

dn_snapshot_state() {
  local stamp ts_dir
  stamp="$(date -u +%Y%m%dT%H%M%SZ)"
  ts_dir="${DOTNEXT_ROLLBACK_DIR}/${stamp}"
  mkdir -p "$ts_dir"
  for bin in dotnext dotnext-agent dotnext-installer dotnext-webdav dotnext-ssh-shell; do
    if [[ -x "/usr/local/bin/${bin}" ]]; then
      cp -a "/usr/local/bin/${bin}" "${ts_dir}/${bin}" || true
    fi
  done
  if [[ -f /etc/dotnext/config.toml ]]; then
    cp -a /etc/dotnext/config.toml "${ts_dir}/config.toml" || true
  fi
  ln -sfn "$ts_dir" "${DOTNEXT_ROLLBACK_DIR}/latest"
  dn_info "saved rollback snapshot at ${ts_dir}"
}

dn_rollback_latest() {
  local latest="${DOTNEXT_ROLLBACK_DIR}/latest"
  if [[ ! -d "$latest" ]]; then
    dn_err "no rollback snapshot available in ${DOTNEXT_ROLLBACK_DIR}"
    return 1
  fi
  dn_warn "restoring binaries & config from ${latest}"
  for bin in dotnext dotnext-agent dotnext-installer dotnext-webdav dotnext-ssh-shell; do
    if [[ -f "${latest}/${bin}" ]]; then
      install -m 0755 "${latest}/${bin}" "/usr/local/bin/${bin}"
    fi
  done
  if [[ -f "${latest}/config.toml" ]]; then
    cp -a "${latest}/config.toml" /etc/dotnext/config.toml
  fi
  systemctl daemon-reload || true
  for u in $(dn_units); do
    systemctl restart "$u" || true
  done
}

dn_verify_health() {
  local url="${1:-http://127.0.0.1:8443/health}"
  local tries=15
  for _ in $(seq 1 "$tries"); do
    if curl -fsSk "$url" >/dev/null 2>&1; then
      dn_info "health check passed: ${url}"
      return 0
    fi
    sleep 1
  done
  dn_err "health check failed after ${tries}s: ${url}"
  return 1
}
