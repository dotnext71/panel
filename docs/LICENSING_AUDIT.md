# dotNext Licensing & Attribution Audit Report

Date: 2026-11-18  
Target License: **MIT License (`MIT`)**

> **2026-11-18 update:** the first-party licence was moved from `GPL-3.0-only` to
> `MIT` (project decision after the fork-authorisation review in §1).
> Component-level metadata (SPDX headers, `frontend/package.json`, UI footer,
> README badge, [THIRD_PARTY_NOTICES](../THIRD_PARTY_NOTICES)) was migrated with it.

## 1. Authorization & Relicensing Basis

dotNext is derived from the Jabali Panel codebase (`upstream-jabali` @ `f0a46afa6`), originally distributed under `AGPL-3.0`. Explicit authorization from the original project rights holders permits:

1. Full product rebranding as **dotNext**.
2. Removal of original product branding and structural reorganization of the codebase.
3. Relicensing of eligible project source code from `AGPL-3.0` to **`MIT`**.

## 2. Contributor & Ownership Audit

Git commit history across the repository (`git shortlog -sn --all`) was inspected prior to license migration:

| Author | Commits | Audit Assessment | Resolution |
|---|---:|---|---|
| Primary upstream author (`Shuki Vaknin` / `shuki` / `shukivaknin`) | 6,245 | Covered by rights-holder relicensing authorization | Eligible for `MIT` |
| dotNext maintainer (`Khalequzzaman`) | 11 | Authoritative fork commits (M1–M12) | `MIT` |
| Automated bots (`dependabot[bot]`, `jabali-catalog-bot`) | 4 | Mechanical dependency/digest bumps (`go.mod`, catalog pins) | Non-copyrightable / covered by upstream |
| `Viktor Szépe` | 1 (`aec341cb9e`) | Modified only 8 lines in `README.md` (`<picture>` / `#gh-*-mode-only` image markup) | Zero code authored; superseded by dotNext `README.md` |

## 3. Component License Classification

### 3.1 Relicensed to `MIT`
- Control plane API & CLI (`panel-api/`, migrating to `cmd/dotnext-api/`, `cmd/dotnext-cli/`, `core/`, `internal/`)
- Privileged host agent & helper binaries (`panel-agent/`, migrating to `cmd/dotnext-agent/`, `agent/`)
- Shared Go packages (`pkg/`, `internal/`, `agentwire/`)
- Web UI (`panel-ui/`, migrating to `frontend/`)
- Installer TUI, modular shell libraries & scripts (`installer/`, `install.sh`, `deploy/`, `install/` except third-party notices listed below)
- Auxiliary first-party Go services (`hostedsvc/`, `support-claim/`)

### 3.2 Retained Under Original Licenses (Not Relicensed)
- **`install/crowdsec/stalwart/LICENSE`**: MIT License (`Copyright (c) 2026 Nethanja`). Preserved untouched and recorded in [THIRD_PARTY_NOTICES](../THIRD_PARTY_NOTICES).
- **`wp-plugins/dotnext-cache/` & `wp-plugins/dotnext-migrator/`**: Retained under **`GPL-2.0-or-later`** per WordPress plugin distribution standards (`wp-plugins/dotnext-cache/LICENSE` preserved).
- **Third-Party Libraries (`go.mod`, `support-claim/go.mod`, `frontend/package.json`, composer locks)**: Retain their respective upstream OSI licenses.
- **External Runtime Binaries & Packages** (`Stalwart Mail` AGPL-3.0, `Bulwark` AGPL-3.0, `Ory Kratos` Apache-2.0, `CrowdSec` MIT, `Snuffleupagus` LGPL-3.0, `Restic` BSD-2-Clause, `YARA-X` BSD-3-Clause, `LMD` GPL-2.0, `phpMyAdmin` GPL-2.0, `Adminer` Apache-2.0/GPL-2.0, OS packages): Downloaded/installed on target hosts and remain subject to upstream licenses only.

## 4. Verification Checklist

- [x] Official MIT Licence text placed verbatim in [LICENSE](../LICENSE).
- [x] Legacy `NOTICE` (referencing AGPL-3.0) replaced by [THIRD_PARTY_NOTICES](../THIRD_PARTY_NOTICES).
- [x] `SPDX-License-Identifier: MIT` headers on primary Go entrypoints, `install.sh`, UI `package.json`, UI footer, and docs.
- [x] Third-party MIT notice (`install/crowdsec/stalwart/LICENSE`) and `wp-plugins/` `GPL-2.0-or-later` verified intact.
