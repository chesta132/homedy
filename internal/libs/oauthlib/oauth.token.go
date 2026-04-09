package oauthlib

import (
	"homedy/internal/libs/cookielib"
	"homedy/internal/libs/cryptolib"
	"homedy/internal/models/payloads"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateState() (string, http.Cookie) {
	state := cryptolib.RandomString(32)
	cookie := cookielib.ToCookie("oauth_state", state, time.Minute*5)
	return state, cookie
}

func ValidateState(c *gin.Context, payload payloads.RequestGithubOAuthCallback) bool {
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != payload.State {
		return false
	}
	return true
}
