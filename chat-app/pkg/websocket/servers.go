package websocket

import (
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/esmailemami/chess/shared/websocket"
	"github.com/gin-gonic/gin"
)

var (
	GlobalRoomWss = websocket.NewServer(GlobalRoomOnMessage)
)

func Run() {
	go GlobalRoomWss.Run()
}

func init() {
	GlobalRoomWss.OnRegister(GlobalRoomOnRegister)
	GlobalRoomWss.OnUnregister(GlobalRoomOnUnregister)
}

func InitializeRoutes(r *gin.Engine) {
	ws := r.Group("/ws")
	ws.Use(middleware.Authorization())

	ws.GET("/global-room", GlobalRoomWss.HandleWS)
}
