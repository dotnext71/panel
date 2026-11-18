//go:build !linux

package fsperm

import (
	"errors"
	"os"
)

var errLinuxOnly = errors.New("fsperm: fd-relative permission operations require Linux")

// ReownStats reports one ReownTree walk.
type ReownStats struct {
	Reowned    map[uint32]int
	Scanned    int
	Stopped    bool
	Failed     int
	FirstError string
	Linked     int
}

// ReownRoot describes a target root for ReownTrees.
type ReownRoot struct {
	Anchor string
	Root   string
	Match  func(owner uint32) bool
}

func GroupSetgidTree(string, int) error                  { return errLinuxOnly }
func RepairDocrootGroup(string, string, int) error       { return errLinuxOnly }
func ChownTree(string, string, int, int) error           { return errLinuxOnly }
func OpenDirUnder(string, string, uint32, int, int) (*os.File, error) {
	return nil, errLinuxOnly
}
func MoveDirUnder(string, string, string, string) error  { return errLinuxOnly }
func RemoveEmptyDirUnder(string, string) error           { return errLinuxOnly }
func OpenExistingDirUnder(string, string) (*os.File, error) {
	return nil, errLinuxOnly
}
func ReownTree(string, string, func(uint32) bool, int, func(uint32) int, func() bool) (ReownStats, error) {
	return ReownStats{}, errLinuxOnly
}
func ReownTrees([]ReownRoot, int, func(uint32) int, func() bool) (ReownStats, error) {
	return ReownStats{}, errLinuxOnly
}
