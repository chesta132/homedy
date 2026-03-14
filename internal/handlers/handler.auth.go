package handlers

import (
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/validatorlib"
	"homedy/internal/payloads"
	"homedy/internal/services/auth"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{db}
}

func (h *Auth) SignUp(c *gin.Context) {
	svc := auth.NewService(h.db, c.Request.Context())
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := validatorlib.BindJSONAndValidate[payloads.RequestSignUp](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user := payload.ToUser()
	err = svc.SignUp(&user, payload.RememberMe)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(user).CreatedJSON()
}
