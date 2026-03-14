package repos

import (
	"context"

	"gorm.io/gorm"
)

type archivable[T any] struct {
	db *gorm.DB
}

func (r *archivable[T]) Archive(ctx context.Context, where any, args ...any) (err error) {
	result, err := gorm.G[T](r.db).Where(where, args...).Delete(ctx)
	if result == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *archivable[T]) Restore(ctx context.Context, where any, args ...any) (err error) {
	result, err := gorm.G[T](r.db.Unscoped()).Where(where, args...).Update(ctx, "deleted_at", nil)
	if result == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *archivable[T]) Delete(ctx context.Context, where any, args ...any) (rowsAffected int, err error) {
	return gorm.G[T](r.db.Unscoped()).Where(where, args...).Delete(ctx)
}

func (r *archivable[T]) DeleteByID(ctx context.Context, id string) (success bool, err error) {
	rowsAffected, err := r.Delete(ctx, "id = ?", id)
	return rowsAffected > 0, err
}

func (r *archivable[T]) DeleteByIDs(ctx context.Context, ids []string) (rowsAffected int, err error) {
	return r.Delete(ctx, "id IN ?", ids)
}
