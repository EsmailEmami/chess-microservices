package rabbitmq

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/esmailemami/chess/user/internal/app/service"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	userExchange                = "user_exchange"
	userLastConnectionDateQueue = "last_connection_date_queue"
	userLastConnectionRouteKey  = "last_connection"
)

func initializeUserRabbitMQ() {
	amqp.DeclareExchange(userExchange, rabbitmq.Direct)

	queue, err := amqp.DeclareQueue(userLastConnectionDateQueue, true, false, false)

	if err != nil {
		logging.FatalE("failed to declare 'last_connection_date_queue' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, userExchange, userLastConnectionRouteKey); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	go consumeLastConnection()
}

type lastConnection struct {
	UserID uuid.UUID `json:"userId"`
	Date   time.Time `json:"date"`
}

func consumeLastConnection() {
	msgsCh, err := amqp.ConsumeMessages(userLastConnectionDateQueue)

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	userService := service.NewUserService()

	for msg := range msgsCh {
		var data lastConnection
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			logging.ErrorE("failed to unmarshal last connection consumer", err)
		}

		if err := userService.UpdateLastConnection(context.Background(), data.UserID, data.Date); err != nil {
			logging.Error("failed to update user connection", "userId", data.UserID)
		}
	}
}
