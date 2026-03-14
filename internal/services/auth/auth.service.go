package auth

import (
	"context"
	"homedy/internal/libs/dblib"
	"homedy/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	ctx context.Context
}

func NewService(db *gorm.DB, ctx context.Context) *Service {
	return &Service{db, ctx}
}

func (s *Service) SignUp(newUser *models.User, rememberMe bool) error {
	// validate email and username
	if err := validateEmailAndUsername(s.db, s.ctx, newUser.Email, newUser.Username); err != nil {
		return err
	}

	// create user (hash in before create)
	if err := gorm.G[models.User](s.db).Create(s.ctx, newUser); err != nil {
		return dblib.GormErrorToReplyError(err, newUser)
	}

	return nil
}
