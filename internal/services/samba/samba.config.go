package samba

func GetConfiguration() (ShareMap, error) {
	shares, err := loadSmbConfMap()
	if err != nil {
		return nil, err
	}

	global, ok := shares["global"]
	if !ok {
		return nil, ErrConfigNotExist
	}
	return global, nil
}

func UpdateConfig(update ShareMap) error {
	shares, err := loadSmbConfMap()
	if err != nil {
		return err
	}

	shares["global"] = update
	return saveSmbConfMap(shares)
}
