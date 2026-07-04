// Package migrate is a dependency-free SQL migration runner. It lives in its own
// leaf package (importing no feature modules) so tests can run migrations
// against a real database without creating an import cycle through platform.
package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Run applies migrations from the default "migrations" directory (relative to
// the working directory).
func Run(db *sql.DB) error { return RunFrom(db, "migrations") }

// RunFrom applies all *.sql migrations from dir. It is idempotent: files already
// recorded in schema_migrations are skipped, so it is safe to call repeatedly
// (including from tests).
func RunFrom(db *sql.DB, dir string) error {
	log.Println("Running migrations...")

	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		filename := filepath.Base(file)
		if err := runMigration(db, file, filename); err != nil {
			return fmt.Errorf("migration failed: %s - %w", filename, err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	const q = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id          bigserial PRIMARY KEY,
		filename    text UNIQUE NOT NULL,
		executed_at timestamptz NOT NULL DEFAULT now()
	)`
	_, err := db.Exec(q)
	return err
}

func runMigration(db *sql.DB, filePath, filename string) error {
	var exists int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", filename).Scan(&exists); err != nil {
		return err
	}
	if exists > 0 {
		log.Printf("Migration already applied: %s", filename)
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Each migration file runs in a single transaction so a failure rolls back cleanly.
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Strip "--" line comments BEFORE splitting on ";" — a comment may itself
	// contain a semicolon, which would otherwise split a statement incorrectly.
	for _, stmt := range strings.Split(stripLineComments(string(content)), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute statement in %s: %w", filename, err)
		}
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (filename) VALUES ($1)", filename); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Migration applied: %s", filename)
	return nil
}

// stripLineComments removes "--" line comments from each line. Our migrations
// never contain "--" inside string literals, so this is safe and avoids
// splitting on semicolons that appear inside comment text.
func stripLineComments(content string) string {
	var b strings.Builder
	for _, line := range strings.Split(content, "\n") {
		if i := strings.Index(line, "--"); i >= 0 {
			line = line[:i]
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}
