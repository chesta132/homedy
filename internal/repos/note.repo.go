package repos

import (
	"context"
	"fmt"
	"homedy/config"
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err = r.db.WithContext(ctx).Model(new(models.Note)).Select("UserID").Where("id = ?", id).Pluck("user_id", &userID).Error
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
	err = r.db.WithContext(ctx).Model(new(models.Note)).Select("UserID").Where("id IN ?", ids).Pluck("user_id", &userID).Error
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

func (r *Note) GetNotesWithPayload(ctx context.Context, userID string, payload payloads.RequestGetManyNote) (notes []models.Note, err error) {
	query := r.db.WithContext(ctx).Model(new(models.Note)).Where("user_id = ?", userID).
		// config.LIMIT_RESOURCE_PER_PAGINATION + 1 to cursor pagination
		Offset(payload.Offset).Limit(config.LIMIT_RESOURCE_PER_PAGINATION + 1).
		// safe as long as payload.Sort is validated
		Order(fmt.Sprintf("updated_at %s, id %s", payload.Sort, payload.Sort))
	if payload.Recycled {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}

	err = query.Find(&notes).Error
	return
}

func (r *Note) RestoreNotesWithPayload(ctx context.Context, payload payloads.RequestRestoreManyNote) (notes []models.Note, err error) {
	tx := r.db.WithContext(ctx).Unscoped().Model(new(models.Note)).
		Clauses(clause.Returning{}).
		Order(fmt.Sprintf("updated_at %s, id %s", payload.Sort, payload.Sort)).
		Where("id IN ?", payload.IDs).
		Update("deleted_at", nil).Scan(&notes)
	notes = slicelib.Map(notes, func(idx int, note models.Note) models.Note { note.Decrypt(); return note })

	err = tx.Error
	if tx.RowsAffected == 0 && err == nil {
		err = gorm.ErrRecordNotFound
	}
	return
}
