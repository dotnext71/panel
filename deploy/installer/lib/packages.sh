#!/usr/bin/env bash
# deploy/installer/lib/packages.sh — OS detection & dependency preflight.

dn_detect_os() {
  OS_ID="debian"
  OS_MAJOR=""
  OS_FAMILY="debian"
  if [[ -r /etc/os-release ]]; then
    local _id="" _ver="" _like=""
    # shellcheck disable=SC1091
    eval "$(. /etc/os-release 2>/dev/null; printf '_id=%q _ver=%q _like=%q' "${ID:-debian}" "${VERSION_ID:-}" "${ID_LIKE:-}")"
    OS_ID="${_id,,}"
    OS_MAJOR="${_ver%%.*}"
    case "$OS_ID" in
      almalinux|rocky|rhel|centos|ol|fedora) OS_FAMILY="rhel" ;;
      debian|ubuntu) OS_FAMILY="debian" ;;
      *)
        if [[ "${_like,,}" == *rhel* || "${_like,,}" == *fedora* || "${_like,,}" == *centos* ]]; then
          OS_FAMILY="rhel"
        fi
        ;;
    esac
  fi
  export OS_ID OS_MAJOR OS_FAMILY
}

dn_preflight_checks() {
  dn_detect_os

  local min_disk_kb=$((4 * 1024 * 1024))
  local avail_kb
  avail_kb="$(df -Pk / | awk 'NR==2 {print $4}')"
  if [[ -n "$avail_kb" && "$avail_kb" -lt "$min_disk_kb" ]]; then
    dn_warn "low root filesystem disk space: $((avail_kb / 1024)) MB available (4096 MB recommended)"
  fi

  local mem_kb
  mem_kb="$(awk '/MemTotal/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)"
  if [[ "$mem_kb" -gt 0 && "$mem_kb" -lt 900000 ]]; then
    dn_warn "system RAM is $((mem_kb / 1024)) MB; >= 1024 MB is recommended for full stack"
  fi

  for bin in systemctl curl openssl; do
    if ! command -v "$bin" >/dev/null 2>&1; then
      dn_warn "missing baseline utility: $bin"
    fi
  done
}
