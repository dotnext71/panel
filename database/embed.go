// SPDX-License-Identifier: MIT

// Package database embeds SQL schema migrations for runtime application.
package database

import (
	"embed"
	"io/fs"
)

//go:embed all:migrations/*.sql
var MigrationsFS embed.FS

// MigrationsRoot returns a sub-FS rooted at migrations/ so golang-migrate
// sees migration filenames directly.
func MigrationsRoot() (fs.FS, error) {
	return fs.Sub(MigrationsFS, "migrations")
}
