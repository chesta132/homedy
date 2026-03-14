package main

import (
	"homedy/database"
	_ "homedy/flags"
	"log"
	"os"
	"time"

	_ "homedy/config"

	"homedy/internal/routes"
	_ "homedy/internal/services/samba"

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

	{
		routes.RegisterWebsocket(g.Group("/ws"))
		routes.RegisterSamba(g.Group("/samba"))
		routes.RegisterAuth(g.Group("/auth"), db)
	}

	g.Run()
}
