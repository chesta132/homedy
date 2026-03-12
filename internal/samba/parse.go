package samba

import (
	"homedy/config"

	"gopkg.in/ini.v1"
)

func loadConf(path string) (shares Shares, err error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return
	}

	err = cfg.MapTo(shares)
	return
}

func loadSmbConf() (Shares, error) {
	return loadConf(config.SMB_CONF_PATH)
}

func save(path string, shares Shares) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, &shares)
	if err != nil {
		return err
	}
	return cfg.SaveTo(path)
}

func saveSmbConf(shares Shares) error {
	return save(config.SMB_CONF_PATH, shares)
}
