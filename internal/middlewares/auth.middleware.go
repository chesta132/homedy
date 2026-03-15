package middlewares

import (
	"errors"
	"homedy/config"
	"homedy/internal/libs/authlib"
	"homedy/internal/libs/replylib"
	"net/http"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (mw *Auth) protected(c *gin.Context) (claims authlib.Claims, err error) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	accessCookie, err := c.Request.Cookie(config.ACCESS_TOKEN_KEY)
	if err == nil {
		claims, err = authlib.ParseAccessToken(accessCookie.Value)
		if err == nil {
			return
		}
	}

	// read & validate refresh token
	refreshCookie, err := c.Request.Cookie(config.REFRESH_TOKEN_KEY)
	if err != nil {
		err = errors.New("no refresh token provided")
		return
	}

	// TODO: check if token is revoked

	claims, err = authlib.ParseRefreshToken(refreshCookie.Value)
	if err != nil {
		return
	}

	// token cookies
	cookies := []http.Cookie{}

	// update access token
	cookies = append(cookies, authlib.CreateAccessCookie(claims.UserID, claims.RememberMe))
	// rotate refresh token
	if claims.RotateAt.Before(time.Now()) {
		cookies = append(cookies, authlib.CreateRefreshCookie(claims.UserID, claims.RememberMe))
	}

	rp.SetCookies(cookies...)

	c.Set("userID", claims.UserID)
	return
}

// ensureAuthenticated to ensure protected middleware run, return false if not authenticated and aborted
func (mw *Auth) ensureAuthenticated(c *gin.Context, rp *reply.Reply) bool {
	// if prev middleware is protecting return true
	if _, exists := c.Get("userID"); exists {
		return true
	}

	// validate token
	_, err := mw.protected(c)
	if err != nil {
		rp.Error(replylib.CodeUnauthorized, err.Error()).FailJSON()
		c.Abort()
		return false
	}

	return true
}

// Protected middleware basic auth
func (mw *Auth) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))
		if mw.ensureAuthenticated(c, rp) {
			c.Next()
		}
	}
}
