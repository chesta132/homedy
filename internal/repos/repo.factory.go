package repos

import "gorm.io/gorm"

type Repos struct {
	db   *gorm.DB
	user *User
}

func New(db *gorm.DB) *Repos {
	return &Repos{db: db}
}

func (r *Repos) User() *User {
	if r.user == nil {
		r.user = NewUser(r.db)
	}
	return r.user
}
