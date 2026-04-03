package services

import (
	"context"
	"errors"
	"homedy/config"
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

// TODO: refactor (DRY)
// TODO: add crypto for content and title column in note model

func NewNote(noteRepo *repos.Note) *Note {
	return &Note{noteRepo}
}

func (s *Note) AttachContext(c *gin.Context) *ContextedNote {
	return &ContextedNote{*s, c, c.Request.Context()}
}

func (s *ContextedNote) CreateOne(payload payloads.RequestCreateNote) (*models.Note, error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	note := payload.ToNote(userID)
	err := s.noteRepo.Create(s.ctx, note)
	return note, err
}

func (s *ContextedNote) GetOne(payload payloads.RequestGetOneNote) (*models.Note, error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	note, err := s.noteRepo.GetByID(s.ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if note.UserID != userID && note.Visibility != models.NotePublic {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "you can not read this note",
		}
	}

	return &note, nil
}

func (s *ContextedNote) GetNotes(payload payloads.RequestGetManyNote) (notes []models.Note, err error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	// config.LIMIT_RESOURCE_PER_PAGINATION + 1 to cursor pagination
	query := s.noteRepo.DB().Where("user_id = ?", userID).Offset(payload.Offset).Limit(config.LIMIT_RESOURCE_PER_PAGINATION + 1)
	if payload.Recycled {
		query = query.Unscoped().Where("deleted_at != NULL")
	}
	err = query.Find(&notes).Error
	return
}

func (s *ContextedNote) UpdateOne(payload payloads.RequestUpdateNote) (note *models.Note, err error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	s.noteRepo.DB().Transaction(func(tx *gorm.DB) error {
		noteRepo := s.noteRepo.WithContext(tx)

		note = payload.ToNote()
		*note, err = noteRepo.UpdateAndGet(s.ctx, *note, "id = ?", payload.ID)
		if err != nil {
			return err
		}

		if note.UserID != userID {
			err = &reply.ErrorPayload{
				Code:    replylib.CodeForbidden,
				Message: "you can not update this note",
			}
			return err
		}

		return nil
	})
	return
}

func (s *ContextedNote) DeleteOne(payload payloads.RequestDeleteOneNote) error {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return errors.New("middleware skipped")
	}

	noteUserID, err := s.noteRepo.GetUserIDByID(s.ctx, payload.ID)
	if err != nil {
		return err
	}

	if noteUserID != userID {
		return &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "you can not delete this note",
		}
	}

	return s.noteRepo.Archive(s.ctx, "id = ?", payload.ID)
}

func (s *ContextedNote) DeleteMany(payload payloads.RequestDeleteManyNote) (err error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return errors.New("middleware skipped")
	}

	noteUserIDs, err := s.noteRepo.GetUserIDsByIDs(s.ctx, payload.IDs)
	if err != nil {
		return
	}

	if slices.ContainsFunc(noteUserIDs, func(id string) bool { return id != userID }) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "you can not delete these notes",
		}
	}

	return s.noteRepo.Archive(s.ctx, "id IN ?", payload.IDs)
}

func (s *ContextedNote) RestoreOne(payload payloads.RequestRestoreOneNote) (*models.Note, error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	noteUserID, err := s.noteRepo.GetUserIDByRecycledID(s.ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if noteUserID != userID {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "you can not restore this note",
		}
	}

	notes, err := s.noteRepo.RestoreAndGet(s.ctx, "id = ?", payload.ID)
	if err != nil {
		return nil, err
	}

	if len(notes) > 0 {
		return &notes[0], nil
	}

	logger.Error("restoring note and success with no error but no notes returning", logger.Fields("note id", payload.ID))
	return nil, gorm.ErrRecordNotFound
}

func (s *ContextedNote) RestoreMany(payload payloads.RequestRestoreManyNote) ([]models.Note, error) {
	userID, ok := middlewares.GetUserID(s.c)
	if !ok {
		return nil, errors.New("middleware skipped")
	}

	noteUserIDs, err := s.noteRepo.GetUserIDsByRecycledIDs(s.ctx, payload.IDs)
	if err != nil {
		return nil, err
	}

	if slices.ContainsFunc(noteUserIDs, func(id string) bool { return id != userID }) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "you can not restore these notes",
		}
	}

	return s.noteRepo.RestoreAndGet(s.ctx, "id IN ?", payload.IDs)
}
