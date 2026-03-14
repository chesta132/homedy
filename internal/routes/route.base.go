package routes

import (
	"homedy/internal/repos"

	"gorm.io/gorm"
)

type Router struct {
	db    *gorm.DB
	repos *repos.Repos
}

func New(db *gorm.DB, repos *repos.Repos) *Router {
	return &Router{db, repos}
}
