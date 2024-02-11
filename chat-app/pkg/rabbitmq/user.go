package rabbitmq

import (
	"context"
	"time"

	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/google/uuid"
)

const (
	userExchange = "chat_user_ex"
)

func initializeUserRabbitMQ() {
	amqp.DeclareExchange(userExchange, rabbitmq.Topic)

	queue, err := amqp.DeclareQueue("chat_user_queue", true, false, false)

	if err != nil {
		logging.FatalE("failed to declare 'chat_user_queue' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, userExchange, "chat_user.*"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}
}

func PublishUserLastConnection(ctx context.Context, userID uuid.UUID, date time.Time) error {
	return amqp.PublishMessage(ctx, userExchange, "chat_user.last_connection", &struct {
		UserID uuid.UUID `json:"userId"`
		Date   time.Time `json:"date"`
	}{
		UserID: userID,
		Date:   date,
	})
}
