<?php
/**
 * DotNext Cache — uninstall cleanup.
 *
 * Runs when the plugin is deleted from the WordPress admin. Removes drop-ins,
 * the generated config file, and the stored settings option. Cached data in
 * Redis is scoped to this site's prefix and expires under LRU; we also make a
 * best-effort prefix-scoped purge here (never FLUSHDB — DB 1 is shared).
 *
 * @package DotNext_Cache
 */

defined( 'WP_UNINSTALL_PLUGIN' ) || exit;

require_once __DIR__ . '/includes/lib.php';
require_once __DIR__ . '/includes/class-settings.php';
require_once __DIR__ . '/includes/class-dropin-manager.php';

// Best-effort: purge this site's keys from Redis.
$dotnext_cache_cfg    = DotNext_Cache_Config::load();
$dotnext_cache_client = new DotNext_Cache_Client( $dotnext_cache_cfg );
if ( $dotnext_cache_client->connect() ) {
	$dotnext_cache_client->delete_by_pattern( $dotnext_cache_cfg['prefix'] . '*' );
	$dotnext_cache_client->close();
}

// Remove drop-ins + generated config.
$dotnext_cache_mgr = new DotNext_Cache_Dropin_Manager( __DIR__ );
$dotnext_cache_mgr->remove();
DotNext_Cache_Settings::delete_config_file();

// Drop settings.
delete_option( DotNext_Cache_Settings::OPTION );
