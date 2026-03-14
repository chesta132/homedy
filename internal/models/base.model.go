package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string    `json:"id" gorm:"default:gen_random_uuid()" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
	CreatedAt time.Time `json:"created_at,omitzero" gorm:"autoCreateTime;not null" example:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt time.Time `json:"updated_at,omitzero" gorm:"autoUpdateTime;not null" example:"2006-01-02T15:04:05Z07:00"`
}

type BaseRecyclable struct {
	Base
	DeleteAt gorm.DeletedAt `json:"delete_at" gorm:"index" example:"2006-01-02T15:04:05Z07:00"`
}
