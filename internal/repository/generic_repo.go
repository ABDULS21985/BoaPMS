package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository is a generic GORM repository that mirrors .NET's IRepo<T>.
// T must be a struct with a SoftDeleted bool field for soft-delete filtering.
type Repository[T any] struct {
	db *gorm.DB
}

// NewRepository creates a new generic repository.
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// DB returns the underlying GORM instance (for advanced queries).
func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

// baseQuery returns a query with soft-delete filtering applied.
func (r *Repository[T]) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("soft_deleted = ?", false)
}

// GetByID retrieves a record by primary key.
func (r *Repository[T]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	err := r.baseQuery(ctx).First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

// GetByStringID retrieves a record by a string primary key field.
func (r *Repository[T]) GetByStringID(ctx context.Context, pkColumn string, id string) (*T, error) {
	var entity T
	err := r.baseQuery(ctx).First(&entity, fmt.Sprintf("%s = ?", pkColumn), id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

// GetAll retrieves all non-deleted records.
func (r *Repository[T]) GetAll(ctx context.Context) ([]T, error) {
	var results []T
	err := r.baseQuery(ctx).Find(&results).Error
	return results, err
}

// GetAllPaginated retrieves a page of records.
func (r *Repository[T]) GetAllPaginated(ctx context.Context, offset, limit int) ([]T, error) {
	var results []T
	err := r.baseQuery(ctx).Offset(offset).Limit(limit).Find(&results).Error
	return results, err
}

// Count returns the count of non-deleted records.
func (r *Repository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.baseQuery(ctx).Model(new(T)).Count(&count).Error
	return count, err
}

// CountWhere returns the count of records matching conditions.
func (r *Repository[T]) CountWhere(ctx context.Context, query interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := r.baseQuery(ctx).Model(new(T)).Where(query, args...).Count(&count).Error
	return count, err
}

// FirstOrDefault retrieves the first matching record, or nil if none found.
func (r *Repository[T]) FirstOrDefault(ctx context.Context, query interface{}, args ...interface{}) (*T, error) {
	var entity T
	err := r.baseQuery(ctx).Where(query, args...).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

// Exists checks whether any record matches the condition.
func (r *Repository[T]) Exists(ctx context.Context, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := r.baseQuery(ctx).Model(new(T)).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// Where retrieves records matching a condition.
func (r *Repository[T]) Where(ctx context.Context, query interface{}, args ...interface{}) ([]T, error) {
	var results []T
	err := r.baseQuery(ctx).Where(query, args...).Find(&results).Error
	return results, err
}

// WherePaginated retrieves a page of records matching a condition.
func (r *Repository[T]) WherePaginated(ctx context.Context, offset, limit int, query interface{}, args ...interface{}) ([]T, error) {
	var results []T
	err := r.baseQuery(ctx).Where(query, args...).Offset(offset).Limit(limit).Find(&results).Error
	return results, err
}

// Create inserts a new record.
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// CreateBatch inserts multiple records.
func (r *Repository[T]) CreateBatch(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

// Update saves changes to an existing record.
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete performs a soft delete by setting SoftDeleted = true.
func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Model(entity).Update("soft_deleted", true).Error
}

// HardDelete permanently removes a record from the database.
func (r *Repository[T]) HardDelete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Unscoped().Delete(entity).Error
}

// Query returns a GORM query builder with soft-delete filtering for custom queries.
// Usage: repo.Query(ctx).Preload("Relation").Where("x = ?", y).Find(&results)
func (r *Repository[T]) Query(ctx context.Context) *gorm.DB {
	return r.baseQuery(ctx).Model(new(T))
}

// Preload returns a query with eager-loaded associations.
// Equivalent to .NET's include properties pattern.
func (r *Repository[T]) Preload(ctx context.Context, associations ...string) *gorm.DB {
	q := r.baseQuery(ctx)
	for _, assoc := range associations {
		q = q.Preload(assoc)
	}
	return q
}
