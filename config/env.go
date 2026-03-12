package config

import "os"

var (
	GO_ENV      = os.Getenv("GO_ENV")
	SERVER_PORT = os.Getenv("PORT")
)
