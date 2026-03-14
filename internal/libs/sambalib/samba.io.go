package sambalib

import (
	"homedy/config"
	"homedy/internal/models"

	"gopkg.in/ini.v1"
)

// models.ShareMaps

func LoadConfMap(path string) (models.ShareMaps, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	cfgMap := make(models.ShareMaps)
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		val := make(models.ShareMap)
		for _, key := range section.Keys() {
			val[key.Name()] = key.Value()
		}
		cfgMap[section.Name()] = val
	}
	return cfgMap, nil
}

func LoadSmbConfMap() (models.ShareMaps, error) {
	return LoadConfMap(config.SMB_CONF_PATH)
}

func SaveMap(path string, maps models.ShareMaps) error {
	cfg, err := ini.Load(path)
	if err != nil {
		cfg = ini.Empty()
	}
	for name, share := range maps {
		section, err := cfg.NewSection(name)
		if err != nil {
			return err
		}
		for k, v := range share {
			section.NewKey(k, v)
		}
	}
	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}
	return restartService()
}

func SaveSmbConfMap(maps models.ShareMaps) error {
	return SaveMap(config.SMB_CONF_PATH, maps)
}

// models.Shares

func LoadConf(path string) (models.Shares, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	shares := make(models.Shares)
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		var share models.Share
		if err := section.MapTo(&share); err != nil {
			return nil, err
		}
		shares[section.Name()] = share
	}
	return shares, nil
}

func LoadSmbConf() (models.Shares, error) {
	return LoadConf(config.SMB_CONF_PATH)
}

func Save(path string, shares models.Shares) error {
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
	return restartService()
}

func SaveSmbConf(shares models.Shares) error {
	return Save(config.SMB_CONF_PATH, shares)
}

// remove

func Remove(path string, name string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		cfg = ini.Empty()
	}

	cfg.DeleteSection(name)

	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}
	return restartService()
}

func RemoveSmbConf(name string) error {
	return Remove(config.SMB_CONF_PATH, name)
}
