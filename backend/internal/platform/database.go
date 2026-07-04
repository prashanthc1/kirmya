package platform

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenDatabase connects to PostgreSQL using the pgx stdlib driver.
//
// Connection is configured via DATABASE_URL, e.g.
//
//	postgres://user:pass@localhost:5432/kirmya?sslmode=disable
//
// If DATABASE_URL is unset, it is assembled from POSTGRES_* env vars with
// sensible local-development defaults.
func OpenDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			getenvDefault("POSTGRES_USER", "postgres"),
			os.Getenv("POSTGRES_PASSWORD"),
			getenvDefault("POSTGRES_HOST", "127.0.0.1"),
			getenvDefault("POSTGRES_PORT", "5432"),
			getenvDefault("POSTGRES_DB", "kirmya"),
			getenvDefault("POSTGRES_SSLMODE", "disable"),
		)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
