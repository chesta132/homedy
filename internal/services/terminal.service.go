package services

import (
	"context"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Terminal struct{}

type ContextedTerminal struct {
	Terminal
	c   *gin.Context
	ctx context.Context
}

func NewTerminal() *Terminal {
	return &Terminal{}
}

func (s *Terminal) AttachContext(c *gin.Context) *ContextedTerminal {
	return &ContextedTerminal{*s, c, c.Request.Context()}
}

func (s *Terminal) Spawn() (ptmx *os.File, cmd *exec.Cmd, err error) {
	cmd = exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-color")

	ptmx, err = pty.Start(cmd)
	return
}

func (s *Terminal) SendPTYOutput(ptmx *os.File, ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			break
		}
		ws.WriteMessage(websocket.BinaryMessage, buf[:n])
	}
}

func (s *Terminal) InputToPTY(ptmx *os.File, ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		ptmx.Write(msg)
	}
}
