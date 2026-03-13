package samba

func GetConfiguration() (map[string]string, error) {
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
