package repos

import (
	"context"
	"homedy/internal/libs/logger"
	"homedy/internal/models"

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
	logger.Debug(userID)
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
