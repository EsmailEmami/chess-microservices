package rabbitmq

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
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
}

func PublishUserLastConnection(ctx context.Context, userID uuid.UUID, date time.Time) error {
	return amqp.PublishMessage(ctx, userExchange, userLastConnectionRouteKey, &struct {
		UserID uuid.UUID `json:"userId"`
		Date   time.Time `json:"date"`
	}{
		UserID: userID,
		Date:   date,
	})
}
