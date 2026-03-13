package samba

import (
	"fmt"
	"homedy/internal/libs/iolib"
	"os"
)

func AddShare(name string, share Share) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[name]; ok {
		return nil, fmt.Errorf("%w: %s", ErrShareAlreadyExist, name)
	} else if isPathExist(shares, share) {
		return nil, fmt.Errorf("%w: %s", ErrPathAlreadyExist, share.Path)
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
		return nil, fmt.Errorf("%w: %s", ErrShareNotExist, name)
	}

	delete(shares, name)
	if isPathExist(shares, share) {
		return nil, fmt.Errorf("%w: %s", ErrPathAlreadyExist, share.Path)
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
		return nil, fmt.Errorf("%w: %s", ErrShareNotExist, name)
	}
	delete(shares, name)

	_ = os.RemoveAll(share.Path)

	return shares, saveSmbConf(FilterShares(shares))
}
