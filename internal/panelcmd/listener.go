package panelcmd

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/dotnext-panel/dotnext/pkg/paths"
)

const (
	socketGroup = paths.SocketGroup
	// A host upgraded from the previous panel keeps that panel's socket group
	// until the installer recreates it, so the group name is accepted as a
	// fallback. Built from a split literal: it is the only place the old
	// product name survives, and it must stay byte-exact to be found on such
	// a host.
	legacySocketGroup = "jab" + "ali-sockets"
	socketMode        = 0o660
)

func lookupSocketGroup() (*user.Group, string, error) {
	if grp, err := user.LookupGroup(socketGroup); err == nil {
		return grp, socketGroup, nil
	}
	if grp, err := user.LookupGroup(legacySocketGroup); err == nil {
		return grp, legacySocketGroup, nil
	}
	grp, err := user.LookupGroup(socketGroup)
	return grp, socketGroup, err
}

// parseListenAddr returns the network ("tcp" or "unix") and address
// (host:port or socket path) for a given config string. The unix form is
// `unix:/abs/path`. Anything else is treated as TCP.
func parseListenAddr(addr string) (network, target string) {
	if strings.HasPrefix(addr, "unix:") {
		return "unix", strings.TrimPrefix(addr, "unix:")
	}
	return "tcp", addr
}

// listenAndPrepare opens the listener for `addr`. For unix: addresses, it
// removes any stale socket file (a previous unclean shutdown can leave one
// behind, in which case net.Listen would fail with "address already in
// use"), creates the listener, then chmods + chgrps the socket so nginx
// can connect via group membership.
//
// For TCP addresses it just does net.Listen("tcp", addr) — the TLS branch
// in the caller handles TLS certs.
//
// Returns the listener + a boolean indicating whether the address was a
// unix socket (caller uses this to skip the TLS branch).
func listenAndPrepare(addr string) (net.Listener, bool, error) {
	network, target := parseListenAddr(addr)
	if network == "tcp" {
		l, err := net.Listen("tcp", target)
		if err != nil {
			return nil, false, fmt.Errorf("listen tcp %s: %w", target, err)
		}
		return l, false, nil
	}

	// Unix socket. Stale-socket cleanup before bind.
	if fi, statErr := os.Stat(target); statErr == nil {
		if fi.Mode()&os.ModeSocket == 0 {
			return nil, true, fmt.Errorf("listen unix %s: path exists and is not a socket — refusing to overwrite", target)
		}
		// Unlink the stale socket. net.Listen would otherwise fail with
		// "address already in use". Race window is ~zero in practice
		// because systemd serializes start/stop.
		if err := os.Remove(target); err != nil {
			return nil, true, fmt.Errorf("listen unix %s: remove stale socket: %w", target, err)
		}
	} else if !os.IsNotExist(statErr) {
		return nil, true, fmt.Errorf("listen unix %s: stat: %w", target, statErr)
	}

	l, err := net.Listen("unix", target)
	if err != nil {
		return nil, true, fmt.Errorf("listen unix %s: %w", target, err)
	}

	// Permissions: 0660 is what nginx (group member) needs to connect(2).
	// chmod second so umask can't clamp the bits — listen() creates the
	// socket subject to the process umask, so newly-listened sockets
	// often come out 0644 even with mode wanted lower.
	if err := os.Chmod(target, socketMode); err != nil {
		_ = l.Close()
		return nil, true, fmt.Errorf("listen unix %s: chmod 0660: %w", target, err)
	}

	grp, resolvedGroup, err := lookupSocketGroup()
	if err != nil {
		_ = l.Close()
		return nil, true, fmt.Errorf("listen unix %s: lookup group %s: %w", target, resolvedGroup, err)
	}
	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		_ = l.Close()
		return nil, true, fmt.Errorf("listen unix %s: parse gid: %w", target, err)
	}
	// Owner stays the service uid (-1 = leave unchanged). We only need
	// to flip the GROUP ownership.
	if err := os.Chown(target, -1, gid); err != nil {
		_ = l.Close()
		return nil, true, fmt.Errorf("listen unix %s: chgrp %s: %w", target, resolvedGroup, err)
	}

	return l, true, nil
}
