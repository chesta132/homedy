package models

import (
	"homedy/internal/libs/authlib"

	"gorm.io/gorm"
)

type UserStatus string

const (
	UserActive  UserStatus = "active"
	UserPending UserStatus = "pending"
)

type User struct {
	BaseRecyclable
	Username string     `json:"username" gorm:"not null,unique" example:"chesta_ardiona"`
	Email    string     `json:"email" gorm:"not null,unique" example:"chestaardi4@gmail.com"`
	Password string     `json:"-" gorm:"not null"`
	Status   UserStatus `json:"status" gorm:"not null;default:'pending';index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if !authlib.IsHashed(u.Password) {
		hashed, err := authlib.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}
