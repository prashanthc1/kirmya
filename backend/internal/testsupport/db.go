//go:build integration

// Package testsupport provides helpers for integration/repository tests that run
// against a real PostgreSQL database. These tests are compiled and run only with
// the `integration` build tag (e.g. `go test -tags=integration ./...`) and skip
// when TEST_DATABASE_URL is not set.
package testsupport

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"workspace-app/internal/platform/migrate"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenTestDB connects to TEST_DATABASE_URL, applies all migrations, and returns
// a database with a clean slate (user-owned data truncated). Seeded reference
// data (roles, communities) is preserved. The test is skipped when
// TEST_DATABASE_URL is unset.
func OpenTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping test db: %v", err)
	}

	if err := migrate.RunFrom(db, filepath.Join(moduleRoot(t), "migrations")); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	// Clean slate for user-owned data; FK cascades wipe profiles/jobs/referrals/etc.
	if _, err := db.Exec("TRUNCATE users RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate users: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// InsertUser creates a minimal active user and returns its id, for use as a
// foreign-key fixture in repository tests.
func InsertUser(t *testing.T, db *sql.DB, email, fullName string) string {
	t.Helper()
	var id string
	err := db.QueryRow(`
		INSERT INTO users (email, full_name, email_verified, status)
		VALUES ($1, $2, true, 'active') RETURNING id`, email, fullName).Scan(&id)
	if err != nil {
		t.Fatalf("insert user %s: %v", email, err)
	}
	return id
}

// moduleRoot walks up from the working directory to the directory containing
// go.mod (the backend module root), so migrations resolve regardless of which
// package's tests are running.
func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}
