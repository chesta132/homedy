package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterWebsocket(group *gin.RouterGroup) {
	terminalSvc := services.NewTerminal()
	h := handlers.NewWsTerminal(terminalSvc)

	group.GET("/terminal", h.Handle)
}
