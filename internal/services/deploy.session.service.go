package services

import (
	"homedy/internal/middlewares"
	"homedy/internal/models/payloads"
)

func (s *ContextedDeploy) CreateSession() (string, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return "", err
	}
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return "", err
	}
	user, _, err := client.Users.Get(s.ctx, "")
	if err != nil {
		return "", err
	}

	id, err := s.deploySessionRepo.CreateSession(s.ctx, userID, user.GetLogin())
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *ContextedDeploy) InvalidateSession(payload payloads.TemplateWithSession) error {
	return s.deploySessionRepo.RemoveSession(s.ctx, payload.Session)
}
