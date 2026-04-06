package services

import (
	"context"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/gin-gonic/gin"
)

type User struct {
	userRepo *repos.User
}

type ContextedUser struct {
	User
	c   *gin.Context
	ctx context.Context
}

func NewUser(userRepo *repos.User) *User {
	return &User{userRepo}
}

func (s *User) AttachContext(c *gin.Context) *ContextedUser {
	return &ContextedUser{*s, c, c.Request.Context()}
}

func (s *ContextedUser) GetUser(payload payloads.RequestGetUser) (*payloads.ResponseGetUser, error) {
	user, err := s.userRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	return &payloads.ResponseGetUser{Base: user.Base, Username: user.Username, Email: user.Email}, nil
}
