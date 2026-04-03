package payloads

import "homedy/internal/models"

type RequestCreateNote struct {
	Title      string                `json:"title" validate:"required"`
	Content    string                `json:"content" validate:"required"`
	Visibility models.NoteVisibility `json:"visibility" validate:"note_visibility"`
}

func (p *RequestCreateNote) ToNote(userID string) *models.Note {
	return &models.Note{Title: p.Title, Content: p.Content, Visibility: p.Visibility, UserID: userID}
}

type RequestUpdateNote struct {
	ID         string                `uri:"id" validate:"required,uuid4" swaggerignore:"true"`
	Title      string                `json:"title"`
	Content    string                `json:"content"`
	Visibility models.NoteVisibility `json:"visibility" validate:"omitempty,note_visibility"`
}

// ToNote returns *[models.Note] without ID
func (p *RequestUpdateNote) ToNote() *models.Note {
	return &models.Note{Title: p.Title, Content: p.Content, Visibility: p.Visibility}
}

type RequestGetOneNote struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestGetManyNote struct {
	Offset   int  `form:"offset" validate:"omitempty,min=0"`
	Recycled bool `form:"recycled"`
}

type RequestDeleteOneNote struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestDeleteManyNote struct {
	IDs []string `json:"ids" validate:"required,dive,uuid4"`
}

type RequestRestoreOneNote struct {
	ID string `uri:"id" validate:"required,uuid4"`
}

type RequestRestoreManyNote struct {
	IDs []string `json:"ids" validate:"required,dive,uuid4"`
}
