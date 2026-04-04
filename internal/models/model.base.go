package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseID struct {
	ID string `json:"id" gorm:"default:gen_random_uuid()" example:"479b5b5f-81b1-4669-91a5-b5bf69e597c6"`
}

type Base struct {
	BaseID
	CreatedAt time.Time `json:"createdAt,omitzero" gorm:"autoCreateTime;not null" example:"2006-01-02T15:04:05Z07:00"`
	UpdatedAt time.Time `json:"updatedAt,omitzero" gorm:"autoUpdateTime;not null" example:"2006-01-02T15:04:05Z07:00"`
}

type BaseRecyclable struct {
	Base
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index" example:"2006-01-02T15:04:05Z07:00" swaggertype:"string" format:"date"`
}

type Sort string

const (
	ASC  Sort = "asc"
	DESC Sort = "desc"
)

var Sorts = []Sort{ASC, DESC}
