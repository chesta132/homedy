package models

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
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
}

// for redis
// session on key
type DeploySession struct {
	UserID string `redis:"userId"`
	Repos  string `redis:"repos"`
}
