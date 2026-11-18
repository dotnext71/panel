package server

import (
	"net"

	"github.com/dotnext-panel/dotnext/pkg/peercred"
)

// peerUID extracts the connecting process's UID from a Unix socket connection.
func peerUID(conn net.Conn) (uint32, error) {
	creds, err := peercred.FromConn(conn)
	if err != nil {
		return 0, err
	}
	return creds.UID, nil
}
