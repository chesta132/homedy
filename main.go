package main

import (
	"embed"
	_ "homedy/flags"
	"io/fs"
	"net/http"
	"strings"

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

	{
		router := routes.New(db, repos.New(db))
		router.RegisterAuth(api.Group("/auth"))
		router.RegisterWebsocket(api.Group("/ws"))
		router.RegisterSamba(api.Group("/samba"))
		router.RegisterAuth(api.Group("/auth"))
	}

	dist, _ := fs.Sub(frontendFiles, "ui/dist")
	fileServer := http.FileServer(http.FS(dist))

	g.Use(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			// Check if the requested file exists
			_, err := fs.Stat(dist, strings.TrimPrefix(c.Request.URL.Path, "/"))
			if os.IsNotExist(err) {
				// If the file does not exist, serve index.html
				c.Request.URL.Path = "/"
			}

			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})

	g.Run(":" + config.SERVER_PORT)
}
