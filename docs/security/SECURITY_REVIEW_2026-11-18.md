# Security Review — dotNext Control Panel

**Date:** 2026-11-18
**Scope:** full repository — Go control plane (`internal/`, `core/`, `cmd/`, `agent/`),
React SPA (`frontend/`), installer (`install.sh`, `bootstrap.sh`), migrations, WordPress plugins.
**Method:** automated scanning (Go vuln DB, `npm audit`, secret grep, authorization route
audit, duplicate-list sweep) plus manual review of every trust boundary.
**Verdict:** the security architecture is **strong** — the two defects found were both
regressions from the rebrand, and both are now fixed.

---

## 1. What was found and fixed

### 1.1 HIGH — admin File Manager lost its legacy deny-list entries

`internal/filesafe/filesafe.go` holds `DefaultAdminDeniedPrefixes`, the hard deny-list
for the root-scoped admin File Manager (GH #1184). The rebrand rewrote the list's
`legacy path → new path` pairs into **duplicates**:

```
before rename:                     after rename:
"/etc/jabali"     -> "/etc/dotnext"   "/etc/dotnext",   // x3
"/run/jabali"     -> "/run/dotnext"   "/run/dotnext",   // x2
"/var/lib/jabali" -> "/var/lib/dotnext" "/var/lib/dotnext", // x2
```

The new path stayed (repeated), the legacy path **vanished**. On a host upgraded from the
previous panel those paths can still hold live config, secrets, agent sockets and
binaries — the installer migrates them only when it runs, and an unmigrated or
partially-migrated host keeps them. That is precisely what the deny-list exists to
prevent an admin UI from reading or rewriting.

36 entries, 27 unique: 8 duplicated, 4 legacy families missing.

**Fixed:** list de-duplicated and grouped by purpose; a separate
`legacyControlPlanePrefixes` restores the four families; a new `AdminDeniedPrefixes()`
returns both so a future entry cannot be added to one list and forgotten in the other.
Three regression tests pin it (legacy families denied, no duplicates, no overlap with
the mutable-root allow-list).

### 1.2 MEDIUM — reachable advisory in the shipped SPA

`npm audit --omit=dev` reported **1 high** production vulnerability:
`follow-redirects@1.16.0` (GHSA-q3jh-wrc6-jgq8 and friends — array-form headers bypass
cross-origin credential stripping; sensitiveHeaders regex escaping can expose credential
headers). It is reachable: `axios` is the only consumer and the SPA calls with
`withCredentials: true` (`frontend/src/kratos.ts`).

**Fixed:** npm override to `^1.16.3` (satisfies axios's `^1.16.0`).
`npm audit --omit=dev` → **0 vulnerabilities**.

### 1.3 INFO — devDependency findings, not shipped

The other 18 findings are all devDependencies (`vitest`, `happy-dom`, `eslint`,
`playwright`, `minimatch`, `picomatch`, `js-yaml`, `fdir`, `source-map-js`,
`tinyglobby`, `brace-expansion`). None is imported by `frontend/src` except `vitest`'s
own test APIs, none reaches the Vite bundle, and none runs on a production host.
Tracked, not urgent; upgrade them on the next dependency pass.

### 1.4 INFO — a stale Go advisory list, not applicable

The Go vulnerability DB lists advisories against 18 of our modules, but only
`github.com/docker/docker@28.5.2` shows advisories that *affect the pinned version*.
Docker is **test-only**: `go mod why` traces it to `dhui/dktest`, reachable from
`golang-migrate`'s test binary. It appears in zero production binaries
(`go list -deps ./cmd/... | grep docker/docker` → 0).

---

## 2. Architecture — what is done well

Worth recording, because these are the properties that should not be regressed.

### Trust boundaries are explicit and kernel-enforced

| Boundary | Control | Assessment |
|---|---|---|
| Panel ⇄ agent | Unix socket, `SO_PEERCRED` uid gate + `AllowedUIDs` allow-list; `peerErr != nil ⇒ adminCapable = false` (fail closed) | Correct |
| Public ⇄ panel | nginx → unix socket; `RequireLocalhost()` rejects any request carrying proxy headers, then checks the kernel peer uid | Correct, and it documents a real past exploit it closed |
| Panel ⇄ MariaDB/Postgres | unix socket, peer-auth; `sudo -u postgres psql` scoped | Correct |
| File operations | `openat2(2)` `RESOLVE_BENEATH`, base fd opened `O_NOFOLLOW`, **fails closed on ENOSYS** rather than falling back to the unsafe logical path | Excellent — this is the right design for a root-scoped file manager |
| Tenant outbound URLs | `ssrfguard`: resolve → reject private/loopback/link-local/CGNAT/reserved → **pin the TCP connection to the validated IP** closing the DNS-rebind TOCTOU, on both IPv4 and IPv6 | Excellent |

### Authorization is centralised, not sprinkled

866 route registrations, 54 `/admin` groups. Every one of them carries
`RequireAdmin()` either inline or on the group. The route audit found **zero**
production admin routes without a guard; the only unguarded matches were four test
harnesses and one route guarded by the strictly stronger `RequireLocalhost()`.

`RequireAdmin()` returns **401 when claims are absent** — it cannot be bypassed by
mounting a handler before the auth middleware — and `RequireOwner()` compares against
the panel ULID, not the Kratos UUID, so ownership checks did not silently break at the
identity cutover.

### Injection surfaces are parameterised or validated

The seven `psql -c` interpolations are all backed by an allow-list regex
(`^[a-zA-Z][a-zA-Z0-9_-]{0,63}$`) or a hardcoded list. MariaDB literals go through
`EscapeMariaDBLiteral`. There is **no** `exec.Command("bash","-c", userInput)` anywhere;
28 `sudo` invocations pass argv vectors, never a concatenated string.

### The SSH sandbox hides what it should

`dotnext-ssh-shell` builds a bwrap jail that ro-binds `/usr`, `/etc/ssl`, `/etc/resolv.conf`
and selected tool paths, and **omits** `/etc/shadow` and `/etc/dotnext` — exposing only
`/etc/dotnext/snuffleupagus` (rules the tenant needs read). Panel secrets stay invisible
to tenant shells by construction, not by file permission tuning.

### Frontend

`vite` bundle contains only `dependencies` (not devDependencies). The panel CSP is
`default-src 'self'` with the support-chat origin added only when an admin opts in.
The one broad directive, `script-src 'unsafe-inline'`, is inherited and should be
tightened to a nonce or hash (see §4).

### Secrets

No secret is committed. `install.sh` **generates** every credential at install time
(`openssl rand -hex 32` for the Redis panel token, JWT secret, PowerDNS password, …).
`.env.example` holds placeholders only; real `.env*` is gitignored. `gitleaks` runs in CI.

Error responses: 400s carry fixed validator strings (`"too large (max 8 chars)"`), 500s
carry `"internal"` / `"internal_error"`. No stack traces or driver errors reach a client.

---

## 3. Dependency posture

| Ecosystem | Production advisories | Notes |
|---|---:|---|
| Go | 0 applicable | 18 modules have DB entries; only docker's affect the pin, and it is test-only |
| npm (shipped SPA) | **0** after override | 1 was reachable and is fixed |
| npm (dev) | 18 | vitest/happy-dom/eslint/playwright; never ship |

Recommendation: add `npm audit --omit=dev` and a Go vuln check to CI as a **hard gate**.
`govulncheck` currently crashes on this module size (a `golang.org/x/tools` SSA bug, not a
project fault); the DB-query approach used here works and could be scripted instead.

---

## 4. Open items — not defects, but the next improvements

| # | Item | Why it matters | Effort |
|---|---|---|---|
| 1 | Tighten panel CSP `script-src` from `'unsafe-inline'` to a per-request nonce | `'unsafe-inline'` removes most of CSP's XSS protection; it is the single biggest remaining hardening win | Medium — needs a nginx `sub_filter` nonce or a build-time hash for the SPA |
| 2 | Add `npm audit --omit=dev` + Go vuln check to CI | The two issues found here were both caught by hand, not by the pipeline | Small |
| 3 | Rate limiter is **in-memory** (`golang.org/x/time/rate` per process) | Correct for the single-process panel; it silently becomes useless if a second instance is ever added behind the load balancer | Small — document the constraint, or move to the Redis store if scaled |
| 4 | Upgrade the 18 devDependency findings | Build-machine only today, but the toolchain still executes them | Small |
| 5 | `follow-redirects` override should be re-checked on each axios bump | Axios may move the requirement | Trivial, just remember it |
| 6 | Tenant personal data: confirm retention + deletion covers backups and caches | GDPR-style obligations; the panel deletes accounts, but archived backups are a separate path | Medium — audit `backup`/`restore` retention against `user delete` |

None of these is a vulnerability. Item 1 is the one worth scheduling.

---

## 5. Verification performed

```
go build ./... (linux/amd64)        clean
go vet ./...                        clean
gofmt -l                            clean
internal/filesafe, core/files,
internal/api tests                  pass
tsc -b                              clean
vitest                              258 files / 1608 tests pass
vite build                          ok
npm audit --omit=dev                0 vulnerabilities
git grep -iE "password|secret|token|api_?key\s*[:=]\s*\"[^\"]{12,}\""
                                    only openssl-generated locals in install.sh
admin route authorization audit     0 unguarded production routes
duplicate-list sweep                fixed (was 1 — the deny-list)
```

A regression test now fails if the deny-list ever gets a duplicate again, or if a
control-plane family stops being denied.
