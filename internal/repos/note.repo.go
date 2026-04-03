package repos

import (
	"context"
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
	err = gorm.G[models.Note](r.db).Select("UserID").Where("id = ?", id).Scan(ctx, &userID)
	return
}

func (r *Note) GetUserIDsByIDs(ctx context.Context, ids []string) (userID []string, err error) {
	err = gorm.G[models.Note](r.db).Select("UserID").Where("id IN ?", ids).Scan(ctx, &userID)
	return
}

func (r *Note) GetUserIDsByIDsUnscoped(ctx context.Context, ids []string) (userID []string, err error) {
	err = gorm.G[models.Note](r.db.Unscoped()).Select("UserID").Where("id IN ?", ids).Scan(ctx, &userID)
	return
}
