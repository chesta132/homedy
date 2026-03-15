package services

import (
	"context"
	"fmt"
	"homedy/internal/libs/iolib"
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/sambalib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"os"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Samba struct{}

type ContextedSamba struct {
	Samba
	c   *gin.Context
	ctx context.Context
}

func NewSamba() *Samba {
	return &Samba{}
}

func (s *Samba) AttachContext(c *gin.Context) *ContextedSamba {
	return &ContextedSamba{*s, c, c.Request.Context()}
}

// share

func (s *Samba) CreateShare(payload payloads.RequestCreateShare) (models.Shares, error) {
	shares, err := sambalib.LoadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[payload.Name]; ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", sambalib.ErrShareAlreadyExist.Error(), payload.Name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s already exist", payload.Name)},
		}
	} else if sambalib.IsPathExist(shares, payload.Share) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", sambalib.ErrPathAlreadyExist.Error(), payload.Path),
			Fields:  reply.FieldsError{"path": fmt.Sprintf("%s already exist", payload.Path)},
		}
	}

	shares[payload.Name] = payload.Share

	err = sambalib.SaveSmbConf(sambalib.FilterShares(shares))
	if err != nil {
		return nil, err
	}

	err = iolib.MakeDirWithPerm(payload.Path, payload.Permissions)
	if err != nil {
		return nil, err
	}

	return shares, nil
}

func (s *Samba) GetShares() (models.Shares, error) {
	return sambalib.LoadSmbConf()
}

func (s *Samba) UpdateShare(payload payloads.RequestUpdateShare) (models.Shares, error) {
	shares, err := sambalib.LoadSmbConf()
	if err != nil {
		return nil, err
	}

	oldShare, ok := shares[payload.Name]
	if !ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", sambalib.ErrShareNotExist.Error(), payload.Name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s not exist", payload.Name)},
		}
	}

	delete(shares, payload.Name)
	if sambalib.IsPathExist(shares, payload.Share) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", sambalib.ErrPathAlreadyExist.Error(), payload.Path),
			Fields:  reply.FieldsError{"path": fmt.Sprintf("%s already exist", payload.Path)},
		}
	}

	shares[payload.Name] = payload.Share

	err = sambalib.SaveSmbConf(sambalib.FilterShares(shares))
	if err != nil {
		return nil, err
	}

	if oldShare.Path != payload.Path {
		err = iolib.MakeDirWithPerm(payload.Path, payload.Permissions)
		if err != nil {
			return nil, err
		}
	}

	return shares, nil
}

func (s *Samba) DeleteShare(payload payloads.RequestDeleteShare) (models.Shares, error) {
	shares, err := sambalib.LoadSmbConf()
	if err != nil {
		return nil, err
	}

	share, ok := shares[payload.Name]
	if !ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", sambalib.ErrShareNotExist.Error(), payload.Name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s not exist", payload.Name)},
		}
	}

	err = sambalib.RemoveSmbConf(payload.Name)
	if err != nil {
		return nil, err
	}

	_ = os.RemoveAll(share.Path)

	delete(shares, payload.Name)
	return shares, nil
}

// configuration

func (s *Samba) GetConfig() (models.ShareMap, error) {
	shares, err := sambalib.LoadSmbConfMap()
	if err != nil {
		return nil, err
	}

	global, ok := shares["global"]
	if !ok {
		return nil, sambalib.ErrConfigNotExist
	}
	return global, nil
}

func (s *Samba) UpdateConfig(update models.ShareMap) error {
	shares, err := sambalib.LoadSmbConfMap()
	if err != nil {
		return err
	}

	shares["global"] = update
	return sambalib.SaveSmbConfMap(shares)
}

// backup

func (s *Samba) Backup() error {
	return sambalib.Backup()
}

func (s *Samba) Restore() (models.Shares, error) {
	err := sambalib.Restore()
	if err != nil {
		return nil, err
	}
	return sambalib.LoadSmbConf()
}
