package routes

import (
	"homedy/internal/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	g     *gin.Engine
	db    *gorm.DB
	repos *repos.Repos
}

func New(g *gin.Engine, db *gorm.DB, repos *repos.Repos) *Router {
	return &Router{g, db, repos}
}
