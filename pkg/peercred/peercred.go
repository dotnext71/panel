// Package peercred reads kernel-authenticated Unix domain socket peer
// credentials (UID and primary GID) without trusting any caller-supplied data.
package peercred

import (
	"errors"
	"net"
)

// Credentials identifies the peer process that opened a Unix socket connection.
type Credentials struct {
	UID uint32
	GID uint32
}

// ErrNotUnixConn is returned when FromConn is passed a non-Unix net.Conn.
var ErrNotUnixConn = errors.New("peercred: not a unix connection")

// FromConn reads the peer UID and primary GID from a *net.UnixConn.
func FromConn(conn net.Conn) (Credentials, error) {
	uc, ok := conn.(*net.UnixConn)
	if !ok {
		return Credentials{}, ErrNotUnixConn
	}
	return readUnixCred(uc)
}
