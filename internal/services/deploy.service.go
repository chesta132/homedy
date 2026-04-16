// flow:
// client bind github first
// client send create deployment session -> server create session with caching (redis/pg)
// client check available repositories -> server check it on cache first and fallback to fetch github API
// client select repository -> server send available branches
// client select branch -> server check is valid and cache it and response with available docker services
//
//	validation:
//	- have docker-compose.(yaml/yml)
//	- have valid docker-compose
//
// client select KV of docker service-domain (subdomain of ENV(HOMEDY_CF_TUNNEL_DOMAIN)) -> server validate and transform, then response with new deploy log
//
//	validation:
//	- has valid exposed ports in service
//		- if more than 1, response with select one of them
//		- if doesn't have one, reposne with create one port (internal port)
//	transform
//	- transform selected ports expose with validated random port
//		validation:
//		- check db is port already used
//		- check system is port already used
//
// client open log with given session (ws) -> server open log and stream it (TODO: idk how yet)
//
// after deployment:
// TODO: cloudflared
package services

import (
	"context"
	"homedy/internal/libs/replylib"
	"homedy/internal/middlewares"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Deploy struct {
	oAuthRepo         *repos.OAuth
	deployRepoRepo    *repos.DeployRepo
	deployLogRepo     *repos.DeployLog
	deploySessionRepo *repos.DeploySession
}

type ContextedDeploy struct {
	Deploy
	c   *gin.Context
	ctx context.Context
}

func NewDeploy(oAuthRepo *repos.OAuth, deployRepoRepo *repos.DeployRepo, deployLogRepo *repos.DeployLog, deploySessionRepo *repos.DeploySession) *Deploy {
	return &Deploy{oAuthRepo, deployRepoRepo, deployLogRepo, deploySessionRepo}
}

func (s *Deploy) AttachContext(c *gin.Context) *ContextedDeploy {
	return &ContextedDeploy{*s, c, c.Request.Context()}
}

func (s *ContextedDeploy) GetRepos(payload payloads.TemplateWithSession) ([]models.FilteredGHRepo, error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	return s.deploySessionRepo.GetSessionOrFetch(s.ctx, payload.Session, client)
}

func (s *ContextedDeploy) SelectRepo(payload payloads.RequestSelectRepo) (*models.FilteredGHRepo, error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	repos, err := s.deploySessionRepo.GetSessionOrFetch(s.ctx, payload.Session, client)
	if err != nil {
		return nil, err
	}

	var selectedRepo *models.FilteredGHRepo
	for _, repo := range repos {
		if repo.ID == payload.ID {
			selectedRepo = &repo
		}
	}
	if selectedRepo == nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: "repository not found",
		}
	}

	return selectedRepo, s.deploySessionRepo.SetSelectedRepo(s.ctx, payload.Session, *selectedRepo)
}

func (s *ContextedDeploy) GetSelectedRepo(payload payloads.TemplateWithSession) (*models.FilteredGHRepo, error) {
	return s.deploySessionRepo.GetSelectedRepo(s.ctx, payload.Session)
}
