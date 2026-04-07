package repos

import (
	"context"
	"homedy/internal/models"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
	create[models.User]
	read[models.User]
	update[models.User]
	archivable[models.User]
}

func NewUser(db *gorm.DB) *User {
	return &User{db, create[models.User]{db}, read[models.User]{db}, update[models.User]{db}, archivable[models.User]{db}}
}

func (r *User) DB() *gorm.DB {
	return r.db
}

func (r *User) WithContext(tx *gorm.DB) *User {
	return NewUser(tx)
}

func (r *User) GetEmailOrUsername(ctx context.Context, emailVal, usernameVal string) (email, username string, err error) {
	user, err := gorm.G[models.User](r.db).
		Select("Email", "Username").
		Where("email = ? OR username = ?", emailVal, usernameVal).
		First(ctx)
	return user.Email, user.Username, err
}
