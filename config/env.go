package config

import "os"

var (
	GO_ENV      = os.Getenv("GO_ENV")
	SERVER_PORT = os.Getenv("PORT")

	DB_PATH = os.Getenv("DB_PATH")

	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
	ACCESS_SECRET  = os.Getenv("ACCESS_SECRET")
)

func ReloadEnv() {
	GO_ENV = os.Getenv("GO_ENV")
	SERVER_PORT = os.Getenv("PORT")
	DB_PATH = os.Getenv("DB_PATH")
	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
	ACCESS_SECRET  = os.Getenv("ACCESS_SECRET")
}
