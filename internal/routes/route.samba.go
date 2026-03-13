package routes

import (
	"homedy/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterSamba(group *gin.RouterGroup) {
	h := handlers.NewSamba()

	group.POST("/", h.AddShare)
	group.GET("/", h.GetAll)
	group.GET("/:name", h.GetOne)
	group.PUT("/:name", h.UpdateOne)
	group.DELETE("/:name", h.DeleteOne)

	registerSambaConfig(group.Group("/config"))
}

func registerSambaConfig(group *gin.RouterGroup) {
	h := handlers.NewSamba()

	group.GET("/", h.GetConfiguration)
}
