package samba

import (
	"fmt"
	"homedy/internal/libs/iolib"
	"homedy/internal/libs/replylib"
	"os"

	"github.com/chesta132/goreply/reply"
)

func AddShare(name string, share Share) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[name]; ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", ErrShareAlreadyExist.Error(), name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s already exist", name)},
		}
	} else if isPathExist(shares, share) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", ErrPathAlreadyExist.Error(), share.Path),
			Fields:  reply.FieldsError{"path": fmt.Sprintf("%s already exist", share.Path)},
		}
	}

	shares[name] = share

	err = iolib.MakeDirWithPerm(share.Path, share.Permissions)
	if err != nil {
		return nil, err
	}

	return shares, saveSmbConf(FilterShares(shares))
}

func ReadShare() (Shares, error) {
	return loadSmbConf()
}

func UpdateShare(name string, share Share) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	oldShare, ok := shares[name]
	if !ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", ErrShareNotExist.Error(), name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s not exist", name)},
		}
	}

	delete(shares, name)
	if isPathExist(shares, share) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", ErrPathAlreadyExist.Error(), share.Path),
			Fields:  reply.FieldsError{"path": fmt.Sprintf("%s already exist", share.Path)},
		}
	}

	shares[name] = share

	if oldShare.Path != share.Path {
		err = iolib.MakeDirWithPerm(share.Path, share.Permissions)
		if err != nil {
			return nil, err
		}
	}

	return shares, saveSmbConf(FilterShares(shares))
}

func DeleteShare(name string) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	share, ok := shares[name]
	if !ok {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: fmt.Sprintf("%s: %s", ErrShareNotExist.Error(), name),
			Fields:  reply.FieldsError{"name": fmt.Sprintf("%s not exist", name)},
		}
	}
	delete(shares, name)

	_ = os.RemoveAll(share.Path)

	return shares, removeSmbConf(name)
}
