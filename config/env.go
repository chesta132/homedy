package config

import "os"

var (
	GO_ENV      = os.Getenv("HOMEDY_ENV")
	SERVER_PORT = os.Getenv("HOMEDY_PORT")

	DB_HOST = os.Getenv("HOMEDY_DB_PATH")
	DB_PORT = os.Getenv("HOMEDY_DB_PORT")
	DB_USER = os.Getenv("HOMEDY_DB_USER")
	DB_PASS = os.Getenv("HOMEDY_DB_PASS")
	DB_NAME = os.Getenv("HOMEDY_DB_NAME")

	REFRESH_SECRET = os.Getenv("HOMEDY_REFRESH_SECRET")
	ACCESS_SECRET  = os.Getenv("HOMEDY_ACCESS_SECRET")
)

func ReloadEnv() {
	GO_ENV = os.Getenv("GO_ENV")
	SERVER_PORT = os.Getenv("PORT")

	DB_HOST = os.Getenv("HOMEDY_DB_PATH")
	DB_PORT = os.Getenv("HOMEDY_DB_PORT")
	DB_USER = os.Getenv("HOMEDY_DB_USER")
	DB_PASS = os.Getenv("HOMEDY_DB_PASS")
	DB_NAME = os.Getenv("HOMEDY_DB_NAME")

	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
	ACCESS_SECRET = os.Getenv("ACCESS_SECRET")
}
