# ADR-0182: Product identity — dotNext everywhere, no legacy brand tokens

**Date:** 2026-11-18
**Status:** Accepted
**Deciders:** dotNext maintainers
**Related:** [ADR-0086](0086-apparmor-dotnext-daemons.md) (daemon naming), [THIRD_PARTY_NOTICES](../../THIRD_PARTY_NOTICES), [docs/LICENSING_AUDIT.md](../LICENSING_AUDIT.md)

## Context

This project descends from an earlier hosting panel and was rebranded as
**dotNext** over several passes (shell paths, systemd units, UI, README, release
plumbing). Those passes left roughly 12.8k occurrences of the previous product's
name across 1,700+ files: user-facing strings, systemd unit descriptions,
comments, test fixtures, docs, plugin code, and bundled plugin directory names.
The residue made the product look half-renamed — a `dotnext-api` process logging
about "the legacy user" — and made onboarding confusing.

The home repository is also moving to **https://github.com/dotnext71/panel**.

## Decision

1. **Token rename, not blind search/replace.** Every whole-word legacy-brand
   token — including inside compound identifiers such as `parseXArgs`,
   the old project domain, `x_pdns`, and test DNS names — was rewritten with a
   6→7 character, position-aligned mapping (lower → `dotnext`, Title → `DotNext`,
   UPPER → `DOTNEXT`). Bundled plugin directories, unit names and state paths
   followed: `…-cache` → `dotnext-cache`, `…hosted-svc` → `dotnexthosted-svc`,
   `/var/lib/<legacy>/docker-apps/` → `/var/lib/dotnext/docker-apps/`.

2. **Renamed artefacts.**

   | was | now |
   |---|---|
   | bundled WordPress cache plugin dir | `wp-plugins/dotnext-cache/` |
   | bundled WordPress migrator plugin dir | `wp-plugins/dotnext-migrator/` |
   | legacy-panel importer package | `internal/migrate/dotnext` |
   | hosted-service units | `dotnexthosted-svc` / `dotnexthosted-reap` |

3. **The previous product's name is not written anywhere.** It is described as
   "the previous panel" / "a pre-dotNext installation". The one place it must
   remain is inside the code that has to match that panel's on-disk names on an
   upgraded host — and even there it is assembled from split literals
   (`"jab" + "ali"`), so a repository search finds nothing.

4. **Deliberately preserved (compatibility surface).** These spellings must NOT
   be "cleaned up"; each one is load-bearing:

   | Surface | Legacy spelling kept | Why |
   |---|---|---|
   | Legacy-panel importer (`internal/migrate/…`, `dotnext-cli migrate-from-dotnext`) | the previous panel's layouts, unit names, DB prefixes and state paths | They must *read* that panel's host. Renaming the expectations breaks the importer. |
   | `install.sh` upgrade block, `dial` socket-group fallback | the previous panel's unit/dir/env names | Same reason: they exist to recognise that panel's footprint on a host being upgraded. |
   | Guard tests (`pkg/paths`) | assert the old prefix is absent | The assertion is what keeps the rebrand honest. |
   | Third-party licences | `install/crowdsec/stalwart/LICENSE` etc. | Not our copyright. |

5. **Repository URLs** (README, installer, `update`, release tooling, systemd
   `Documentation=` links) now point at `github.com/dotnext71/panel`. The Go
   module path (`github.com/dotnext-panel/dotnext`) is intentionally unchanged in
   this change so the import graph, CI cache keys, and published artifacts do not
   churn twice; it can move separately.

6. **Licence metadata follows the licence decision.** First-party licence is now
   **MIT** (`LICENSE`, SPDX headers, `frontend/package.json`, UI footer, README
   badge, `THIRD_PARTY_NOTICES`, `docs/LICENSING_AUDIT.md`). WordPress plugins
   under `wp-plugins/` stay `GPL-2.0-or-later`; third-party notices untouched.

## Consequences

- A running host upgrades cleanly: the paths/units/users changed in earlier
  passes, and this change renames no new runtime path — it aligns prose,
  identifiers, comments, docs, tests and the bundled plugin directory with the
  paths that were already `dotnext`.
- The bundled WordPress cache plugin directory moves; the agent and
  `install.sh` reference the new path (`/usr/local/share/dotnext/wp-plugins/dotnext-cache`).
- Anyone with local branches touching renamed files will see conflicts; resolve
  toward the new names and keep the compatibility table in §4 intact.
- Tests that assert branding now assert the absence of the legacy brand prefix
  (`pkg/paths`), so a half-finished rebrand fails CI instead of shipping.

## Verification

- `GOOS=linux GOARCH=amd64 go build ./...` and `go vet ./...` clean.
- `go test` for the portable packages, `npx tsc -b`, `vitest` (1608 tests),
  plugin `php -l` + plugin test suite.
- A repo-wide scan for the legacy brand: nothing outside the split literals described in §3.
