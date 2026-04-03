package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterNote(group *gin.RouterGroup) {
	noteSvc := services.NewNote(rt.repos.Note())
	h := handlers.NewNote(noteSvc)

	amw := middlewares.NewAuth(rt.repos.Revoke())
	group.Use(amw.Protected())

	group.POST("/", h.CreateOne)

	group.GET("/:id", h.GetOne)
	group.GET("/", h.GetMany)

	group.PUT("/:id", h.UpdateOne)

	group.DELETE("/:id", h.DeleteOne)
	group.DELETE("/", h.DeleteMany)

	group.PATCH("/restore/:id", h.RestoreOne)
	group.PATCH("/restore", h.RestoreMany)
}
