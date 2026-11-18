<?php
/**
 * DotNext Cache — sample wp-config.php constants.
 *
 * You normally do NOT need any of these: the plugin auto-detects the dotnext
 * panel Redis (unix socket /run/redis/redis.sock, database 1).
 *
 * The early drop-ins read these `DOTNEXT_CACHE_*` constants from wp-config.php
 * (plus built-in defaults) — they load before WordPress boots, so they do NOT
 * read the admin/CLI settings option. There is NO generated
 * wp-content/dotnext-cache-config.php (an older build wrote one; it is now
 * legacy and is deleted on save/activation/uninstall).
 *
 * Define any of these in wp-config.php (ABOVE the "stop editing" line) to
 * override the defaults.
 *
 * @package DotNext_Cache
 */

// Master kill switch (useful for staging): true disables all caching.
// define( 'DOTNEXT_CACHE_DISABLED', false );

// Connection — unix socket is the dotnext default and recommended path.
// define( 'DOTNEXT_CACHE_SOCKET', '/run/redis/redis.sock' );
// define( 'DOTNEXT_CACHE_DB', 1 );          // ADR-0059: database 1 reserved for WP.
// define( 'DOTNEXT_CACHE_PASSWORD', '' );    // the panel socket has no AUTH.

// TCP fallback (only if you run a non-dotnext Redis). Setting a host switches
// the scheme to tcp automatically.
// define( 'DOTNEXT_CACHE_HOST', '127.0.0.1' );
// define( 'DOTNEXT_CACHE_PORT', 6379 );

// Explicit per-site key prefix. Leave unset to auto-derive a stable, unique
// prefix from the site URL + path (recommended on shared Redis).
// define( 'DOTNEXT_CACHE_PREFIX', 'mysite' );

// Cap object-cache key TTL in seconds (0 = no expiry; rely on Redis LRU).
// define( 'DOTNEXT_CACHE_MAXTTL', 0 );

// Enable the OPTIONAL full-page cache (GH #625). WP_CACHE makes WordPress load
// advanced-cache.php; DOTNEXT_CACHE_PAGE_CACHE turns the dotnext page cache on and
// DOTNEXT_CACHE_PAGE_TTL sets its TTL. These constants are honored by the early
// drop-ins (same as the object-cache constants above).
// NOTE: dotnext already ships an nginx fastcgi microcache (ADR-0108); only
// enable the WP page cache if you disabled that at the panel.
// define( 'WP_CACHE', true );
// define( 'DOTNEXT_CACHE_PAGE_CACHE', true ); // WP full-page cache on/off
// define( 'DOTNEXT_CACHE_PAGE_TTL', 300 );    // page-cache TTL (seconds)
