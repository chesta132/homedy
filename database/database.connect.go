package database

import (
	"homedy/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DB_PATH))
	if err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
