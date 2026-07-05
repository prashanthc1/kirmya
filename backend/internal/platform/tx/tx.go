package tx

import (
	"context"
	"database/sql"
)

type txKey struct{}

// DBTX is the common interface shared by *sql.DB and *sql.Tx.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// WithTx embeds a transaction in the context.
func WithTx(ctx context.Context, sqlTx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, sqlTx)
}

// GetTx retrieves the transaction from the context, if present.
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	sqlTx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return sqlTx, ok
}

// GetExecutor returns the active transaction if present in the context;
// otherwise, it falls back to the base *sql.DB.
func GetExecutor(ctx context.Context, db *sql.DB) DBTX {
	if sqlTx, ok := GetTx(ctx); ok {
		return sqlTx
	}
	return db
}

// TxDB wraps *sql.DB and implements DBTX by dynamically choosing between the
// transaction stored in the context and the base database connection pool.
type TxDB struct {
	db *sql.DB
}

// NewTxDB wraps the database connection pool in a transaction-aware executor.
func NewTxDB(db *sql.DB) *TxDB {
	return &TxDB{db: db}
}

func (t *TxDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return GetExecutor(ctx, t.db).ExecContext(ctx, query, args...)
}

func (t *TxDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return GetExecutor(ctx, t.db).PrepareContext(ctx, query)
}

func (t *TxDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return GetExecutor(ctx, t.db).QueryContext(ctx, query, args...)
}

func (t *TxDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return GetExecutor(ctx, t.db).QueryRowContext(ctx, query, args...)
}

// TxManager coordinates the transaction lifecycle.
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a transaction coordinator.
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// RunInTx runs the provided callback inside a database transaction. If the context
// already carries an active transaction, it is reused. If the callback returns
// an error, the transaction is rolled back; otherwise, it commits.
func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := GetTx(ctx); ok {
		return fn(ctx)
	}

	sqlTx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = sqlTx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(WithTx(ctx, sqlTx)); err != nil {
		_ = sqlTx.Rollback()
		return err
	}

	return sqlTx.Commit()
}
