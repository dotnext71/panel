//go:build !linux

package filesafe

import (
	"context"
	"errors"
	"os"
	"time"
)

var errLinuxOnly = errors.New("filesafe: openat2 confinement requires Linux")

// LeafInfo mirrors the Linux openat2 stat metadata for cross-platform compilation.
type LeafInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	IsLink  bool
	UID     uint32
	GID     uint32
}

// ExtractDir is a non-Linux stub for archive extraction.
type ExtractDir struct{}

func (s *Scope) OpenInScope(string, int, os.FileMode) (*os.File, error) {
	return nil, errLinuxOnly
}

func (s *Scope) OpenDirInScope(string) (*os.File, error) {
	return nil, errLinuxOnly
}

func (s *Scope) OpenParentInScope(string) (*os.File, string, error) {
	return nil, "", errLinuxOnly
}

func (s *Scope) StatInScope(string) (*LeafInfo, error) {
	return nil, errLinuxOnly
}

func (s *Scope) ExistsInScope(string) (*LeafInfo, error) {
	return nil, errLinuxOnly
}

func (s *Scope) ReadDirInScope(string) ([]LeafInfo, error) {
	return nil, errLinuxOnly
}

func (s *Scope) CreateExclInScope(string, os.FileMode) (*os.File, error) {
	return nil, errLinuxOnly
}

func (s *Scope) RenameInScope(string, string) error {
	return errLinuxOnly
}

func (s *Scope) RenameIntoInScope(string, string) error {
	return errLinuxOnly
}

func (s *Scope) RemoveInScope(string, bool) error {
	return errLinuxOnly
}

func (s *Scope) RemoveAllInScope(context.Context, string) error {
	return errLinuxOnly
}

func (s *Scope) CopyTreeInScope(context.Context, string, string, int, int) (int64, error) {
	return 0, errLinuxOnly
}

func (s *Scope) CopyTreeInScopeProgress(context.Context, string, string, int, int, func(int64)) (int64, error) {
	return 0, errLinuxOnly
}

func (s *Scope) CountTreeBytesInScope(context.Context, string) (int64, error) {
	return 0, errLinuxOnly
}

func (s *Scope) ReadlinkInScope(string) (string, error) {
	return "", errLinuxOnly
}

func (s *Scope) ChmodInScope(string, os.FileMode) error {
	return errLinuxOnly
}

func (s *Scope) MkdirInScope(string, bool, os.FileMode) (bool, error) {
	return false, errLinuxOnly
}

func (s *Scope) NewExtractDir(string, int, int) (*ExtractDir, error) {
	return nil, errLinuxOnly
}

func (e *ExtractDir) Close() error                       { return nil }
func (e *ExtractDir) MkdirAll(string, os.FileMode) error { return errLinuxOnly }
func (e *ExtractDir) Create(string, os.FileMode) (*os.File, error) {
	return nil, errLinuxOnly
}
func (e *ExtractDir) Remove(string) error { return errLinuxOnly }

func (s *Scope) MkdirChownInScope(string, os.FileMode, bool, int, int) error {
	return errLinuxOnly
}
