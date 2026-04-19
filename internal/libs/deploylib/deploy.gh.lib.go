package deploylib

import (
	"context"
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/google/go-github/v68/github"
)

func FilterGHRepo(repo *github.Repository) models.FilteredGHRepo {
	return models.FilteredGHRepo{ID: repo.GetID(), Name: repo.GetName(), FullName: repo.GetFullName()}
}

func FilterGHRepos(repos []*github.Repository) []models.FilteredGHRepo {
	return slicelib.Map(repos, func(i int, r *github.Repository) models.FilteredGHRepo { return FilterGHRepo(r) })
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
	composePath := append(cli.DefaultFileNames, cli.DefaultOverrideFileNames...)
	_, dir, _, err := ghClient.Repositories.GetContents(ctx, ghUsername, repoName, "", nil)
	if err != nil {
		return "", err
	}

	filesInRoot := make(map[string]string)
	for _, f := range dir {
		filesInRoot[f.GetName()] = f.GetPath()
	}

	for _, name := range composePath {
		if path, ok := filesInRoot[name]; ok {
			return GetGHContent(ctx, ghClient, ghUsername, repoName, path)
		}
	}
	return "", &reply.ErrorPayload{Code: replylib.CodeUnprocessableEntity, Message: "no docker compose"}
}
