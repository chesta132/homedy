package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterSamba(group *gin.RouterGroup) {
	sambaSvc := services.NewSamba()
	h := handlers.NewSamba(sambaSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	group.POST("/", h.CreateShare)
	group.GET("/", h.GetShares)
	group.GET("/:name", h.GetShare)
	group.PUT("/:name", h.UpdateShare)
	group.DELETE("/:name", h.DeleteShare)

	group.Use(amw.AppProtected(middlewares.AppProtectQuery()))

	group.POST("/backup", h.Backup)
	group.POST("/restore", h.Restore)

	rt.registerSambaConfig(group.Group("/config"))
}

func (rt *Router) registerSambaConfig(group *gin.RouterGroup) {
	sambaSvc := services.NewSamba()
	h := handlers.NewSamba(sambaSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	group.Use(amw.AppProtected(middlewares.AppProtectQuery()))

	group.GET("/", h.GetConfig)
	group.PUT("/", h.UpdateConfig)
}
