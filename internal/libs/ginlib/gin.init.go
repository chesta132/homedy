package ginlib

import (
	"homedy/config"

	"github.com/gin-gonic/gin"
)

func init() {
	if config.IsEnvProd() {
		gin.SetMode(gin.ReleaseMode)
	}
}
