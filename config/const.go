package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// general config

	APP_NAME = "Homedy"
	TEMP_DIR = "/tmp/homedy"

	// samba config

	SAMBA_PKG_NAME       string = "samba"
	SAMBA_PKG_VER        string = "2:4.19.5+dfsg-4ubuntu9.4"
	SAMBA_PKG            string = SAMBA_PKG_NAME + "=" + SAMBA_PKG_VER
	SMB_CONF_PATH        string = "/etc/samba/smb.conf"
	SMB_CONF_BACKUP_PATH string = "/etc/samba/smb.conf.backup"

	// libreoffice config

	LIBRE_PKG_NAME = "libreoffice"
	LIBRE_PKG_VER  = "4:24.2.7-0ubuntu0.24.04.4"
	LIBRE_PKG      = LIBRE_PKG_NAME + "=" + LIBRE_PKG_VER

	// logger config

	LOGGER_TIME_FORMAT string = "2006-01-02 15:04:05"

	// auth token (jwt)

	REFRESH_TOKEN_KEY = "refresh_token"
	ACCESS_TOKEN_KEY  = "access_token"

	// batch

	MAX_CREATE_BATCH              = 100 // max data per batch when inserting multiple records
	LIMIT_RESOURCE_PER_PAGINATION = 40

	// terminal remote

	TERMINAL_RESTRICTED bool = false

	// uploads

	LIMIT_UPLOAD_SIZE = 200 << 20 // 200MB per request

	// app secret

	APP_SECRET_WS_SUBPROTOCOL_KEY = "app-secret"
	APP_SECRET_HEADER_KEY         = "X-APP-SECRET"
)

var (
	// auth token (jwt)

	SIGN_METHOD                jwt.SigningMethod = jwt.SigningMethodHS256
	REFRESH_TOKEN_EXPIRY       time.Duration     = (time.Hour * 24 * 7 * 2) + (time.Hour * 24 * 3) // 2 weeks 3 days
	ACCESS_TOKEN_EXPIRY        time.Duration     = time.Minute * 5                                 // 5 minutes
	ROTATE_REFRESH_TOKEN_AFTER time.Duration     = time.Hour * 24 * 7 * 2                          // 2 weeks

	// convert

	ConvertFileLimits = map[string]int64{
		"html": 2 << 20,  // 2MB
		"md":   2 << 20,  // 2MB
		"csv":  5 << 20,  // 5MB
		"xlsx": 10 << 20, // 10Mb
		"docx": 10 << 20, // 10MB
		"pptx": 10 << 20, // 10MB
		"pdf":  15 << 20, // 15MB
	}
)
