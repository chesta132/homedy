package repos

import (
	"context"
	"encoding/json"
	"homedy/internal/libs/deploylib"
	"homedy/internal/models"
	"time"

	"github.com/google/go-github/v68/github"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DeploySession struct {
	rdb *redis.Client
}

func NewDeploySession(rdb *redis.Client) *DeploySession {
	return &DeploySession{rdb}
}

func deploySessionKey(session string) string {
	return "deploy:session:" + session
}

// returns session id and error
func (r *DeploySession) CreateSession(ctx context.Context, userID string) (string, error) {
	id := uuid.NewString()
	key := deploySessionKey(id)

	err := r.rdb.HSet(ctx, key, models.DeploySession{UserID: userID}).Err()
	if err != nil {
		return "", err
	}

	err = r.rdb.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		return "", err
	}

	return id, err
}

func (r *DeploySession) GetRepos(ctx context.Context, session string) (repos []models.FilteredGHRepo, err error) {
	var reposStr string
	err = r.rdb.HGet(ctx, deploySessionKey(session), "repos").Scan(&reposStr)
	if err != nil {
		return nil, err
	}
	if reposStr == "" {
		return nil, redis.Nil
	}

	err = json.Unmarshal([]byte(reposStr), &repos)
	return
}

func (r *DeploySession) SetRepos(ctx context.Context, session string, repos []models.FilteredGHRepo) error {
	reposBytes, err := json.Marshal(repos)
	if err != nil {
		return err
	}
	return r.rdb.HSet(ctx, deploySessionKey(session), "repos", string(reposBytes)).Err()
}

func (r *DeploySession) GetSessionOrFetch(ctx context.Context, session string, ghClient *github.Client) (repos []models.FilteredGHRepo, err error) {
	repos, err = r.GetRepos(ctx, session)
	if err == nil {
		return
	}

	ghRepos, _, err := ghClient.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		return nil, err
	}
	repos = deploylib.FilterGHRepos(ghRepos)

	err = r.SetRepos(ctx, session, repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *DeploySession) SetSelectedRepo(ctx context.Context, session string, repo models.FilteredGHRepo) error {
	repoBytes, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	return r.rdb.HSet(ctx, deploySessionKey(session), "selectedRepo", string(repoBytes)).Err()
}

func (r *DeploySession) GetSelectedRepo(ctx context.Context, session string) (repo *models.FilteredGHRepo, err error) {
	var repoStr string
	err = r.rdb.HGet(ctx, deploySessionKey(session), "selectedRepo").Scan(&repoStr)
	if err != nil {
		return nil, err
	}
	if repoStr == "" {
		return nil, redis.Nil
	}

	err = json.Unmarshal([]byte(repoStr), &repo)
	return
}
