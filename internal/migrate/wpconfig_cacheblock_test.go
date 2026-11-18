package migrate

import (
	"strings"
	"testing"
)

func TestStripDotNextCacheBlock(t *testing.T) {
	in := `<?php
define( 'DB_NAME', 'wp' );
// BEGIN DotNext WP Cache
if ( ! defined( 'DOTNEXT_CACHE_SOCKET' ) ) define( 'DOTNEXT_CACHE_SOCKET', '/run/redis/redis.sock' );
if ( ! defined( 'DOTNEXT_CACHE_PREFIX' ) ) define( 'DOTNEXT_CACHE_PREFIX', 'src:01ABC' );
// END DotNext WP Cache
require_once ABSPATH . 'wp-settings.php';
`
	out := StripDotNextCacheBlock(in)
	if strings.Contains(out, "DOTNEXT_CACHE_") || strings.Contains(out, "BEGIN DotNext WP Cache") {
		t.Fatalf("cache block not stripped:\n%s", out)
	}
	if !strings.Contains(out, "DB_NAME") || !strings.Contains(out, "wp-settings.php") {
		t.Fatalf("stripped too much:\n%s", out)
	}
	if StripDotNextCacheBlock(out) != out {
		t.Fatal("not idempotent")
	}
}
