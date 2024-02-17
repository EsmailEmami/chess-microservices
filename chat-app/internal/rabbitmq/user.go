package rabbitmq

import (
	"github.com/esmailemami/chess/shared/logging"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	userExchange = "chat_user_events"
)

var (
	UserProfileChangedCh = make(chan *UserProfileChangedMessage, 256)
)

func initializeUserRabbitMQ() {
	go consumeUserProfileChanged()
}

func consumeUserProfileChanged() {
	queue, err := amqp.DeclareQueue("chat_user_profile_changed", true, false, false)
	if err != nil {
		logging.FatalE("failed to declare 'chat_user_profile_changed' queue", err)
	}

	if err := amqp.BindQueueToExchange(queue.Name, userExchange, "chat_user.user.profile.changed"); err != nil {
		logging.FatalE("failed to bind queue", err)
	}

	messageBus, err := amqp.ConsumeMessages(queue.Name, false)
	if err != nil {
		logging.FatalE("failed to consume 'chat_user_profile_changed' queue", err)
	}

	for msg := range messageBus {
		var data UserProfileChangedMessage
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			logging.ErrorE("failed to unmarshal user profile consumer", err)
		}

		logging.Debug("received user profile changed")

		UserProfileChangedCh <- &data
	}
}

type UserProfileChangedMessage struct {
	UserID      uuid.UUID `json:"userID"`
	ProfilePath string    `json:"profilePath"`
}
