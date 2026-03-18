package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterAuth(group *gin.RouterGroup) {
	authSvc := services.NewAuth(rt.repos.User(), rt.repos.Revoke())
	h := handlers.NewAuth(authSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	group.POST("/signup", h.SignUp)
	group.POST("/signin", h.SignIn)

	group.Use(amw.Protected())
	group.POST("/signout", h.SignOut)
	group.GET("/me", h.Me)
}
