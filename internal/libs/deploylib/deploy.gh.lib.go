package deploylib

import (
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"

	"github.com/google/go-github/v68/github"
)

func FilterGHRepo(repo *github.Repository) models.FilteredGHRepo {
	return models.FilteredGHRepo{ID: *repo.ID, Name: *repo.Name, FullName: *repo.FullName}
}

func FilterGHRepos(repos []*github.Repository) []models.FilteredGHRepo {
	return slicelib.Map(repos, func(i int, r *github.Repository) models.FilteredGHRepo { return FilterGHRepo(r) })
}
