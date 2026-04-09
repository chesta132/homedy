package repos

import (
	"homedy/internal/models"

	"gorm.io/gorm"
)

type OAuth struct {
	db *gorm.DB
	create[models.OAuth]
	read[models.OAuth]
	update[models.OAuth]
	delete[models.OAuth]
}

func NewOAuth(db *gorm.DB) *OAuth {
	return &OAuth{db, create[models.OAuth]{db}, read[models.OAuth]{db}, update[models.OAuth]{db}, delete[models.OAuth]{db}}
}

func (r *OAuth) DB() *gorm.DB {
	return r.db
}

func (r *OAuth) WithContext(tx *gorm.DB) *OAuth {
	return NewOAuth(tx)
}
