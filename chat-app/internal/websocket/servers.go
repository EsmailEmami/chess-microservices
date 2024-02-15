package websocket

import (
	"context"
	"time"

	"github.com/esmailemami/chess/chat/pkg/rabbitmq"
	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/middleware"
	"github.com/esmailemami/chess/shared/websocket"
	"github.com/gin-gonic/gin"
)

var (
	GlobalRoomWss      = websocket.NewServer(GlobalRoomOnMessage)
	PublicChatRoomWss  = websocket.NewServer(PublicChatRoomOnMessage)
	PrivateChatRoomWss = websocket.NewServer(PublicChatRoomOnMessage)
)

func Run() {
	go GlobalRoomWss.Run()
	go PublicChatRoomWss.Run()
	go PrivateChatRoomWss.Run()
}

func init() {
	GlobalRoomWss.OnRegister(GlobalRoomOnRegister)
	GlobalRoomWss.OnUnregister(GlobalRoomOnUnregister)

	PublicChatRoomWss.OnRegister(PublicChatRoomOnRegister)
	PublicChatRoomWss.OnUnregister(PublicChatRoomOnUnregister)

	PrivateChatRoomWss.OnRegister(PrivateChatRoomOnRegister)
	PrivateChatRoomWss.OnUnregister(PrivateChatRoomOnUnregister)

	PublicChatRoomWss.OnPong(func(c *websocket.Client) {
		if err := rabbitmq.PublishUserLastConnection(context.Background(), c.UserID, time.Now()); err != nil {
			logging.WarnE("failed to send user last connection date", err)
		}
	})
}

func InitializeRoutes(r *gin.Engine) {
	ws := r.Group("/ws")
	ws.Use(middleware.Authorization())

	ws.GET("/global", GlobalRoomWss.HandleWS)
	ws.GET("/public", PublicChatRoomWss.HandleWS)
	ws.GET("/private", PrivateChatRoomWss.HandleWS)
}
