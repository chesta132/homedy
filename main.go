package main

import (
	"homedy/database"
	_ "homedy/flags"

	_ "homedy/config"

	"homedy/internal/routes"
	_ "homedy/internal/services/samba"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := database.Connect()
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	{
		routes.RegisterWebsocket(g.Group("/ws"))
		routes.RegisterSamba(g.Group("/samba"))
	}

	g.Run()
}
