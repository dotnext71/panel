#!/usr/bin/env bash
# deploy/installer/lib/common.sh — shared logging, OS detection, and lock helpers
# for dotNext installer and lifecycle scripts.

DOTNEXT_LOG_FILE="${DOTNEXT_LOG_FILE:-/var/log/dotnext-install.log}"

dn_log_init() {
  local log_dir
  log_dir="$(dirname "$DOTNEXT_LOG_FILE")"
  mkdir -p "$log_dir"
  touch "$DOTNEXT_LOG_FILE"
  chmod 0640 "$DOTNEXT_LOG_FILE" || true
}

dn_info() {
  printf "\033[1;36m[dotnext]\033[0m %s\n" "$*"
  [[ -w "$DOTNEXT_LOG_FILE" ]] && printf "[INFO] %s %s\n" "$(date -u +%FT%TZ)" "$*" >> "$DOTNEXT_LOG_FILE" || true
}

dn_warn() {
  printf "\033[1;33m[dotnext]\033[0m %s\n" "$*" >&2
  [[ -w "$DOTNEXT_LOG_FILE" ]] && printf "[WARN] %s %s\n" "$(date -u +%FT%TZ)" "$*" >> "$DOTNEXT_LOG_FILE" || true
}

dn_err() {
  printf "\033[1;31m[dotnext]\033[0m ERROR: %s\n" "$*" >&2
  [[ -w "$DOTNEXT_LOG_FILE" ]] && printf "[ERROR] %s %s\n" "$(date -u +%FT%TZ)" "$*" >> "$DOTNEXT_LOG_FILE" || true
}

dn_require_root() {
  if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
    dn_err "this installer must run as root (try: sudo $0 $*)"
    exit 1
  fi
}

dn_detect_os() {
  if [[ ! -r /etc/os-release ]]; then
    dn_err "cannot detect OS (/etc/os-release missing)"
    return 1
  fi
  # shellcheck disable=SC1091
  . /etc/os-release
  export DOTNEXT_OS_ID="${ID:-unknown}"
  export DOTNEXT_OS_VERSION="${VERSION_ID:-unknown}"
  case "$DOTNEXT_OS_ID:$DOTNEXT_OS_VERSION" in
    debian:12|debian:13|ubuntu:22.04|ubuntu:24.04|ubuntu:26.04)
      return 0
      ;;
    *)
      dn_warn "untested OS ${DOTNEXT_OS_ID} ${DOTNEXT_OS_VERSION}; Debian 12/13, Ubuntu 22.04/24.04/26.04 LTS, AlmaLinux/Rocky 8/9/10 recommended"
      return 0
      ;;
  esac
}
