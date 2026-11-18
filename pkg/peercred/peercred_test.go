package peercred

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestFromConn_UnixSocket(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("/tmp", "dnpc-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	sock := filepath.Join(dir, "test.sock")
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	done := make(chan Credentials, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, aerr := ln.Accept()
		if aerr != nil {
			errCh <- aerr
			return
		}
		defer conn.Close()
		c, cerr := FromConn(conn)
		if cerr != nil {
			errCh <- cerr
			return
		}
		done <- c
	}()

	client, err := net.Dial("unix", sock)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer client.Close()

	select {
	case cerr := <-errCh:
		t.Fatalf("server err: %v", cerr)
	case creds := <-done:
		if int(creds.UID) != os.Getuid() {
			t.Fatalf("uid=%d, want %d", creds.UID, os.Getuid())
		}
	}
}
