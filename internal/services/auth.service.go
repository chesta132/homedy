package services

import (
	"context"
	"errors"
	"homedy/internal/libs/dblib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo *repos.User
}

type ContextedAuth struct {
	Auth
	c   *gin.Context
	ctx context.Context
}

func NewAuth(userRepo *repos.User) *Auth {
	return &Auth{userRepo}
}

func (s *Auth) AttachContext(c *gin.Context) *ContextedAuth {
	return &ContextedAuth{*s, c, c.Request.Context()}
}

func (s *ContextedAuth) SignUp(payload payloads.RequestSignUp) (*models.User, error) {
	// validate email and username
	email, username, err := s.userRepo.GetEmailOrUsername(s.ctx, payload.Email, payload.Password)
	isErrNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !isErrNotFound {
		return nil, err
	}
	if !isErrNotFound {
		fe := make(reply.FieldsError)
		if email == payload.Email {
			fe["email"] = "email already registered"
		}
		if username == payload.Username {
			fe["username"] = "username already registered"
		}
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "email or username already registered",
			Fields:  fe,
		}
	}

	// create user (hash in before create)
	newUser := payload.ToUser()
	if err := s.userRepo.Create(s.ctx, &newUser); err != nil {
		return nil, dblib.GormErrorToReplyError(err, &newUser)
	}

	return &newUser, nil
}
