package routes

import (
	"homedy/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterWebsocket(group *gin.RouterGroup) {
	terminal := handlers.NewWsTerminal()

	group.GET("/terminal", terminal.Handle)
}
