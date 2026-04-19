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

	// oauth validate
	group.Use(dmw.Protected())

	// session related
	group.POST("/session", h.CreateSession)
	group.Use(dmw.SessionProtected())
	group.DELETE("/session/:session/invalidate", h.InvalidateSession)

	// user's github repositories
	group.GET("/session/:session/repos", h.GetRepos)
	// user's github branches of repository
	group.GET("/session/:session/repos/:id/branches", h.GetBranches)

	// user's selected repository in cache
	group.POST("/session/:session/selected-repo", h.SetSelectedRepo)
	group.GET("/session/:session/selected-repo", h.GetSelectedRepo)

	// user's global env and selected repo env
	group.GET("/session/:session/env", h.GetEnv)
	group.POST("/session/:session/env", h.SetEnv)
}
