package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(opts ...gorm.Option) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(appDSN()), opts...)
	if err != nil {
		return nil, err
	}

	if ok, _ := isSuperUser(db); !ok {
		return nil, errors.New("connected user is not a super user")
	}

	if err = migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
