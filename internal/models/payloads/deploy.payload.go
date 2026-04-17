package payloads

import (
	"homedy/internal/libs/slicelib"
	"homedy/internal/models"

	"github.com/google/go-github/v68/github"
)

type TemplateWithSession struct {
	Session string `uri:"session" json:"session" validate:"required,uuid4"`
}

func FilterGHRepo(repo *github.Repository) models.FilteredGHRepo {
	return models.FilteredGHRepo{ID: *repo.ID, Name: *repo.Name, FullName: *repo.FullName}
}

func FilterGHRepos(repos []*github.Repository) []models.FilteredGHRepo {
	return slicelib.Map(repos, func(i int, r *github.Repository) models.FilteredGHRepo { return FilterGHRepo(r) })
}

type RequestSetSelectedRepo struct {
	TemplateWithSession
	ID     int64  `json:"id" validate:"required"` // repo id
	Branch string `json:"branch" validate:"required"`
}

type RequestGetBranches struct {
	TemplateWithSession
	ID int64 `uri:"id" validate:"required"`
}
