//go:build darwin

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
	var cred *unix.Xucred
	var credErr error
	if err := raw.Control(func(fd uintptr) {
		cred, credErr = unix.GetsockoptXucred(int(fd), unix.SOL_LOCAL, unix.LOCAL_PEERCRED)
	}); err != nil {
		return Credentials{}, fmt.Errorf("peercred: control: %w", err)
	}
	if credErr != nil {
		return Credentials{}, fmt.Errorf("peercred: getsockopt: %w", credErr)
	}
	var gid uint32
	if cred.Ngroups > 0 {
		gid = cred.Groups[0]
	}
	return Credentials{UID: cred.Uid, GID: gid}, nil
}
