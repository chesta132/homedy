package main

import (
	_ "homedy/flags"

	"homedy/config"

	"homedy/database"
	"log"
	"os"
	"time"

	_ "homedy/internal/libs/ginlib"
	_ "homedy/internal/libs/sambalib"
	_ "homedy/internal/libs/validatorlib"
	"homedy/internal/repos"
	"homedy/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
	router := routes.New(db, repos.New(db))

	{
		router.RegisterWebsocket(g.Group("/ws"))
		router.RegisterSamba(g.Group("/samba"))
		router.RegisterAuth(g.Group("/auth"))
	}

	g.Run(":" + config.SERVER_PORT)
}
