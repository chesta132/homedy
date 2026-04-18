// flow:
// client bind github first
// client send create deployment session -> server create session with caching (redis/pg)
// client check available repositories -> server check it on cache first and fallback to fetch github API
// client select repository in client side
// client check available branch from selected repository -> server check it on cache first and fallback to fetch github API
// client select branch (send selected repo and branch) -> server check is valid and cache it and response with available docker services
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
	"homedy/internal/libs/deploylib"
	"homedy/internal/libs/replylib"
	"homedy/internal/middlewares"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/chesta132/goreply/reply"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Deploy struct {
	oAuthRepo         *repos.OAuth
	deployRepoRepo    *repos.DeployRepo
	deployLogRepo     *repos.DeployLog
	deploySessionRepo *repos.DeploySession
	composeService    api.Compose
}

type ContextedDeploy struct {
	Deploy
	c   *gin.Context
	ctx context.Context
}

func NewDeploy(oAuthRepo *repos.OAuth, deployRepoRepo *repos.DeployRepo, deployLogRepo *repos.DeployLog, deploySessionRepo *repos.DeploySession, composeService api.Compose) *Deploy {
	return &Deploy{oAuthRepo, deployRepoRepo, deployLogRepo, deploySessionRepo, composeService}
}

func (s *Deploy) AttachContext(c *gin.Context) *ContextedDeploy {
	return &ContextedDeploy{*s, c, c.Request.Context()}
}

// repos -------------

func (s *ContextedDeploy) GetRepos(payload payloads.TemplateWithSession) ([]models.FilteredGHRepo, error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	return s.deploySessionRepo.GetReposOrFetch(s.ctx, payload.Session, client)
}

// branch -------------

func (s *ContextedDeploy) GetBranches(payload payloads.RequestGetBranches) ([]models.FilteredGHRepoBranch, error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	repos, err := s.deploySessionRepo.GetReposOrFetch(s.ctx, payload.Session, client)
	if err != nil {
		return nil, err
	}

	var selectedRepo *models.FilteredGHRepo
	for _, repo := range repos {
		if repo.ID == payload.ID {
			selectedRepo = &repo
			break
		}
	}
	if selectedRepo == nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: "repository not found",
		}
	}

	branches, err := s.deploySessionRepo.GetBranchesOrFetch(s.ctx, payload.Session, selectedRepo, client)
	if err != nil {
		return nil, err
	}

	return branches, err
}

// selected repo -------------

func (s *ContextedDeploy) SetSelectedRepo(payload payloads.RequestSetSelectedRepo) (*models.SelectedRepoInSession, error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	// get repos
	repos, err := s.deploySessionRepo.GetReposOrFetch(s.ctx, payload.Session, client)
	if err != nil {
		return nil, err
	}

	// select repo
	var selectedRepo *models.SelectedRepoInSession
	for _, repo := range repos {
		if repo.ID == payload.ID {
			branches, err := s.deploySessionRepo.GetBranchesOrFetch(s.ctx, payload.Session, &repo, client)
			if err != nil {
				return nil, err
			}
			for _, branch := range branches {
				if branch.Name == payload.Branch {
					selectedRepo = &models.SelectedRepoInSession{FilteredGHRepo: repo, Branch: branch}
					break
				}
			}
		}
	}
	if selectedRepo == nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: "repository or branch not found",
		}
	}
	ghUsername := deploylib.GetGHUsernameFromRepo(selectedRepo.FilteredGHRepo)

	// check compose
	composes, err := s.deploySessionRepo.GetComposes(s.ctx, payload.Session)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil || err == nil {
		var skipAdd bool
		// is compose cached
		for _, compose := range composes {
			if compose.RepoID == selectedRepo.ID {
				skipAdd = true
				break
			}
		}

		// compose add to cache if compose not cached
		if !skipAdd {
			content, err := deploylib.GetDockerCompose(s.ctx, client, ghUsername, selectedRepo.Name)
			if err != nil {
				return nil, err
			}

			// validate docker compose
			_, err = deploylib.LoadDockerCompose(s.ctx, s.composeService, payload.Session, content)
			if err != nil {
				return nil, err
			}

			composes = append(composes, models.DeploySessionCompose{RepoID: selectedRepo.ID, Content: content})
			err = s.deploySessionRepo.SetComposes(s.ctx, payload.Session, composes)
			if err != nil {
				return nil, err
			}
		}
	}

	// set selected repo
	return selectedRepo, s.deploySessionRepo.SetSelectedRepo(s.ctx, payload.Session, *selectedRepo)
}

func (s *ContextedDeploy) GetSelectedRepo(payload payloads.TemplateWithSession) (*models.SelectedRepoInSession, error) {
	return s.deploySessionRepo.GetSelectedRepo(s.ctx, payload.Session)
}
