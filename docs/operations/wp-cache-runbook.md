# WordPress Cache — Operator Troubleshooting & Recovery (GH #628)

Ties together every moving part of the DotNext WordPress cache. Two layers:

- **Object cache** — Redis via the `dotnext-cache` plugin (drop-ins + phpredis).
- **Page cache** — nginx FastCGI micro-cache (per-domain vhost, ADR-0108).

## Source of truth

The panel owns everything. Runtime config lives in **wp-config.php constants**
stamped by `setWPConfigCacheConstants` (panel-agent/`wordpress_cache.go`) inside
the `// BEGIN/END DotNext WP Cache` marker block. The drop-ins read those
constants, NEVER the DB option (they load before WP boots). The plugin never
writes wp-config (wp.org compliance). If a constant is wrong, fix it in the
panel and re-enable cache — do not hand-edit wp-config.

Constants: `DOTNEXT_CACHE_SOCKET/DB/PASSWORD/USER/PREFIX` (connection),
`DOTNEXT_CACHE_MAXTTL/PAGE_CACHE/PAGE_TTL/MAXMEMORY_MB/DISABLED` (behavior).

## Fast diagnostics

```
# object cache health + stats (as the tenant)
sudo -u <user> wp --path=<docroot> dotnext-cache diagnose
sudo -u <user> wp --path=<docroot> dotnext-cache stats     # {connected,driver,keys,...}
sudo -u <user> wp --path=<docroot> dotnext-cache verify     # live Redis read/write

# page cache — is nginx serving HITs?
curl -ksI --resolve <domain>:443:127.0.0.1 https://<domain>/ | grep -i x-dotnext-cache
```

Panel: the app row's **Cache settings** drawer shows the object-cache stats card
(connected / driver / keys, + admin-only server hit-ratio/memory) and the page
hit ratio. Admin **Cache overview** (`/dotnext-admin/cache`) ranks domains by
page-cache hit ratio (low-hit first).

## Symptom → cause → fix

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| `dotnext-cache diagnose` = "not connected" | Redis socket perms / tenant not in `dotnext-redis-clients` group | re-enable cache (re-runs `usermod -aG` + restarts the per-user FPM master); check `/run/redis/redis.sock` group |
| `driver: pure-PHP` (not phpredis) | `php-redis` missing on this PHP version | `dotnext-cli update` (installs + reloads per-user FPM, #606) |
| `stats` connects but `keys: 0` | cache enabled but not warmed | Warm cache action / auto-warm on enable (#615) |
| Page cache always MISS on `/blog?utm=…` | (fixed #629) gate now on `$uri` | ensure vhost re-rendered (`dotnext-cli update` + reconcile) |
| Homepage never HITs on repeat | a cookie/plugin forcing bypass | run the advisor (#620) — it verifies HIT + flags risky endpoints |
| Wrong tenant's data after clone/migrate | source `DOTNEXT_CACHE_*` carried over | clone/migration/restore strip it (#621); re-enable cache to re-stamp |
| Redis using too much memory | one site's object cache dominating | set a per-install Redis memory budget (#612) → `wp dotnext-cache trim` enforces it |
| Cross-tenant cache eviction via purge | (fixed #630) purge validates ownership from the nginx vhost | n/a |

## Purge

- Manual: app-row **Purge cache** (page cache) — targeted md5-path unlink (#631).
- Automatic: the WP plugin drops a request to `/run/dotnext-wp-purge`; the agent
  watcher validates ownership against the host's **nginx vhost** docroot (#630)
  before purging. A tenant cannot purge another tenant's host.
- After a purge, the domain auto-warms (#632) if auto-warm is on.

## Full recovery (object cache broken)

1. `wp dotnext-cache diagnose` → identify the failing layer.
2. Toggle cache OFF then ON in the panel (re-stamps constants, re-provisions the
   Redis ACL user, re-grants socket group, restarts the per-user FPM master).
3. `wp dotnext-cache stats` → confirm `connected:true, driver:phpredis`.
4. `wp dotnext-cache verify` → confirm live read/write.

## Nothing to lose

The cache is best-effort (Redis is `allkeys-lru`; the page cache is a
micro-cache). Flushing or rebuilding never loses site data — worst case is a
brief cold-cache TTFB spike until it re-warms.
