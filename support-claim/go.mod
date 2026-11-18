// support-claim: standalone operator-run claim/redeem service for DotNext
// diagnostic bundles. Stdlib-only; its own module so it builds + deploys
// independently of the panel (the root go.mod's `go build ./...` ignores it).
module github.com/dotnext-panel/dotnext/support-claim

go 1.24
