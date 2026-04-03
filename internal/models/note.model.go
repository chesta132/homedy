package models

import (
	"homedy/config"
	"homedy/internal/libs/cryptolib"

	"gorm.io/gorm"
)

type NoteVisibility string

const (
	NotePublic  NoteVisibility = "public"
	NotePrivate NoteVisibility = "private"
)

var NoteVisibilities = []NoteVisibility{NotePublic, NotePrivate}

type Note struct {
	BaseRecyclable
	Title      string         `json:"title" gorm:"not null"`
	Content    string         `json:"content" gorm:"not null"`
	Visibility NoteVisibility `json:"visibility" gorm:"not null,default:'private'"`

	UserID string `json:"user_id"`
	User   User   `json:"user,omitzero"`
}

func cryptoNote(note *Note, f func([]byte, []byte) (string, error)) error {
	if note.Title != "" {
		title, err := f([]byte(note.Title), []byte(config.NOTE_CRYPTO_KEY))
		if err != nil {
			return err
		}
		note.Title = title
	}

	if note.Content != "" {
		content, err := f([]byte(note.Content), []byte(config.NOTE_CRYPTO_KEY))
		if err != nil {
			return err
		}
		note.Content = content
	}

	return nil
}

func (n *Note) Encrypt() error {
	return cryptoNote(n, cryptolib.EncryptGCM)
}

func (n *Note) Decrypt() error {
	return cryptoNote(n, cryptolib.DecryptGCM)
}

// hooks

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	return n.Encrypt()
}

func (n *Note) AfterCreate(tx *gorm.DB) error {
	return n.Decrypt()
}

// gorm not reflecting first before fire hooks, make sure to manual encrypt
// func (n *Note) BeforeUpdate(tx *gorm.DB) error {
// 	logger.Debug("check before update note", logger.Fields("title", n.Title, "content", n.Content))
// 	return n.Encrypt()
// }

// gorm not reflecting first before fire hooks, make sure to manual decrypt
// func (n *Note) AfterUpdate(tx *gorm.DB) error {
// 	logger.Debug("check after update note", logger.Fields("title", n.Title, "content", n.Content))
// 	return n.Decrypt()
// }

func (n *Note) AfterFind(tx *gorm.DB) error {
	return n.Decrypt()
}
