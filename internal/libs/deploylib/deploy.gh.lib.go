package deploylib

import (
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"

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
