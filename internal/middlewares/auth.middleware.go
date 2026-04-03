package middlewares

import (
	"errors"
	"homedy/config"
	"homedy/internal/libs/authlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/ws"
	"homedy/internal/repos"
	"net/http"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Auth struct {
	revokedRepo *repos.Revoke
}

func NewAuth(revokedRepo *repos.Revoke) *Auth {
	return &Auth{revokedRepo}
}

func (mw *Auth) protected(c *gin.Context) (claims authlib.Claims, err error) {
	ctx := c.Request.Context()
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	accessCookie, err := c.Request.Cookie(config.ACCESS_TOKEN_KEY)
	if err == nil {
		claims, err = authlib.ParseAccessToken(accessCookie.Value)
		if err == nil {
			c.Set("userID", claims.UserID)
			return
		}
	}

	// read & validate refresh token
	refreshCookie, err := c.Request.Cookie(config.REFRESH_TOKEN_KEY)
	if err != nil {
		err = errors.New("no refresh token provided")
		return
	}

	// check if token is revoked
	if revoked, revErr := mw.revokedRepo.GetFirst(ctx, "value = ?", refreshCookie.Value); !errors.Is(revErr, gorm.ErrRecordNotFound) {
		if revErr != nil {
			err = revErr
			return
		}
		err = errors.New(revoked.Reason)
		return
	}

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

type secretGetter func(c *gin.Context) string

func SecretGetterWs() secretGetter {
	return func(c *gin.Context) string {
		protocols := websocket.Subprotocols(c.Request)
		return ws.GetSubprotocolValue(protocols, config.APP_SECRET_WS_SUBPROTOCOL_KEY)
	}
}

func SecretGetterHeader() secretGetter {
	return func(c *gin.Context) string { return c.GetHeader(config.APP_SECRET_HEADER_KEY) }
}

// App protect middleware compare [app_secret] in payload with [config.APP_SECRET]
func (mw *Auth) AppProtected(secretGetter secretGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))

		secret := secretGetter(c)
		if secret != config.APP_SECRET {
			rp.Error(replylib.CodeForbidden, "invalid app secret", reply.WithFields(reply.FieldsError{
				"app_secret": "invalid app secret",
			})).FailJSON()
			c.Abort()
			return
		}

		c.Next()
	}
}

// return [ErrMiddlewareSkipped] on error
func GetUserID(c *gin.Context) (string, error) {
	userIDIfc, _ := c.Get("userID")
	userID, ok := userIDIfc.(string)
	if !ok {
		return "", ErrMiddlewareSkipped
	}
	return userID, nil
}
