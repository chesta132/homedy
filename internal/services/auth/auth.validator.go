package auth

import (
	"context"
	"errors"
	"homedy/internal/libs/replylib"
	"homedy/internal/models"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
)

func validateEmailAndUsername(db *gorm.DB, ctx context.Context, email, username string) error {
	user, err := gorm.G[models.User](db).
		Select("Email", "Username").
		Where("email = ? OR username = ?", email, username).
		First(ctx)
	isErrNotFound := errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil && !isErrNotFound {
		return err
	}

	if !isErrNotFound {
		fe := make(reply.FieldsError)
		if user.Email == email {
			fe["email"] = "email already registered"
		}
		if user.Username == username {
			fe["username"] = "username already registered"
		}
		return &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "email or username already registered",
			Fields:  fe,
		}
	}
	return nil
}
