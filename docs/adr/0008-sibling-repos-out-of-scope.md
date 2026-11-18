# 0008 — Sibling repos are out-of-scope for panel

## Status
Accepted — 2026-04-16

## Context
The DotNext ecosystem includes several independent services: dotnext-manager (central control plane), dotnext-security (WAF daemon), dotnext-isolator (nspawn container manager), dotnext-tunnel (WSS sidecar), and Bulwark webmail (Next.js JMAP client). Each has its own release cadence and audience.

This decision enforces scope: this repo (dotnext2) develops the panel only. Sibling repos are installed as optional addons but not developed here.

## Decision
`dotnext-manager`, `dotnext-security`, `dotnext-isolator`, `dotnext-tunnel`, and Bulwark webmail are separate repositories. The panel may define wire protocols (bridge endpoints, socket proxies) for integration but does NOT implement the sibling daemons. Their APIs are pinned via `~/projects/dotnext-shared/CONTEXT.md`.

## Consequences

### Positive
- Clear repo boundaries; independent release cycles
- Reduced merge conflict and test surface
- Teams can own separate services independently
- Easier rollback (one service doesn't block another)

### Negative
- Inter-service API changes require coordination
- Multi-node features (manager, tunnel) are deferred
- Install complexity grows (optional addons)

### Neutral
- `install.sh` may invoke sibling installers (as subshells)

## Alternatives considered

- **Monorepo**: Rejected — couples release cycles, hard to scale teams
- **Vendored subtrees**: Rejected — merge hell and duplicate CI

## References
- `~/projects/dotnext-shared/CONTEXT.md` — wire protocol pinning
- `install.sh` — addon installation order
