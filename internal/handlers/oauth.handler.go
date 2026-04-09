package handlers

import (
	"homedy/config"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/oauthlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"
	"net/http"
	"net/url"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type OAuth struct {
	oAuthSvc *services.OAuth
}

func NewOAuth(oAuthSvc *services.OAuth) *OAuth {
	return &OAuth{oAuthSvc}
}

// @Summary      Redirect to github oauth
// @Tags         OAuth
// @Success      307
// @Router			 /oauth/github   [get]
func (h *OAuth) BindGithub(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	state, cookie := oauthlib.CreateState()
	url := oauthlib.GithubOAuthConfig.AuthCodeURL(state)

	rp.SetCookies(cookie).Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary      Callback of github oauth
// @Tags         OAuth
// @Produce      json
// @Param				 payload  query	payloads.RequestGithubOAuthCallback	true "state of oauth callback"
// @Response     307
// @Router			 /oauth/github/callback [get]
func (h *OAuth) CallbackGithub(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	u, _ := url.Parse(config.FRONTEND_URL)
	u.Path = "/account"

	payload, err := ginlib.BindAndValidate[payloads.RequestGithubOAuthCallback](c.ShouldBindQuery)
	if err != nil {
		u.RawQuery = url.Values{"github": {"error"}, "error": {err.Error()}}.Encode()
		rp.Redirect(http.StatusTemporaryRedirect, u.String())
		return
	}

	err = h.oAuthSvc.AttachContext(c).CallbackGithub(payload)
	if err != nil {
		u.RawQuery = url.Values{"github": {"error"}, "error": {err.Error()}}.Encode()
		rp.Redirect(http.StatusTemporaryRedirect, u.String())
		return
	}

	u.RawQuery = url.Values{"github": {"connected"}}.Encode()
	rp.Redirect(http.StatusTemporaryRedirect, u.String())
}
