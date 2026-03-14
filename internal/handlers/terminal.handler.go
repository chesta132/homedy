package handlers

import (
	"homedy/internal/libs/ws"
	"homedy/internal/services"

	"github.com/gin-gonic/gin"
)

type WsTerminal struct {
	terminalSvc *services.Terminal
}

func NewWsTerminal(terminalSvc *services.Terminal) *WsTerminal {
	return &WsTerminal{terminalSvc}
}

func (h *WsTerminal) Handle(c *gin.Context) {
	ws, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	ptmx, cmd, err := h.terminalSvc.Spawn()
	if err != nil {
		return
	}
	defer ptmx.Close()

	go h.terminalSvc.SendPTYOutput(ptmx, ws)

	h.terminalSvc.InputToPTY(ptmx, ws)

	cmd.Wait()
}
