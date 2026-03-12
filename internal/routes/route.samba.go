package routes

import (
	"homedy/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterSamba(group *gin.RouterGroup) {
	h := handlers.NewSamba()
	group.POST("/", h.AddShare)
	group.GET("/", h.GetAll)
}
