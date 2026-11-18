//go:build linux

package cpanel

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

// readRegularBeneath reads base/rel (base opened O_NOFOLLOW, rel resolved
// strictly beneath it with no symlink in any component), when it is a regular
// file of at most max bytes.
func readRegularBeneath(base, rel string, max int64) ([]byte, error) {
	bfd, err := unix.Open(base, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(bfd)
	fd, err := unix.Openat2(bfd, rel, &unix.OpenHow{
		Flags:   unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_BENEATH | unix.RESOLVE_NO_SYMLINKS,
	})
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), rel)
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file")
	}
	if fi.Size() > max {
		return nil, fmt.Errorf("larger than %d bytes", max)
	}
	return io.ReadAll(io.LimitReader(f, max+1))
}
