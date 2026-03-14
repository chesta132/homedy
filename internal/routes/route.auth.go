package routes

import (
	"homedy/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuth(group *gin.RouterGroup, db *gorm.DB) {
	h := handlers.NewAuth(db)

	group.POST("/signup", h.SignUp)
}
