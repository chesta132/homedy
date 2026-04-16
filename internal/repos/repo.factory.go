package repos

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repos struct {
	db            *gorm.DB
	rdb           *redis.Client
	user          *User
	revoke        *Revoke
	note          *Note
	oAuth         *OAuth
	deployRepo    *DeployRepo
	deployLog     *DeployLog
	deploySession *DeploySession
}

func New(db *gorm.DB, rdb *redis.Client) *Repos {
	return &Repos{db: db, rdb: rdb}
}

func (r *Repos) RDB() *redis.Client {
	return r.rdb
}

func (r *Repos) User() *User {
	if r.user == nil {
		r.user = NewUser(r.db)
	}
	return r.user
}

func (r *Repos) Revoke() *Revoke {
	if r.revoke == nil {
		r.revoke = NewRevoke(r.db)
	}
	return r.revoke
}

func (r *Repos) Note() *Note {
	if r.note == nil {
		r.note = NewNote(r.db)
	}
	return r.note
}

func (r *Repos) OAuth() *OAuth {
	if r.oAuth == nil {
		r.oAuth = NewOAuth(r.db)
	}
	return r.oAuth
}

func (r *Repos) DeployRepo() *DeployRepo {
	if r.deployRepo == nil {
		r.deployRepo = NewDeployRepo(r.db)
	}
	return r.deployRepo
}

func (r *Repos) DeployLog() *DeployLog {
	if r.deployLog == nil {
		r.deployLog = NewDeployLog(r.db)
	}
	return r.deployLog
}

func (r *Repos) DeploySession() *DeploySession {
	if r.deploySession == nil {
		r.deploySession = NewDeploySession(r.rdb)
	}
	return r.deploySession
}
