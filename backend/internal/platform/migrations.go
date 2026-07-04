package platform

import (
	"database/sql"

	"workspace-app/internal/platform/migrate"
)

// RunMigrations applies migrations from the default "migrations" directory.
func RunMigrations(db *sql.DB) error { return migrate.Run(db) }
