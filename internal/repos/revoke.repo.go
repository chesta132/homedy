package repos

import (
	"context"
	"homedy/internal/libs/authlib"
	"homedy/internal/models"

	"gorm.io/gorm"
)

type Revoke struct {
	db *gorm.DB
	create[models.Revoke]
	read[models.Revoke]
	update[models.Revoke]
	delete[models.Revoke]
}

func NewRevoke(db *gorm.DB) *Revoke {
	return &Revoke{db, create[models.Revoke]{db}, read[models.Revoke]{db}, update[models.Revoke]{db}, delete[models.Revoke]{db}}
}

func (r *Revoke) DB() *gorm.DB {
	return r.db
}

func (r *Revoke) WithContext(tx *gorm.DB) *Revoke {
	return NewRevoke(tx)
}
func (r *Revoke) RevokeToken(ctx context.Context, token, reason string) error {
	claims, err := authlib.ParseRefreshToken(token)
	if err != nil {
		return err
	}
	revoke := models.Revoke{
		Value:       token,
		Reason:      reason,
		RevokeUntil: claims.ExpiresAt.Time,
	}
	return gorm.G[models.Revoke](r.db).Create(ctx, &revoke)
}
