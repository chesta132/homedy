package models

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
