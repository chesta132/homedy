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
	err = r.hGetWithParse(ctx, deploySessionKey(session), "selectedRepo", &repo)
	return
}

func (r *DeploySession) SetSelectedRepo(ctx context.Context, session string, repo models.SelectedRepoInSession) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"selectedRepo": repo})
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

func (r *DeploySession) GetRepoAndBranchFromRepos(ctx context.Context, session string, ghClient *github.Client, repos []models.FilteredGHRepo, repoID int64) (*models.FilteredGHRepo, []string, error) {
	for _, repo := range repos {
		if repo.ID == repoID {
			branches, err := r.GetBranchesOrFetch(ctx, session, &repo, ghClient)
			if err != nil {
				return nil, nil, err
			}
			return &repo, branches, nil
		}
	}
	return nil, nil, &reply.ErrorPayload{Code: replylib.CodeNotFound, Message: "repository not found"}
}

// branches --------------

func (r *DeploySession) GetBranches(ctx context.Context, session string, repoID int64) ([]string, error) {
	repoBranches := make(models.DeploySessionRepoBranches)
	err := r.hGetWithParse(ctx, deploySessionKey(session), "branches", &repoBranches)
	if err != nil {
		return nil, err
	}
	branches, ok := repoBranches[repoID]
	if !ok {
		return nil, redis.Nil
	}
	return branches, nil
}

func (r *DeploySession) GetAllBranches(ctx context.Context, session string) (models.DeploySessionRepoBranches, error) {
	repoBranches := make(models.DeploySessionRepoBranches)
	err := r.hGetWithParse(ctx, deploySessionKey(session), "branches", &repoBranches)
	return repoBranches, err
}

func (r *DeploySession) SetBranches(ctx context.Context, session string, repoBranches models.DeploySessionRepoBranches) error {
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"branches": repoBranches})
}

func (r *DeploySession) GetBranchesOrFetch(ctx context.Context, session string, repo *models.FilteredGHRepo, ghClient *github.Client) (branches []string, err error) {
	repoBranches, err := r.GetAllBranches(ctx, session)
	if err == nil {
		if branches, ok := repoBranches[repo.ID]; ok {
			return branches, nil
		}
	}

	ghUsername := deploylib.GetGHUsernameFromRepo(*repo)
	ghBranch, _, err := ghClient.Repositories.ListBranches(ctx, ghUsername, repo.Name, nil)
	if err != nil {
		return nil, err
	}
	branches = slicelib.Map(ghBranch, func(i int, b *github.Branch) string { return b.GetName() })
	repoBranches[repo.ID] = branches

	err = r.SetBranches(ctx, session, repoBranches)
	if err != nil {
		return nil, err
	}

	return branches, nil
}

// compose --------------

func (r *DeploySession) GetComposes(ctx context.Context, session string) (models.DeploySessionCompose, error) {
	composes := make(models.DeploySessionCompose)
	err := r.hGetWithParse(ctx, deploySessionKey(session), "composes", &composes)
	return composes, err
}

func (r *DeploySession) SetComposes(ctx context.Context, session string, composes models.DeploySessionCompose) error {
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
	if compose, ok := composes[repo.ID]; ok {
		// set services and validate docker compose
		project, err := deploylib.LoadDockerCompose(ctx, session, compose)
		if err != nil {
			return nil, err
		}
		return project, nil
	}

	// get from github and append/create compose to cache if compose not cached
	compose, err := deploylib.GetDockerCompose(ctx, ghClient, ghUsername, repo.Name)
	if err != nil {
		return nil, err
	}

	// set services and validate docker compose
	project, err := deploylib.LoadDockerCompose(ctx, session, compose)
	if err != nil {
		return nil, err
	}

	composes[repo.ID] = compose
	err = r.SetComposes(ctx, session, composes)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// env --------------

func (r *DeploySession) GetEnv(ctx context.Context, session string) (env models.DeploySessionEnv, err error) {
	err = r.hGetWithParse(ctx, deploySessionKey(session), "env", &env)
	if err == nil {
		var decryptedEnv *models.DeploySessionEnv
		decryptedEnv, err = deploylib.DecryptSessionEnv(env)
		if err != nil {
			return
		}
		env = *decryptedEnv
	}
	return
}

func (r *DeploySession) SetEnv(ctx context.Context, session string, env models.DeploySessionEnv) error {
	encryptedEnv, err := deploylib.EncryptSessionEnv(env)
	if err != nil {
		return err
	}
	return r.hSetWithParse(ctx, deploySessionKey(session), map[string]any{"env": encryptedEnv})
}
