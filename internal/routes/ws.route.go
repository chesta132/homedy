package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterWebsocket(group *gin.RouterGroup) {
	terminalSvc := services.NewTerminal()
	h := handlers.NewWsTerminal(terminalSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	group.GET("/terminal", amw.AppProtected(middlewares.SecretGetterWs()), h.Handle)
}
