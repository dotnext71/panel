//go:build !linux

package cpanel

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// readRegularBeneath provides a component-wise Lstat + containment check on
// non-Linux developer hosts (e.g. macOS unit tests). Production Linux builds
// compile mail_shadow_linux.go which enforces kernel openat2(RESOLVE_BENEATH).
func readRegularBeneath(base, rel string, max int64) ([]byte, error) {
	if filepath.IsAbs(rel) || strings.ContainsRune(rel, 0) {
		return nil, fmt.Errorf("invalid relative path %q", rel)
	}
	baseInfo, err := os.Lstat(base)
	if err != nil {
		return nil, err
	}
	if baseInfo.Mode()&os.ModeSymlink != 0 || !baseInfo.IsDir() {
		return nil, fmt.Errorf("base is not a real directory")
	}
	cur := base
	parts := strings.Split(filepath.Clean(rel), string(filepath.Separator))
	for i, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			return nil, fmt.Errorf("path traversal rejected")
		}
		cur = filepath.Join(cur, part)
		fi, lerr := os.Lstat(cur)
		if lerr != nil {
			return nil, lerr
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("symlink rejected: %s", part)
		}
		if i < len(parts)-1 && !fi.IsDir() {
			return nil, fmt.Errorf("not a directory: %s", part)
		}
	}
	f, err := os.Open(cur)
	if err != nil {
		return nil, err
	}
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
