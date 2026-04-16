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
	"encoding/json"
	"homedy/internal/middlewares"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Deploy struct {
	rdb            *redis.Client
	oAuthRepo      *repos.OAuth
	deployRepoRepo *repos.DeployRepo
	deployLogRepo  *repos.DeployLog
}

type ContextedDeploy struct {
	Deploy
	c   *gin.Context
	ctx context.Context
}

func NewDeploy(rdb *redis.Client, oAuthRepo *repos.OAuth, deployRepoRepo *repos.DeployRepo, deployLogRepo *repos.DeployLog) *Deploy {
	return &Deploy{rdb, oAuthRepo, deployRepoRepo, deployLogRepo}
}

func (s *Deploy) AttachContext(c *gin.Context) *ContextedDeploy {
	return &ContextedDeploy{*s, c, c.Request.Context()}
}

func (s *ContextedDeploy) GetRepos(payload payloads.TemplateWithSession) (repos []payloads.ResponseGetRepo, err error) {
	client, err := middlewares.GetGithubClient(s.c)
	if err != nil {
		return nil, err
	}

	var reposStr string
	err = s.rdb.HGet(s.ctx, "deploy:session:"+payload.Session, "repos").Scan(&reposStr)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err != redis.Nil {
		if err = json.Unmarshal([]byte(reposStr), &repos); err == nil {
			return
		}
	}

	ghRepos, _, err := client.Repositories.ListByAuthenticatedUser(s.ctx, nil)
	if err != nil {
		return nil, err
	}
	repos = payloads.ToResponseGetRepos(ghRepos)

	reposBytes, err := json.Marshal(repos)
	if err == nil {
		s.rdb.HSet(s.ctx, "deploy:session:"+payload.Session, "repos", string(reposBytes))
	}

	return repos, nil
}
