# ADR-0182: Product rename — every first-party `Jabali` token becomes `DotNext`

**Date:** 2026-11-18
**Status:** Accepted
**Deciders:** dotNext maintainers
**Related:** [ADR-0086](0086-apparmor-dotnext-daemons.md) (daemon naming), [THIRD_PARTY_NOTICES](../../THIRD_PARTY_NOTICES), [docs/LICENSING_AUDIT.md](../LICENSING_AUDIT.md)

## Context

The project was forked from Jabali Panel and rebranded as **dotNext** over several
passes (shell paths, systemd units, UI, README, release plumbing). Those passes
left ~12.8k `jabali` / `Jabali` / `JABALI` occurrences spread across 1,700+ files:
user-facing strings, systemd unit descriptions, comments, test fixtures, docs,
plugin code, and the bundled WordPress plugin directory names. The residue made
the product look half-renamed (a `dotnext-api` process logging about "the jabali
user") and made onboarding confusing.

The home repository is also moving to **https://github.com/dotnext71/panel**.

## Decision

1. **Token rename, not blind search/replace.** Every whole-word `jabali` token
   (including inside compound identifiers such as `parseJabaliArgs`,
   `jabali-panel.com`, `jabali_pdns`, `smoke.jabalitest.com`) is rewritten with
   a 6→7 character, position-aligned mapping:

   | legacy | replacement |
   |---|---|
   | `jabali` | `dotnext` |
   | `Jabali` | `DotNext` (camel-case when it starts a mixed identifier) |
   | `JABALI` | `DOTNEXT` |

   `jabali-cache` therefore becomes `dotnext-cache`,
   `jabalihosted-svc` becomes `dotnexthosted-svc`,
   `/var/lib/jabali/docker-apps/` becomes `/var/lib/dotnext/docker-apps/`.

2. **Renamed artefacts.**

   | was | now |
   |---|---|
   | `wp-plugins/jabali-cache/` (bundled) | `wp-plugins/dotnext-cache/` |
   | `wp-plugins/jabali-migrator/` | `wp-plugins/dotnext-migrator/` |
   | `internal/migrate/jabali` (Go package) | `internal/migrate/dotnext` |
   | `jabalihosted-svc` / `jabalihosted-reap` units | `dotnexthosted-svc` / `dotnexthosted-reap` |
   | `infrastructure/docker-apps/jabali-sounder/` | `infrastructure/docker-apps/dotnext-sounder/` |

3. **Brand-name prose only.** "Jabali Panel" as the *origin* of this fork is
   still written as "Jabali Panel" in provenance documents (this ADR,
   `THIRD_PARTY_NOTICES`, `docs/LICENSING_AUDIT.md`, historical ADRs and
   `CHANGELOG.md`), because those records describe the upstream project, not us.
   The legacy repository name and account appear only where git history
   is quoted.

4. **Deliberately preserved (compatibility surface).** These spellings must NOT
   be "cleaned up"; each one is load-bearing:

   | Surface | Legacy spelling kept | Why |
   |---|---|---|
   | Migrate-from-legacy importers (`internal/migrate/…`, `dotnext-cli migrate-from-dotnext`) | `jabali` panel layouts, unit names, DB prefixes, `/var/lib/jabali*`, `/opt/jabali*` | They *read* a legacy source panel. Renaming the expectations breaks the importer. |
   | WordPress.org plugin slug | `jabali-cache` (SDV/`wp plugin update`) | The directory slug is immutable; changing it orphans installed sites. The plugin's *name* is "dotNext Cache". |
   | Dotted DNS names in tests/examples (`smoke.jabali.site`, `ns1.jabali.local`) | unchanged gate domain | Owned by an external test zone; not ours to rename. |
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
- Tests that assert branding now assert the absence of `jabali`
  (`pkg/paths`), so a half-finished rebrand fails CI instead of shipping.

## Verification

- `GOOS=linux GOARCH=amd64 go build ./...` and `go vet ./...` clean.
- `go test` for the portable packages, `npx tsc -b`, `vitest` (1608 tests),
  plugin `php -l` + plugin test suite.
- A repo-wide scan for remaining `jabali` except the surface in §4.
