package repos

import (
	"context"
	"homedy/config"
	"homedy/internal/models"
	"homedy/internal/models/payloads"

	"gorm.io/gorm"
)

type Note struct {
	db *gorm.DB
	create[models.Note]
	read[models.Note]
	update[models.Note]
	archivable[models.Note]
}

func NewNote(db *gorm.DB) *Note {
	return &Note{db, create[models.Note]{db}, read[models.Note]{db}, update[models.Note]{db}, archivable[models.Note]{db}}
}

func (r *Note) DB() *gorm.DB {
	return r.db
}

func (r *Note) WithContext(tx *gorm.DB) *Note {
	return NewNote(tx)
}

func (r *Note) GetUserIDByID(ctx context.Context, id string) (userID string, err error) {
	err = r.db.WithContext(ctx).Select("UserID").Where("id = ?", id).Pluck("user_id", &userID).Error
	if userID == "" && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *Note) GetUserIDByRecycledID(ctx context.Context, id string) (userID string, err error) {
	err = r.db.WithContext(ctx).Unscoped().Model(new(models.Note)).Select("UserID").Where("id = ? AND deleted_at IS NOT NULL", id).Pluck("user_id", &userID).Error
	if userID == "" && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *Note) GetUserIDsByIDs(ctx context.Context, ids []string) (userID []string, err error) {
	err = r.db.WithContext(ctx).Select("UserID").Where("id IN ?", ids).Pluck("user_id", &userID).Error
	if len(userID) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *Note) GetUserIDsByRecycledIDs(ctx context.Context, ids []string) (userIDs []string, err error) {
	err = r.db.WithContext(ctx).Unscoped().Model(new(models.Note)).Select("UserID").Where("id IN ? AND deleted_at IS NOT NULL", ids).Pluck("user_id", &userIDs).Error
	if len(userIDs) == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *Note) GetNotesWithPayload(userID string, payload payloads.RequestGetManyNote) (notes []models.Note, err error) {
	// config.LIMIT_RESOURCE_PER_PAGINATION + 1 to cursor pagination
	query := r.db.Where("user_id = ?", userID).Offset(payload.Offset).Limit(config.LIMIT_RESOURCE_PER_PAGINATION + 1)
	if payload.Recycled {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}
	err = query.Find(&notes).Error
	return
}
