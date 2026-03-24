package handlers

import (
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

// @Summary      Creates new account
// @Description	 Create new account and ask [config.MAIL_OWNER] to register approval, user will be notified if account successfully review
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 payload  body	payloads.RequestSignUp	true	"data of new account"
// @Success      201  		{object}  replylib.Envelope
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/signup [post]
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

// @Summary      Allow or deny sign up request
// @Description	 Owner only allow or deny sign up request. Identified by app secret
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 X-APP-SECRET header string true "app secret authentication for access"
// @Param				 payload  body	payloads.RequestSignUpApproval	true	"data of sign up request"
// @Success      200  		{object}  replylib.Envelope{data=models.User}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/signup/approval [patch]
func (h *Auth) SignUpApproval(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSignUpApproval](c.ShouldBindBodyWithJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user, err := h.authSvc.AttachContext(c).SignUpApproval(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(user).OkJSON()
}

// @Summary      Check sign up status
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 payload  body	payloads.RequestSignUpApprovalStatus	true "identity to check status"
// @Success      200  		{object}  replylib.Envelope{data=payloads.ResponseSignUpApprovalStatus}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/signup/approval-status [get]
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

// @Summary      Sign in onto existing account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param				 payload  body	payloads.RequestSignIn	true	"sign in identity"
// @Success      200  		{object}  replylib.Envelope{data=models.User}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/signin [post]
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

// @Summary      Sign out account for one device (request device). Have to sign in first
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  		{object}  replylib.Envelope
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/signout [post]
func (h *Auth) SignOut(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	cookies := h.authSvc.AttachContext(c).SignOut()
	rp.Success(nil).SetCookies(cookies...).OkJSON()
}

// @Summary      Get signed in account information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  		{object}  replylib.Envelope{data=models.User}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /auth/me [get]
func (h *Auth) Me(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	user, err := h.authSvc.AttachContext(c).Me()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(user).OkJSON()
}
