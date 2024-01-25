package websocket

import "github.com/esmailemami/chess/shared/websocket"

var (
	ChessWss = websocket.NewServer(ChessOnMessage)
)

func Run() {
	go ChessWss.Run()
}
