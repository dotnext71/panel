<?php
/**
 * Plugin Name: DotNext Cache — Advanced Cache Drop-in
 * Description: Redis-backed full-page cache for the dotnext panel. Installed by the DotNext Cache plugin; do not edit by hand.
 * Version: 1.0.0
 *
 * Copied to wp-content/advanced-cache.php. WordPress includes this only when
 * the WP_CACHE constant is true (wp-config.php). It runs before the rest of
 * WordPress, so it must be self-locating and must never fatal: any problem
 * just means the page is generated normally.
 *
 * @package DotNext_Cache
 */

defined( 'ABSPATH' ) || exit;

// advanced-cache must not interfere with the installer or CLI.
if ( ( defined( 'WP_INSTALLING' ) && WP_INSTALLING ) || ( defined( 'WP_CLI' ) && WP_CLI ) ) {
	return;
}

if ( ! defined( 'DOTNEXT_CACHE_PLUGIN_DIR' ) ) {
	// A drop-in lives in wp-content/, not in the plugin, so it cannot resolve
	// the plugin path with plugin_dir_path( __FILE__ ). The plugin stamps the
	// absolute path into the placeholder below at install; the only fallback is
	// the canonical WP_PLUGIN_DIR for a manual install under the standard slug.
	$dotnext_cache_stamped = '__DOTNEXT_CACHE_PLUGIN_DIR__';
	if ( 0 !== strpos( $dotnext_cache_stamped, '__DOTNEXT' ) && is_dir( $dotnext_cache_stamped ) ) {
		define( 'DOTNEXT_CACHE_PLUGIN_DIR', $dotnext_cache_stamped );
	} elseif ( defined( 'WP_PLUGIN_DIR' ) && is_dir( WP_PLUGIN_DIR . '/dotnext-cache' ) ) {
		define( 'DOTNEXT_CACHE_PLUGIN_DIR', WP_PLUGIN_DIR . '/dotnext-cache' );
	} else {
		define( 'DOTNEXT_CACHE_PLUGIN_DIR', '' );
	}
}

if ( '' === DOTNEXT_CACHE_PLUGIN_DIR ) {
	return;
}

$dotnext_cache_page_engine = DOTNEXT_CACHE_PLUGIN_DIR . '/includes/class-page-cache.php';
if ( ! is_readable( $dotnext_cache_page_engine ) ) {
	return;
}

require_once DOTNEXT_CACHE_PLUGIN_DIR . '/includes/lib.php';
require_once $dotnext_cache_page_engine;

try {
	$dotnext_cache_page_cache = new DotNext_Cache_Page_Cache();
	$dotnext_cache_page_cache->run();
} catch ( \Throwable $e ) {
	// Never let the page cache take the site down.
	return;
}
