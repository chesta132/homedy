package handlers

import (
	"fmt"
	"homedy/config"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"
	"net/http"
	"net/url"

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

	err = h.authSvc.AttachContext(c).SignUp(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(nil).CreatedJSON()
}

func (h *Auth) SignUpApproval(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSignUpApproval](c.ShouldBindQuery)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user, err := h.authSvc.AttachContext(c).SignUpApproval(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Redirect(
		http.StatusTemporaryRedirect,
		fmt.Sprintf("%s/signup/review-approval?username=%s&email=%s&action=%s", config.FRONTEND_URL, url.QueryEscape(user.Username), url.QueryEscape(user.Email), payload.Action),
	)
}

func (h *Auth) SignUpApprovalStatus(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSignUpApprovalStatus](c.ShouldBindQuery)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	res, err := h.authSvc.AttachContext(c).SignUpApprovalStatus(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(res).OkJSON()
}

func (h *Auth) SignIn(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindJSONAndValidate[payloads.RequestSignIn](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user, cookies, err := h.authSvc.AttachContext(c).SignIn(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(user).SetCookies(cookies...).OkJSON()
}

func (h *Auth) SignOut(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	cookies := h.authSvc.AttachContext(c).SignOut()
	rp.Success(nil).SetCookies(cookies...).OkJSON()
}

func (h *Auth) Me(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	user, err := h.authSvc.AttachContext(c).Me()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(user).OkJSON()
}
