package repos

import (
	"context"
	"encoding/json"
	"homedy/internal/libs/deploylib"
	"homedy/internal/libs/slicelib"
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

// session --------------

// returns session id and error
func (r *DeploySession) CreateSession(ctx context.Context, userID string, ghUsername string) (string, error) {
	id := uuid.NewString()
	key := deploySessionKey(id)

	err := r.rdb.HSet(ctx, key, models.DeploySession{UserID: userID, GHUsername: ghUsername}).Err()
	if err != nil {
		return "", err
	}

	err = r.rdb.Expire(ctx, key, time.Hour).Err()
	if err != nil {
		return "", err
	}

	return id, err
}

func (r *DeploySession) RemoveSession(ctx context.Context, session string) error {
	key := deploySessionKey(session)
	return r.rdb.Del(ctx, key).Err()
}

func (r *DeploySession) GetGHUsername(ctx context.Context, session string) (string, error) {
	return r.rdb.HGet(ctx, deploySessionKey(session), "ghUsername").Result()
}

// repos --------------

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

func (r *DeploySession) GetReposOrFetch(ctx context.Context, session string, ghClient *github.Client) (repos []models.FilteredGHRepo, err error) {
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

// selected repo --------------

func (r *DeploySession) SetSelectedRepo(ctx context.Context, session string, repo models.SelectedRepoInSession) error {
	repoBytes, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	return r.rdb.HSet(ctx, deploySessionKey(session), "selectedRepo", string(repoBytes)).Err()
}

func (r *DeploySession) GetSelectedRepo(ctx context.Context, session string) (repo *models.SelectedRepoInSession, err error) {
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

// branches --------------

func (r *DeploySession) GetBranches(ctx context.Context, session string, repoID int64) (branches []models.FilteredGHRepoBranch, err error) {
	var branchesStr string
	err = r.rdb.HGet(ctx, deploySessionKey(session), "branches").Scan(&branchesStr)
	if err != nil {
		return nil, err
	}
	if branchesStr == "" {
		return nil, redis.Nil
	}

	err = json.Unmarshal([]byte(branchesStr), &branches)
	if err != nil {
		return
	}

	branches = slicelib.Filter(branches, func(idx int, val models.FilteredGHRepoBranch) bool { return val.RepoID == repoID })
	// branches must be at least one
	if len(branches) == 0 {
		return nil, redis.Nil
	}
	return
}

func (r *DeploySession) GetBranchesOrFetch(ctx context.Context, session string, repo *models.FilteredGHRepo, ghClient *github.Client) (branches []models.FilteredGHRepoBranch, err error) {
	branches, err = r.GetBranches(ctx, session, repo.ID)
	if err == nil {
		return
	}

	ghUsername := deploylib.GetGHUsernameFromRepo(*repo)
	ghBranch, _, err := ghClient.Repositories.ListBranches(ctx, ghUsername, repo.Name, nil)
	if err != nil {
		return nil, err
	}
	branches = deploylib.FilterGHBranches(ghBranch, repo.ID)

	err = r.SetBranches(ctx, session, branches)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *DeploySession) SetBranches(ctx context.Context, session string, branches []models.FilteredGHRepoBranch) error {
	branchesBytes, err := json.Marshal(branches)
	if err != nil {
		return err
	}
	return r.rdb.HSet(ctx, deploySessionKey(session), "branches", string(branchesBytes)).Err()
}

// compose --------------

func (r *DeploySession) GetCompose(ctx context.Context, session string) (composes []models.DeploySessionCompose, err error) {
	var composeStr string
	err = r.rdb.HGet(ctx, deploySessionKey(session), "composes").Scan(&composeStr)
	if err != nil {
		return nil, err
	}
	if composeStr == "" {
		return nil, redis.Nil
	}

	err = json.Unmarshal([]byte(composeStr), &composes)
	return
}

func (r *DeploySession) SetCompose(ctx context.Context, session string, compose []models.DeploySessionCompose) error {
	composeBytes, err := json.Marshal(compose)
	if err != nil {
		return err
	}
	return r.rdb.HSet(ctx, deploySessionKey(session), "composes", string(composeBytes)).Err()
}
