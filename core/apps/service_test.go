package apps

import (
	"errors"
	"testing"
)

func TestValidateDeploySpec(t *testing.T) {
	t.Parallel()

	spec, err := ValidateDeploySpec(" My-WordPress ", "/blog/", 512)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Slug != "my-wordpress" || spec.MountURL != "/blog" || spec.MemoryMB != 512 {
		t.Fatalf("unexpected spec: %+v", spec)
	}

	if _, err := ValidateDeploySpec("../bad", "/", 256); !errors.Is(err, ErrInvalidSlug) {
		t.Fatalf("expected ErrInvalidSlug, got %v", err)
	}
	if _, err := ValidateDeploySpec("app-1", "/blog/../secret", 256); !errors.Is(err, ErrInvalidSubpath) {
		t.Fatalf("expected ErrInvalidSubpath, got %v", err)
	}
	if _, err := ValidateDeploySpec("app-1", "/", 16); !errors.Is(err, ErrInvalidMemory) {
		t.Fatalf("expected ErrInvalidMemory, got %v", err)
	}
}
