package handlers

import (
	"homedy/internal/services/terminal"
	"homedy/internal/libs/ws"

	"github.com/gin-gonic/gin"
)

type WsTerminal struct{}

func NewWsTerminal() *WsTerminal {
	return &WsTerminal{}
}

func (h *WsTerminal) Handle(c *gin.Context) {
	ws, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	ptmx, cmd, err := terminal.Spawn()
	if err != nil {
		return
	}
	defer ptmx.Close()

	go terminal.SendPTYOutput(ptmx, ws)

	terminal.InputToPTY(ptmx, ws)

	cmd.Wait()
}
