package services

import (
	"homedy/internal/middlewares"
)

func (s *ContextedDeploy) CreateSession() (string, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return "", err
	}

	id, err := s.deploySessionRepo.CreateSession(s.ctx, userID)
	if err != nil {
		return "", err
	}

	return id, nil
}
