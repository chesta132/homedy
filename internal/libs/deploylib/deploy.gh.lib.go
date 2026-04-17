package deploylib

import (
	"context"
	"errors"
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"
	"strings"

	"github.com/google/go-github/v68/github"
)

func FilterGHRepo(repo *github.Repository) models.FilteredGHRepo {
	return models.FilteredGHRepo{ID: repo.GetID(), Name: repo.GetName(), FullName: repo.GetFullName()}
}

func FilterGHRepos(repos []*github.Repository) []models.FilteredGHRepo {
	return slicelib.Map(repos, func(i int, r *github.Repository) models.FilteredGHRepo { return FilterGHRepo(r) })
}

func FilterGHBranch(branch *github.Branch, repoID int64) models.FilteredGHRepoBranch {
	return models.FilteredGHRepoBranch{Name: branch.GetName(), RepoID: repoID}
}

func FilterGHBranches(branches []*github.Branch, repoID int64) []models.FilteredGHRepoBranch {
	return slicelib.Map(branches, func(i int, r *github.Branch) models.FilteredGHRepoBranch { return FilterGHBranch(r, repoID) })
}

func GetGHUsernameFromRepo(repo models.FilteredGHRepo) string {
	ghUsername, _, _ := strings.Cut(repo.FullName, "/")
	return ghUsername
}

func GetGHContent(ctx context.Context, ghClient *github.Client, ghUsername, repoName, path string) (string, error) {
	file, _, _, err := ghClient.Repositories.GetContents(ctx, ghUsername, repoName, path, nil)
	if err != nil {
		return "", err
	}
	return file.GetContent()
}

func GetDockerCompose(ctx context.Context, ghClient *github.Client, ghUsername, repoName string) (string, error) {
	// following docker priority
	composePriority := []string{"compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml"}
	_, dir, _, err := ghClient.Repositories.GetContents(ctx, ghUsername, repoName, "", nil)
	if err != nil {
		return "", err
	}

	filesInRoot := make(map[string]string)
	for _, f := range dir {
		filesInRoot[f.GetName()] = f.GetPath()
	}

	for _, name := range composePriority {
		if path, ok := filesInRoot[name]; ok {
			return GetGHContent(ctx, ghClient, ghUsername, repoName, path)
		}
	}
	return "", errors.New("no docker compose")
}
