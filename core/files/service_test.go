package files

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestOpenTenantAndAdminScopes(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	scope, err := OpenTenantScope("u1", "alice", []string{root})
	if err != nil {
		t.Fatalf("OpenTenantScope: %v", err)
	}
	if _, err := scope.Resolve(filepath.Join(root, "index.php")); err != nil {
		t.Fatalf("Resolve inside docroot: %v", err)
	}
	if _, err := scope.Resolve("/etc/dotnext/config.toml"); err == nil {
		t.Fatal("expected escape outside docroot to fail")
	}

	if _, err := OpenAdminScope("a1", "root", nil); !errors.Is(err, ErrEmptyMutable) {
		t.Fatalf("expected ErrEmptyMutable, got %v", err)
	}
	admin, err := OpenAdminScope("a1", "admin", []string{root})
	if err != nil {
		t.Fatalf("OpenAdminScope: %v", err)
	}
	if _, err := admin.WriteScope().Resolve("/etc/dotnext/config.toml"); err == nil {
		t.Fatal("expected dotNext config directory to remain denied for admin writes")
	}
}

func TestValidateRelativePath(t *testing.T) {
	t.Parallel()

	for _, okCase := range []struct {
		in   string
		want string
	}{
		{"", "."},
		{"wp-content/themes", "wp-content/themes"},
		{"a/./b", "a/b"},
	} {
		got, err := ValidateRelativePath(okCase.in)
		if err != nil || got != okCase.want {
			t.Fatalf("ValidateRelativePath(%q) = (%q, %v), want (%q, nil)", okCase.in, got, err, okCase.want)
		}
	}

	for _, bad := range []string{"../etc/passwd", "/etc/passwd", "foo/../../bar", "ok\x00bad"} {
		if _, err := ValidateRelativePath(bad); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("ValidateRelativePath(%q) error = %v, want ErrUnsafePath", bad, err)
		}
	}
}
