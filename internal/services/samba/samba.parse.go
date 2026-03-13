package samba

import (
	"homedy/config"
	"homedy/internal/libs/cmdlib"

	"gopkg.in/ini.v1"
)

// FIXME: global and other default config reset while save conf

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
	cfg, err := ini.Load(path)
	if err != nil {
		cfg = ini.Empty()
	}
	for name, share := range shares {
		section, err := cfg.NewSection(name)
		if err != nil {
			return err
		}
		if err := section.ReflectFrom(&share); err != nil {
			return err
		}
	}
	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}
	_, err = cmdlib.RestartService("smbd")
	return err
}

func saveSmbConf(shares Shares) error {
	return save(config.SMB_CONF_PATH, shares)
}

func remove(path string, name string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		cfg = ini.Empty()
	}

	cfg.DeleteSection(name)

	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}
	_, err = cmdlib.RestartService("smbd")
	return err
}

func removeSmbConf(name string) error {
	return remove(config.SMB_CONF_PATH, name)
}
