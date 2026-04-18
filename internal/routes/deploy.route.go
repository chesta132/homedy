package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterDeploy(group *gin.RouterGroup, composeService api.Compose) {
	deploySvc := services.NewDeploy(rt.repos.OAuth(), rt.repos.DeployRepo(), rt.repos.DeployLog(), rt.repos.DeploySession(), composeService)
	h := handlers.NewDeploy(deploySvc)

	dmw := middlewares.NewDeploy(rt.repos.RDB(), rt.repos.OAuth())

	group.Use(dmw.Protected())
	group.POST("/new", h.CreateSession)

	group.Use(dmw.SessionProtected())
	// user's github repositories
	group.GET("/:session/repos", h.GetRepos)

	// user's selected repository in cache
	group.POST("/:session/selected-repo", h.SetSelectedRepo)
	group.GET("/:session/selected-repo", h.GetSelectedRepo)

	// user's github branches of repository
	group.GET("/:session/repos/:id/branches", h.GetBranches)
}
