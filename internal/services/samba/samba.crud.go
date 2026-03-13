package samba

import (
	"errors"
	"fmt"
)

func AddConf(name string, share Share) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[name]; ok {
		return nil, errors.Join(ErrShareAlreadyExist, fmt.Errorf(": %s", name))
	}
	shares[name] = share

	return shares, saveSmbConf(FilterShares(shares))
}

func ReadConf() (Shares, error) {
	return loadSmbConf()
}

func UpdateConf(name string, share Share) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[name]; !ok {
		return nil, errors.Join(ErrShareNotExist, fmt.Errorf(": %s", name))
	}
	shares[name] = share

	return shares, saveSmbConf(FilterShares(shares))
}

func DeleteConf(name string) (Shares, error) {
	shares, err := loadSmbConf()
	if err != nil {
		return nil, err
	}

	if _, ok := shares[name]; !ok {
		return nil, errors.Join(ErrShareNotExist, fmt.Errorf(": %s", name))
	}

	return shares, saveSmbConf(FilterShares(shares))
}
