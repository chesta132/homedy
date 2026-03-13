package database

import (
	"fmt"
	"homedy/config"

	"gorm.io/gorm"
)

func createDSN(host, user, pass, name, port string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, pass, name, port)
}

func appDSN() string {
	return createDSN(
		config.DB_HOST,
		config.DB_USER,
		config.DB_PASS,
		config.DB_NAME,
		config.DB_PORT,
	)
}

func isSuperUser(db *gorm.DB) (bool, error) {
	var isSuper bool
	err := db.Raw(`
        SELECT usesuper FROM pg_user WHERE usename = current_user
    `).Scan(&isSuper).Error
	return isSuper, err
}
