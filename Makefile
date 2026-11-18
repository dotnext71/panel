.PHONY: help build run test test-short test-coverage coverage-check lint fmt vet tidy clean \
	test-ui test-e2e test-all ui-install ui-build demo-guard

GO         := go
API_PKG    := ./cmd/... ./core/... ./internal/... ./database/...
AGENT_PKG  := ./agent/... ./cmd/...
WIRE_PKG   := ./agentwire/...
ALL_PKG    := $(API_PKG) $(AGENT_PKG) $(WIRE_PKG)
BIN        := ./bin/dotnext-api
CLI_BIN    := ./bin/dotnext-cli
AGENT_BIN  := ./bin/dotnext-agent
COVER      := coverage.out
MIN_COV    := 80

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Compile dotNext binaries (panel + cli + agent + helpers)
	mkdir -p bin
	$(GO) build -o $(BIN) ./cmd/dotnext-api
	$(GO) build -o $(CLI_BIN) ./cmd/dotnext-cli
	$(GO) build -o $(AGENT_BIN) ./cmd/dotnext-agent
	$(GO) build -o bin/dotnext-installer ./cmd/dotnext-installer
	$(GO) build -o bin/dotnext-webdav ./cmd/dotnext-webdav
	$(GO) build -o bin/dotnext-ssh-shell ./cmd/dotnext-ssh-shell
	$(GO) build -o bin/dotnext-mailhook ./cmd/dotnext-mailhook
	$(GO) build -o bin/dotnext-sendmail ./cmd/dotnext-sendmail

demo-guard: ## JAB-159: prod build must EXCLUDE demo code; demo build must INCLUDE it
	mkdir -p bin
	$(GO) build -tags demo -o bin/dotnext-demo ./cmd/dotnext-api
	$(GO) build -o bin/dotnext-prod ./cmd/dotnext-api
	@grep -aq 'demo_mode' bin/dotnext-demo || { echo 'FAIL: demo build missing demo_mode marker'; exit 1; }
	@if grep -aq 'demo_mode' bin/dotnext-prod; then echo 'FAIL: prod build leaked demo_mode symbol'; exit 1; fi
	@echo 'demo-guard OK: prod excludes demo_mode, demo build includes it'
	$(GO) test -tags demo -count=1 ./internal/middleware/

run: ## Run the panel server (dev)
	$(GO) run ./cmd/dotnext-api

test: ## Run all Go tests across the workspace
	$(GO) test -race -count=1 $(ALL_PKG)

test-short: ## Run only fast unit tests (skip integration)
	$(GO) test -race -count=1 -short $(ALL_PKG)

test-coverage: ## Run tests with coverage (internal packages only)
	$(GO) test -race -count=1 -coverprofile=$(COVER) -covermode=atomic -coverpkg=./internal/... ./internal/...
	$(GO) tool cover -func=$(COVER) | tail -n 1

test-integration: ## Run integration tests (requires DOTNEXT_TEST_DATABASE_URL + real MariaDB)
	$(GO) test -race -count=1 -tags=integration -coverprofile=$(COVER) -covermode=atomic -coverpkg=./internal/... ./internal/...
	$(GO) tool cover -func=$(COVER) | tail -n 1

coverage-check: ## Fail if combined (unit+integration) coverage below MIN_COV
	@if [ -z "$$DOTNEXT_TEST_DATABASE_URL" ]; then \
		echo "coverage-check requires DOTNEXT_TEST_DATABASE_URL (real MariaDB)"; \
		echo "  set it, or run 'make test-coverage' for unit-only (no gate)"; \
		exit 1; \
	fi
	@$(MAKE) --no-print-directory test-integration
	@pct=$$($(GO) tool cover -func=$(COVER) | awk '/total:/ {gsub("%","",$$3); print $$3}'); \
	awk -v p="$$pct" -v m="$(MIN_COV)" 'BEGIN { if (p+0 < m+0) { printf "coverage %s%% below %s%%\n", p, m; exit 1 } else { printf "coverage %s%% OK\n", p } }'

lint: ## Run golangci-lint across the workspace + install.sh phantom-function lint
	golangci-lint run $(ALL_PKG)
	tools/lint-install-sh.sh

fmt: ## Format all Go code
	$(GO) fmt $(ALL_PKG)

vet: ## Run go vet
	$(GO) vet $(ALL_PKG)

tidy: ## Tidy module deps
	$(GO) mod tidy

clean: ## Remove build artefacts
	rm -rf bin $(COVER)

# ---------- frontend ----------

UI_DIR := frontend

ui-install: ## Install frontend npm deps (clean, reproducible)
	cd $(UI_DIR) && rm -rf node_modules && npm ci --no-audit --no-fund

ui-build: ## Build the frontend SPA (required before E2E — tests run against dist/)
	cd $(UI_DIR) && npm run build

test-ui: ## Run frontend unit tests (vitest)
	cd $(UI_DIR) && npx vitest run

test-e2e: ui-build ## Run Playwright E2E suite against the built SPA
	cd $(UI_DIR) && npx playwright test --project=chromium --reporter=list

test-all: test test-ui test-e2e ## Run everything: Go tests + vitest + Playwright

# ---------- aa-smoke (M40.1 AppArmor profile verifier) ----------

aa-smoke: ## Verify every loaded dotnext AppArmor profile reaches its declared sockets (M40.1)
	bash tools/aa-smoke/run.sh
