package models

import (
	"homedy/config"
	"homedy/internal/libs/cryptolib"

	"gorm.io/gorm"
)

// gorm not reflecting first before fire hooks, make sure to manual encrypt and decrypt on update
type OAuth struct {
	Base
	AppID        int64  `gorm:"not null" json:"-"`
	Username     string `gorm:"not null" json:"-"`
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`

	UserID string `gorm:"not null"`
	User   *User  `gorm:"constraint:OnDelete:CASCADE"`
}

func cryptoOAuth(oauth *OAuth, f func([]byte, []byte) (string, error)) error {
	if oauth.AccessToken != "" {
		title, err := f([]byte(oauth.AccessToken), []byte(config.GITHUB_OAUTH_CRYPTO_KEY))
		if err != nil {
			return err
		}
		oauth.AccessToken = title
	}

	if oauth.RefreshToken != "" {
		content, err := f([]byte(oauth.RefreshToken), []byte(config.GITHUB_OAUTH_CRYPTO_KEY))
		if err != nil {
			return err
		}
		oauth.RefreshToken = content
	}

	return nil
}

func (n *OAuth) Encrypt() error {
	return cryptoOAuth(n, cryptolib.EncryptGCM)
}

func (n *OAuth) Decrypt() error {
	return cryptoOAuth(n, cryptolib.DecryptGCM)
}

// hooks

func (n *OAuth) BeforeCreate(tx *gorm.DB) error {
	return n.Encrypt()
}

func (n *OAuth) AfterCreate(tx *gorm.DB) error {
	return n.Decrypt()
}

func (n *OAuth) AfterFind(tx *gorm.DB) error {
	return n.Decrypt()
}
