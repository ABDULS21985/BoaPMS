package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// RawRepository provides raw SQL query capabilities using sqlx.
// This mirrors .NET's IDapperRepo<T> for performance-critical queries.
type RawRepository[T any] struct {
	db        *sqlx.DB
	tableName string
}

// NewRawRepository creates a new raw SQL repository.
func NewRawRepository[T any](db *sqlx.DB, tableName string) *RawRepository[T] {
	return &RawRepository[T]{db: db, tableName: tableName}
}

// GetAll retrieves all records from the table.
func (r *RawRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var results []T
	query := "SELECT * FROM " + r.tableName
	err := r.db.SelectContext(ctx, &results, query)
	return results, err
}

// GetByID retrieves a single record by integer ID.
func (r *RawRepository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	var result T
	query := "SELECT * FROM " + r.tableName + " WHERE id = $1"
	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// QuerySingle executes a custom SQL query returning a single result.
func (r *RawRepository[T]) QuerySingle(ctx context.Context, sql string, params map[string]interface{}) (*T, error) {
	var result T
	rows, err := r.db.NamedQueryContext(ctx, sql, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(&result); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, nil
}

// QueryList executes a custom SQL query returning multiple results.
func (r *RawRepository[T]) QueryList(ctx context.Context, sql string, params map[string]interface{}) ([]T, error) {
	var results []T
	rows, err := r.db.NamedQueryContext(ctx, sql, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item T
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

// Exec executes a non-query SQL statement (INSERT, UPDATE, DELETE).
func (r *RawRepository[T]) Exec(ctx context.Context, sql string, params map[string]interface{}) (int64, error) {
	result, err := r.db.NamedExecContext(ctx, sql, params)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// RawQuery executes any SQL and scans into the target slice.
// For queries that don't map to the repository's type T.
func RawQuery[R any](db *sqlx.DB, ctx context.Context, sql string, args ...interface{}) ([]R, error) {
	var results []R
	err := db.SelectContext(ctx, &results, sql, args...)
	return results, err
}

// RawQuerySingle executes a SQL query and returns a single result.
func RawQuerySingle[R any](db *sqlx.DB, ctx context.Context, sql string, args ...interface{}) (*R, error) {
	var result R
	err := db.GetContext(ctx, &result, sql, args...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
