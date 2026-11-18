//go:build !linux && !darwin

package peercred

import (
	"errors"
	"net"
)

func readUnixCred(*net.UnixConn) (Credentials, error) {
	return Credentials{}, errors.New("peercred: unsupported operating system")
}
