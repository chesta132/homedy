package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterOAuth(group *gin.RouterGroup) {
	oAuthSvc := services.NewOAuth(rt.repos.OAuth())
	h := handlers.NewOAuth(oAuthSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	// binds
	group.GET("/github", h.BindGithub)

	// callbacks
	group.Use(amw.Protected())
	group.GET("/github/callback", h.CallbackGithub)
}
