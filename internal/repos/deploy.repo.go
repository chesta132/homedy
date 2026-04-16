package repos

import (
	"homedy/internal/models"

	"gorm.io/gorm"
)

// deploy repo

type DeployRepo struct {
	db *gorm.DB
	create[models.DeployRepo]
	read[models.DeployRepo]
	update[models.DeployRepo]
	delete[models.DeployRepo]
}

func NewDeployRepo(db *gorm.DB) *DeployRepo {
	return &DeployRepo{db, create[models.DeployRepo]{db}, read[models.DeployRepo]{db}, update[models.DeployRepo]{db}, delete[models.DeployRepo]{db}}
}

func (r *DeployRepo) DB() *gorm.DB {
	return r.db
}

func (r *DeployRepo) WithContext(tx *gorm.DB) *DeployRepo {
	return NewDeployRepo(tx)
}

// deploy log

type DeployLog struct {
	db *gorm.DB
	create[models.DeployLog]
	read[models.DeployLog]
	update[models.DeployLog]
	delete[models.DeployLog]
}

func NewDeployLog(db *gorm.DB) *DeployLog {
	return &DeployLog{db, create[models.DeployLog]{db}, read[models.DeployLog]{db}, update[models.DeployLog]{db}, delete[models.DeployLog]{db}}
}

func (r *DeployLog) DB() *gorm.DB {
	return r.db
}

func (r *DeployLog) WithContext(tx *gorm.DB) *DeployLog {
	return NewDeployLog(tx)
}
