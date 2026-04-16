package payloads

import (
	"homedy/internal/libs/slicelib"

	"github.com/google/go-github/v68/github"
)

type TemplateWithSession struct {
	Session string `uri:"session" json:"session" validate:"required,uuid4"`
}

type ResponseGetRepo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
}

func ToResponseGetRepo(repo *github.Repository) ResponseGetRepo {
	return ResponseGetRepo{ID: *repo.ID, Name: *repo.Name, FullName: *repo.FullName}
}

func ToResponseGetRepos(repos []*github.Repository) []ResponseGetRepo {
	return slicelib.Map(repos, func(i int, r *github.Repository) ResponseGetRepo { return ToResponseGetRepo(r) })
}
