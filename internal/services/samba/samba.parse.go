package samba

import (
	"homedy/config"

	"gopkg.in/ini.v1"
)

func loadConf(path string) (Shares, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	shares := make(Shares)
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		var share Share
		if err := section.MapTo(&share); err != nil {
			return nil, err
		}
		shares[section.Name()] = share
	}
	return shares, nil
}

func loadSmbConf() (Shares, error) {
	return loadConf(config.SMB_CONF_PATH)
}

func save(path string, shares Shares) error {
	cfg := ini.Empty()
	for name, share := range shares {
		section, err := cfg.NewSection(name)
		if err != nil {
			return err
		}
		if err := section.ReflectFrom(&share); err != nil {
			return err
		}
	}
	return cfg.SaveTo(path)
}

func saveSmbConf(shares Shares) error {
	return save(config.SMB_CONF_PATH, shares)
}
