package routes

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (rt *Router) RegisterFrontend(fsys fs.FS, dir string) {
	dist, _ := fs.Sub(fsys, dir)
	fileServer := http.FileServer(http.FS(dist))

	rt.g.Use(func(c *gin.Context) {
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
}
