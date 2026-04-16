package middlewares

import (
	"errors"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v68/github"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Deploy struct {
	rdb       *redis.Client
	oAuthRepo *repos.OAuth
}

func NewDeploy(rdb *redis.Client, oAuthRepo *repos.OAuth) *Deploy {
	return &Deploy{rdb, oAuthRepo}
}

// use auth.Protected first
func (mw *Deploy) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))

		ctx := c.Request.Context()
		userID, err := GetUserID(c)
		if err != nil {
			replylib.HandleError(err, rp)
			return
		}

		oAuthResource, err := mw.oAuthRepo.GetFirst(ctx, "user_id = ?", userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rp.Error(replylib.CodeForbidden, "GitHub account is not linked").FailJSON()
			return
		}
		if err != nil {
			replylib.HandleError(err, rp)
			return
		}

		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oAuthResource.AccessToken})
		client := github.NewClient(oauth2.NewClient(ctx, ts))
		c.Set("gh-client", client)

		c.Next()
	}
}

// return [ErrMiddlewareSkipped] on error
func GetGithubClient(c *gin.Context) (*github.Client, error) {
	ghClientIfc, _ := c.Get("gh-client")
	ghClient, ok := ghClientIfc.(*github.Client)
	if !ok {
		return nil, ErrMiddlewareSkipped
	}
	return ghClient, nil
}

// use auth.Protected first
func (mw *Deploy) SessionProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		rp := replylib.Client.Use(adapter.AdaptGin(c))
		payload, err := ginlib.BindAndValidate[payloads.TemplateWithSession](c.ShouldBindUri)
		if err != nil {
			replylib.HandleError(err, rp)
			return
		}

		userID, err := GetUserID(c)
		if err != nil {
			replylib.HandleError(err, rp)
			return
		}

		var rUserID string
		err = mw.rdb.HGet(c.Request.Context(), "deploy:session:"+payload.Session, "userId").Scan(&rUserID)
		if err != nil {
			replylib.HandleError(err, rp)
			return
		}

		if rUserID != userID {
			rp.Error(replylib.CodeForbidden, "Invalid session").FailJSON()
			return
		}

		c.Next()
	}
}
