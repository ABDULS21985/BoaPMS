package repository

import (
	"context"
	"errors"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"gorm.io/gorm"
)

// PMSRepository extends Repository with PMS-specific query methods.
// Mirrors .NET's IPMSRepo<T> with additional methods beyond the generic repo.
type PMSRepository[T any] struct {
	Repository[T]
}

// NewPMSRepository creates a new PMS-specific repository.
func NewPMSRepository[T any](db *gorm.DB) *PMSRepository[T] {
	return &PMSRepository[T]{
		Repository: Repository[T]{db: db},
	}
}

// TableNoTracking returns a query without change tracking (read-only).
// In GORM this simply returns the base query — GORM doesn't track by default.
func (r *PMSRepository[T]) TableNoTracking(ctx context.Context) *gorm.DB {
	return r.baseQuery(ctx).Model(new(T))
}

// GetAllIncluding retrieves all records with the specified associations preloaded.
func (r *PMSRepository[T]) GetAllIncluding(ctx context.Context, associations ...string) ([]T, error) {
	var results []T
	q := r.baseQuery(ctx)
	for _, assoc := range associations {
		q = q.Preload(assoc)
	}
	err := q.Find(&results).Error
	return results, err
}

// InsertAndSave creates a record and immediately commits (mirrors .NET Insert + SaveChangesAsync).
func (r *PMSRepository[T]) InsertAndSave(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// InsertBatchAndSave creates multiple records and immediately commits.
func (r *PMSRepository[T]) InsertBatchAndSave(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

// UpdateAndSave updates a record and immediately commits.
func (r *PMSRepository[T]) UpdateAndSave(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// GetRecordsWithStatus retrieves records filtered by their status string.
// Mirrors .NET's GetRecordsWithSatus method.
func (r *PMSRepository[T]) GetRecordsWithStatus(ctx context.Context, status enums.Status) ([]T, error) {
	var results []T
	q := r.db.WithContext(ctx).Where("soft_deleted = ?", false)

	if status != enums.StatusAll {
		q = q.Where("status = ?", status.String())
	}

	err := q.Find(&results).Error
	return results, err
}

// FirstOrDefaultWithPreload retrieves the first matching record with preloaded associations.
func (r *PMSRepository[T]) FirstOrDefaultWithPreload(ctx context.Context, associations []string, query interface{}, args ...interface{}) (*T, error) {
	var entity T
	q := r.baseQuery(ctx)
	for _, assoc := range associations {
		q = q.Preload(assoc)
	}
	err := q.Where(query, args...).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

// WhereWithPreload retrieves records matching a condition with preloaded associations.
func (r *PMSRepository[T]) WhereWithPreload(ctx context.Context, associations []string, query interface{}, args ...interface{}) ([]T, error) {
	var results []T
	q := r.baseQuery(ctx)
	for _, assoc := range associations {
		q = q.Preload(assoc)
	}
	err := q.Where(query, args...).Find(&results).Error
	return results, err
}
