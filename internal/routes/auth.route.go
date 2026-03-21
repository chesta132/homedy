package routes

import (
	"homedy/internal/handlers"
	"homedy/internal/libs/mail"
	"homedy/internal/middlewares"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterAuth(group *gin.RouterGroup) {
	authSvc := services.NewAuth(rt.repos.User(), rt.repos.Revoke(), mail.NewAppMailer())
	h := handlers.NewAuth(authSvc)
	amw := middlewares.NewAuth(rt.repos.Revoke())

	group.POST("/signup", h.SignUp)
	group.POST("/signin", h.SignIn)
	// get for email
	group.GET("/signup/approval", h.SignUpApproval)
	group.GET("/signup/approval-status", h.SignUpApprovalStatus)

	group.Use(amw.Protected())
	group.POST("/signout", h.SignOut)
	group.GET("/me", h.Me)
}
