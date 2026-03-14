package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterSamba(group *gin.RouterGroup) {
	sambaSvc := services.NewSamba()
	h := handlers.NewSamba(sambaSvc)

	group.POST("/", h.CreateShare)
	group.GET("/", h.GetShares)
	group.GET("/:name", h.GetShare)
	group.PUT("/:name", h.UpdateShare)
	group.DELETE("/:name", h.DeleteShare)

	rt.registerSambaConfig(group.Group("/config"))
}

func (rt *Router) registerSambaConfig(group *gin.RouterGroup) {
	sambaSvc := services.NewSamba()
	h := handlers.NewSamba(sambaSvc)

	group.GET("/", h.GetConfig)
	group.PUT("/", h.UpdateConfig)
}
