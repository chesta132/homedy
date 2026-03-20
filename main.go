package main

import (
	"embed"
	_ "homedy/flags"

	"homedy/config"

	"homedy/database"
	"log"
	"os"
	"time"

	_ "homedy/internal/libs/ginlib"
	_ "homedy/internal/libs/sambalib"
	_ "homedy/internal/libs/validatorlib"
	"homedy/internal/middlewares"
	"homedy/internal/repos"
	"homedy/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed ui/dist
var frontendFiles embed.FS

func main() {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := database.Connect(&gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(err)
	}

	g := gin.Default()
	api := g.Group("/api")
	api.Use(middlewares.LimitTotalUploadSize(config.LIMIT_UPLOAD_SIZE))
	amw := middlewares.NewAuth(repos.NewRevoke(db))

	{
		router := routes.New(g, db, repos.New(db))
		router.RegisterAuth(api.Group("/auth"))

		api.Use(amw.Protected())
		router.RegisterWebsocket(api.Group("/ws"))
		router.RegisterSamba(api.Group("/samba"))
		router.RegisterConverter(api.Group("/convert"))

		router.RegisterFrontend(frontendFiles, "ui/dist")
	}

	g.Run(":" + config.SERVER_PORT)
}
