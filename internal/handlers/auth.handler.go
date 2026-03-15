package handlers

import (
	"homedy/internal/libs/authlib"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	authSvc *services.Auth
}

func NewAuth(authSvc *services.Auth) *Auth {
	return &Auth{authSvc}
}

func (h *Auth) SignUp(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindJSONAndValidate[payloads.RequestSignUp](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user := payload.ToUser()
	err = h.authSvc.AttachContext(c).SignUp(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(user).SetCookies(authlib.CreateTokenCookie(user.ID, payload.RememberMe)...).CreatedJSON()
}
