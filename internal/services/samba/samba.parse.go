package samba

import (
	"homedy/config"

	"gopkg.in/ini.v1"
)

func loadConfMap(path string) (ShareMaps, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	cfgMap := make(ShareMaps)
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			continue
		}
		val := make(ShareMap)
		for _, key := range section.Keys() {
			val[key.Name()] = key.Value()
		}
		cfgMap[section.Name()] = val
	}
	return cfgMap, nil
}

func loadSmbConfMap() (ShareMaps, error) {
	return loadConfMap(config.SMB_CONF_PATH)
}

func saveMap(path string, maps ShareMaps) error {
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

func saveSmbConfMap(maps ShareMaps) error {
	return saveMap(config.SMB_CONF_PATH, maps)
}

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
	return restartService()
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
	return restartService()
}

func removeSmbConf(name string) error {
	return remove(config.SMB_CONF_PATH, name)
}
