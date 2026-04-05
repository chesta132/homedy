package services

import (
	"context"
	"fmt"
	"homedy/internal/libs/logger"
	"homedy/internal/libs/replylib"
	"homedy/internal/middlewares"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"
	"slices"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Note struct {
	noteRepo *repos.Note
}

type ContextedNote struct {
	Note
	c   *gin.Context
	ctx context.Context
}

func NewNote(noteRepo *repos.Note) *Note {
	return &Note{noteRepo}
}

func (s *Note) AttachContext(c *gin.Context) *ContextedNote {
	return &ContextedNote{*s, c, c.Request.Context()}
}

func (s *Note) assertOwner(noteUserIDs []string, userID, action string) error {
	if slices.ContainsFunc(noteUserIDs, func(id string) bool { return id != userID }) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: fmt.Sprintf("you can not %s", action),
		}
	}
	return nil
}

func (s *ContextedNote) CreateOne(payload payloads.RequestCreateNote) (*models.Note, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	note := payload.ToNote(userID)
	err = s.noteRepo.Create(s.ctx, note)
	return note, err
}

func (s *ContextedNote) GetOne(payload payloads.RequestGetOneNote) (*models.Note, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	note, err := s.noteRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if note.Visibility != models.NotePublic {
		err = s.assertOwner([]string{note.UserID}, userID, "read this note")
		if err != nil {
			return nil, err
		}
	}

	return &note, nil
}

func (s *ContextedNote) GetNotes(payload payloads.RequestGetManyNote) ([]models.Note, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	if payload.Sort == "" {
		payload.Sort = models.DESC
	}

	return s.noteRepo.GetNotesWithPayload(s.ctx, userID, payload)
}

func (s *ContextedNote) UpdateOne(payload payloads.RequestUpdateNote) (note *models.Note, err error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	s.noteRepo.DB().Transaction(func(tx *gorm.DB) error {
		noteRepo := s.noteRepo.WithContext(tx)

		note = payload.ToNote()
		if err = note.Encrypt(); err != nil {
			return err
		}

		*note, err = noteRepo.UpdateAndGet(s.ctx, *note, "id = ?", payload.ID)
		if err != nil {
			return err
		}

		if err = note.Decrypt(); err != nil {
			return err
		}

		err = s.assertOwner([]string{note.UserID}, userID, "update this note")
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (s *ContextedNote) DeleteOne(payload payloads.RequestDeleteOneNote) error {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return err
	}

	noteUserID, err := s.noteRepo.GetUserIDByID(s.ctx, payload.ID)
	if err != nil {
		return err
	}

	err = s.assertOwner([]string{noteUserID}, userID, "delete this note")
	if err != nil {
		return err
	}

	return s.noteRepo.Archive(s.ctx, "id = ?", payload.ID)
}

func (s *ContextedNote) DeleteMany(payload payloads.RequestDeleteManyNote) (err error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return err
	}

	noteUserIDs, err := s.noteRepo.GetUserIDsByIDs(s.ctx, payload.IDs)
	if err != nil {
		return
	}

	err = s.assertOwner(noteUserIDs, userID, "delete these notes")
	if err != nil {
		return err
	}

	return s.noteRepo.Archive(s.ctx, "id IN ?", payload.IDs)
}

func (s *ContextedNote) RestoreOne(payload payloads.RequestRestoreOneNote) (*models.Note, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	noteUserID, err := s.noteRepo.GetUserIDByRecycledID(s.ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	err = s.assertOwner([]string{noteUserID}, userID, "restore this note")
	if err != nil {
		return nil, err
	}

	notes, err := s.noteRepo.RestoreAndGet(s.ctx, "id = ?", payload.ID)
	if err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		notes[0].Decrypt()
		return &notes[0], nil
	}

	logger.Error("restoring note and success with no error but no notes returning", logger.Fields("note id", payload.ID))
	return nil, gorm.ErrRecordNotFound
}

func (s *ContextedNote) RestoreMany(payload payloads.RequestRestoreManyNote) ([]models.Note, error) {
	userID, err := middlewares.GetUserID(s.c)
	if err != nil {
		return nil, err
	}

	noteUserIDs, err := s.noteRepo.GetUserIDsByRecycledIDs(s.ctx, payload.IDs)
	if err != nil {
		return nil, err
	}

	err = s.assertOwner(noteUserIDs, userID, "restore these notes")
	if err != nil {
		return nil, err
	}

	if payload.Sort == "" {
		payload.Sort = models.DESC
	}

	return s.noteRepo.RestoreNotesWithPayload(s.ctx, payload)
}
