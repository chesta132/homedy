package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterConverter(group *gin.RouterGroup) {
	convSvc := services.NewConverter()
	h := handlers.NewConverter(convSvc)

	group.POST("/multiple", h.ConvertMultiple)
	group.POST("/single", h.ConvertOne)
}
