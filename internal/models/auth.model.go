package models

import "time"

type Revoke struct {
	Base        `json:"-"`
	Value       string    `gorm:"uniqueIndex;<-:create;not null" json:"-"`
	Reason      string    `gorm:"not null" json:"-"`
	RevokeUntil time.Time `gorm:"index;<-:create;not null" json:""`
}
