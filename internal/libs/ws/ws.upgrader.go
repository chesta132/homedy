package ws

import (
	"homedy/config"
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin")  == config.FRONTEND_URL
	},
	Subprotocols: []string{config.APP_SECRET_WS_SUBPROTOCOL_KEY},
}
