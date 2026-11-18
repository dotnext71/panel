package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/dotnext-panel/dotnext/agentwire"
	"github.com/dotnext-panel/dotnext/internal/platform/osinfo"
)

// phpVersionListResponse is the output shape for php.version.list.
type phpVersionListResponse struct {
	Versions []string `json:"versions"`
}

func extractVersionFromPath(path string) (string, error) {
	if v, ok := osinfo.ExtractPHPVersionFromPoolDir(path); ok {
		return v, nil
	}
	return "", fmt.Errorf("invalid path format: %s", path)
}

// listInstalledPHPVersions lists all installed PHP versions across Debian
// (/etc/php/*/fpm/pool.d) and Remi EL (/etc/opt/remi/php*/php-fpm.d).
func listInstalledPHPVersions() ([]string, error) {
	matches, err := globPoolDirs()
	if err != nil {
		return nil, fmt.Errorf("glob failed: %w", err)
	}

	versionsMap := make(map[string]bool)
	for _, path := range matches {
		// Verify the path actually exists and is a directory.
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			continue
		}

		version, err := extractVersionFromPath(path)
		if err != nil {
			continue
		}
		versionsMap[version] = true
	}

	// Convert map to sorted slice.
	var versions []string
	for v := range versionsMap {
		versions = append(versions, v)
	}
	sort.Strings(versions)

	return versions, nil
}

func phpVersionListHandler(ctx context.Context, params json.RawMessage) (any, error) {
	// No params expected or required for php.version.list.
	if params != nil && len(params) > 0 {
		var p map[string]interface{}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, &agentwire.AgentError{
				Code:    agentwire.CodeInvalidArgument,
				Message: fmt.Sprintf("failed to parse params: %v", err),
			}
		}
		// Even if params are provided, we ignore them and proceed.
	}

	versions, err := listInstalledPHPVersions()
	if err != nil {
		return nil, &agentwire.AgentError{
			Code:    agentwire.CodeInternal,
			Message: fmt.Sprintf("failed to list PHP versions: %v", err),
		}
	}

	return phpVersionListResponse{
		Versions: versions,
	}, nil
}

func init() {
	Default.Register("php.version.list", phpVersionListHandler)
}
