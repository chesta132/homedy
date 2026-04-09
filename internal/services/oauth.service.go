package services

import (
	"context"
	"homedy/internal/libs/oauthlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/middlewares"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/google/go-github/v68/github"

	"github.com/gin-gonic/gin"
)

type OAuth struct {
	oAuthRepo *repos.OAuth
}

type ContextedOAuth struct {
	OAuth
	c   *gin.Context
	ctx context.Context
}

func NewOAuth(oAuthRepo *repos.OAuth) *OAuth {
	return &OAuth{oAuthRepo}
}

func (s *OAuth) AttachContext(c *gin.Context) *ContextedOAuth {
	return &ContextedOAuth{*s, c, c.Request.Context()}
}

func (s *ContextedOAuth) CallbackGithub(payload payloads.RequestGithubOAuthCallback) error {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return err
	}

	if ok := oauthlib.ValidateState(s.c, payload); !ok {
		return &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid state",
		}
	}

	token, err := oauthlib.GithubOAuthConfig.Exchange(s.ctx, payload.Code)
	if err != nil {
		return err
	}

	httpClient := oauthlib.GithubOAuthConfig.Client(s.ctx, token)
	client := github.NewClient(httpClient)

	user, _, err := client.Users.Get(s.ctx, "")

	newOAuth := models.OAuth{
		AppID:       user.GetID(),
		Username:    user.GetLogin(),
		AccessToken: token.AccessToken,
		UserID:      userID,
	}

	return s.oAuthRepo.Create(s.ctx, &newOAuth)
}
