package database

import (
	"homedy/internal/models"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{}, &models.Revoke{},
	)
}
