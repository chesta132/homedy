package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterUser(group *gin.RouterGroup) {
	userSvc := services.NewUser(rt.repos.User())
	h := handlers.NewUser(userSvc)

	group.GET("/:id", h.GetUser)
}
