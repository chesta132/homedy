package terminal

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func Spawn() (ptmx *os.File, cmd *exec.Cmd, err error) {
	cmd = exec.Command("bash")
	cmd.Env = append(os.Environ(), "TERM=xterm-color")

	ptmx, err = pty.Start(cmd)
	return
}

func SendPTYOutput(ptmx *os.File, ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ptmx.Read(buf)
		if err != nil {
			break
		}
		ws.WriteMessage(websocket.BinaryMessage, buf[:n])
	}
}

func InputToPTY(ptmx *os.File, ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		ptmx.Write(msg)
	}
}
