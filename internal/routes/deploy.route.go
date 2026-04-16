package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterDeploy(group *gin.RouterGroup) {
	deploySvc := services.NewDeploy(rt.repos.RDB(), rt.repos.OAuth(), rt.repos.DeployRepo(), rt.repos.DeployLog())
	h := handlers.NewDeploy(deploySvc)

	dmw := middlewares.NewDeploy(rt.repos.RDB(), rt.repos.OAuth())

	group.Use(dmw.Protected())
	group.POST("/new", h.CreateSession)

	group.Use(dmw.SessionProtected())
	group.GET("/:session/repos", h.GetRepos)
}
