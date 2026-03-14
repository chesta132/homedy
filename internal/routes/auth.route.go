package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterAuth(group *gin.RouterGroup) {
	authSvc := services.NewAuth(rt.repos.User())
	h := handlers.NewAuth(authSvc)

	group.POST("/signup", h.SignUp)
}
