//go:build linux

package peercred

import (
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

func readUnixCred(uc *net.UnixConn) (Credentials, error) {
	raw, err := uc.SyscallConn()
	if err != nil {
		return Credentials{}, fmt.Errorf("peercred: syscall conn: %w", err)
	}
	var cred *unix.Ucred
	var credErr error
	if err := raw.Control(func(fd uintptr) {
		cred, credErr = unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
	}); err != nil {
		return Credentials{}, fmt.Errorf("peercred: control: %w", err)
	}
	if credErr != nil {
		return Credentials{}, fmt.Errorf("peercred: getsockopt: %w", credErr)
	}
	return Credentials{UID: cred.Uid, GID: cred.Gid}, nil
}
