package repos

import (
	"context"
	"homedy/internal/libs/deploylib"
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"
	"time"

	"github.com/chesta132/goreply/reply"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/google/go-github/v68/github"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DeploySession struct {
	rdb *redis.Client
	redisRepo
}

func NewDeploySession(rdb *redis.Client) *DeploySession {
	return &DeploySession{rdb, redisRepo{rdb}}
}

func deploySessionKey(session string) string {
	return "deploy:session:" + session
}

// session --------------

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

func (r *DeploySession) RemoveSession(ctx context.Context, session string) error {
	key := deploySessionKey(session)
	return r.rdb.Del(ctx, key).Err()
}

// repos --------------

func (r *DeploySession) GetRepos(ctx context.Context, session string) (repos []models.FilteredGHRepo, err error) {
	err = r.hGetWithParse(ctx, deploySessionKey(session), "repos", &repos)
	return
}

func (r *DeploySession) SetRepos(ctx context.Context, session string, repos []models.FilteredGHRepo) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"repos": repos})
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

func (r *DeploySession) GetSelectedRepo(ctx context.Context, session string) (repo *models.SelectedRepoInSession, err error) {
	err = r.hGetWithParse(ctx, deploySessionKey(session), "repo", &repo)
	return
}

func (r *DeploySession) SetSelectedRepo(ctx context.Context, session string, repo models.SelectedRepoInSession) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"repo": repo})
}

func (r *DeploySession) LazySetSelectedRepo(ctx context.Context, session string, repo models.SelectedRepoInSession) error {
	oldRepo, err := r.GetSelectedRepo(ctx, session)
	if err != nil && err != redis.Nil {
		return err
	}
	if (err == nil && oldRepo.ID != repo.ID) || (err == redis.Nil) {
		if err = r.SetSelectedRepo(ctx, session, repo); err != nil {
			return err
		}
	}
	return nil
}

func (r *DeploySession) GetRepoAndBranchFromRepos(ctx context.Context, session string, ghClient *github.Client, repos []models.FilteredGHRepo, repoID int64, branchName string) (*models.FilteredGHRepo, *models.FilteredGHRepoBranch, error) {
	for _, repo := range repos {
		if repo.ID == repoID {
			branches, err := r.GetBranchesOrFetch(ctx, session, &repo, ghClient)
			if err != nil {
				return nil, nil, err
			}
			for _, branch := range branches {
				if branch.Name == branchName {
					return &repo, &branch, nil

				}
			}
			return nil, nil, &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "branch not found"}
		}
	}
	return nil, nil, &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "repository not found"}
}

// branches --------------

func (r *DeploySession) GetBranches(ctx context.Context, session string, repoID int64) (branches []models.FilteredGHRepoBranch, err error) {
	err = r.hGetWithParse(ctx, deploySessionKey(session), "branches", &branches)
	if err != nil {
		return nil, err
	}

	branches = slicelib.Filter(branches, func(idx int, val models.FilteredGHRepoBranch) bool { return val.RepoID == repoID })
	// branches must be at least one
	if len(branches) == 0 {
		return nil, redis.Nil
	}
	return
}

func (r *DeploySession) SetBranches(ctx context.Context, session string, branches []models.FilteredGHRepoBranch) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"branches": branches})
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

// compose --------------

func (r *DeploySession) GetComposes(ctx context.Context, session string) (composes []models.DeploySessionCompose, err error) {
	err = r.hGetWithParse(ctx, deploySessionKey(session), "composes", &composes)
	return
}

func (r *DeploySession) SetComposes(ctx context.Context, session string, composes []models.DeploySessionCompose) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"composes": composes})
}

// this func also cache to session
func (r *DeploySession) SearchComposeProjectOfRepo(ctx context.Context, session string, ghClient *github.Client, ghUsername string, repo models.FilteredGHRepo) (*types.Project, error) {
	// check compose
	composes, err := r.GetComposes(ctx, session)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// is compose cached
	for _, compose := range composes {
		if compose.RepoID == repo.ID {
			// set services and validate docker compose
			project, err := deploylib.LoadDockerCompose(ctx, session, compose.Content)
			if err != nil {
				return nil, err
			}
			return project, nil
		}
	}

	// get from github and append/create compose to cache if compose not cached
	content, err := deploylib.GetDockerCompose(ctx, ghClient, ghUsername, repo.Name)
	if err != nil {
		return nil, err
	}

	// set services and validate docker compose
	project, err := deploylib.LoadDockerCompose(ctx, session, content)
	if err != nil {
		return nil, err
	}

	composes = append(composes, models.DeploySessionCompose{RepoID: repo.ID, Content: content})
	err = r.SetComposes(ctx, session, composes)
	if err != nil {
		return nil, err
	}

	return project, nil
}
