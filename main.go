package main

import (
	_ "homedy/config"

	"homedy/internal/routes"
	_ "homedy/internal/services/samba"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	{
		routes.RegisterSamba(g.Group("/samba"))
	}

	g.Run()
}
