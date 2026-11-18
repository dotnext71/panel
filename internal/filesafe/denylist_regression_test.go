// SPDX-License-Identifier: MIT

package filesafe

import "testing"

// The admin File Manager is root-scoped, so its deny-list is the last thing
// standing between an admin UI request and the panel's own secrets, sockets and
// binaries. A brand rename once collapsed the list's "legacy path -> new path"
// pairs into duplicates, which silently DROPPED the legacy entries: /etc/<old>,
// /run/<old> and /var/lib/<old> were no longer denied. These tests pin the
// behaviour rather than the spelling.
func TestAdminDeniedPrefixes_LegacyControlPlaneStillDenied(t *testing.T) {
	denied := AdminDeniedPrefixes()

	// Build the legacy paths the same way the package does, so this test cannot
	// reintroduce the old product name as a literal either.
	legacy := legacyName
	want := []string{
		"/etc/" + legacy,
		"/run/" + legacy,
		"/var/lib/" + legacy,
		"/usr/local/bin/" + legacy,
	}
	for _, w := range want {
		found := false
		for _, d := range denied {
			if d == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("legacy control-plane path %q is not denied", w)
		}
	}
}

// A duplicate entry means a rename collapsed a pair; it is always a mistake
// because the entry it replaced is gone.
func TestAdminDeniedPrefixes_NoDuplicateEntries(t *testing.T) {
	seen := map[string]bool{}
	for _, d := range AdminDeniedPrefixes() {
		if seen[d] {
			t.Errorf("duplicate deny prefix %q — a collapsed legacy/new pair?", d)
		}
		seen[d] = true
	}
}

// The dotNext control plane must stay denied, and the mutable-root allow-list
// must not overlap it.
func TestAdminDeniedPrefixes_ControlPlaneDenied(t *testing.T) {
	denied := AdminDeniedPrefixes()
	for _, must := range []string{
		"/etc/dotnext", "/etc/shadow", "/etc/sudoers", "/etc/ssh",
		"/run/dotnext", "/run/dotnext-panel", "/var/lib/dotnext",
		"/usr/local/bin/dotnext-api", "/usr/local/bin/dotnext-agent",
	} {
		found := false
		for _, d := range denied {
			if d == must {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("control-plane path %q is not denied", must)
		}
	}
	for _, m := range DefaultAdminMutableRoots {
		for _, d := range denied {
			if m == d {
				t.Errorf("mutable root %q is also denied", m)
			}
		}
	}
}
