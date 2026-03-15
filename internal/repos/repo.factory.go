package repos

import "gorm.io/gorm"

type Repos struct {
	db     *gorm.DB
	user   *User
	revoke *Revoke
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

func (r *Repos) Revoke() *Revoke {
	if r.revoke == nil {
		r.revoke = NewRevoke(r.db)
	}
	return r.revoke
}
