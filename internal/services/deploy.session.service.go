package services

import (
	"homedy/internal/middlewares"
	"homedy/internal/models"
	"time"

	"github.com/google/uuid"
)

func (s *ContextedDeploy) CreateSession() (string, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()
	key := "deploy:session:" + id
	err = s.rdb.HSet(s.ctx, key, models.DeploySession{UserID: userID}).Err()
	if err != nil {
		return "", err
	}

	err = s.rdb.Expire(s.ctx, key, time.Hour).Err()
	if err != nil {
		return "", err
	}

	return id, nil
}
