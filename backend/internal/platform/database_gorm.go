package platform

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenGormDB wraps the existing *sql.DB pgx pool inside a GORM session,
// allowing both ORM operations and raw pgx queries to share the same connection pool.
func OpenGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
}
