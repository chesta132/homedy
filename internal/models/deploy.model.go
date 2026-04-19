package models

import "github.com/compose-spec/compose-go/v2/types"

type DeployRepo struct {
	Base
	RepoOwner string `gorm:"not null" json:"repoOwner"`
	RepoName  string `gorm:"not null" json:"repoName"`
	WebhookID int64  `gorm:"not null" json:"-"`
	Branch    string `json:"branch"`

	UserID string `gorm:"not null;index" json:"userId"`
	User   *User  `json:"user,omitempty"`
}

type DeployStatus string

const (
	DeploySuccess    = "success"
	DeployFailed     = "failed"
	DeployProcessing = "processing"
)

type DeployLog struct {
	Base
	Status    DeployStatus `gorm:"not null"`
	Session   string       `gorm:"not nulls"`
	CommitSHA string       `gorm:"not null"`
	CommitMsg string

	DeployRepoID string      `gorm:"not null;index"`
	DeployRepo   *DeployRepo `gorm:"constraint:OnDelete:CASCADE"`
}

type FilteredGHRepo struct {
	ID       int64  `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
}

type DeploySessionRepoBranches map[int64][]string

type SelectedRepoInSession struct {
	FilteredGHRepo
	Branch   string   `json:"branch"`
	Services []string `json:"services"`
}

// for redis
// session on key
type DeploySession struct {
	UserID       string `redis:"userId"`
	Repos        string `redis:"repos"`
	SelectedRepo string `redis:"selectedRepo"`
	Composes     string `redis:"composes"`
}

type DeploySessionCompose map[int64]string

type GlobalEnv types.MappingWithEquals
type ServiceEnv map[string]types.MappingWithEquals
type RepoEnv map[int64]ServiceEnv

type DeploySessionEnv struct {
	Global GlobalEnv
	Repo   RepoEnv // repo -> service -> env
}
