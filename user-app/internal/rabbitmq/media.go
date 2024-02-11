package rabbitmq

import "github.com/esmailemami/chess/shared/logging"

func initializeMediaRabbitMQ() {
	go consumeMediaUserProfile()
}

func consumeMediaUserProfile() {
	msgsCh, err := amqp.ConsumeMessagesFromTopicWithRoutingKey("user_queue", mediaExchange, "media.profile")

	if err != nil {
		logging.ErrorE("failed to consume rabbit MQ", err)
		return
	}

	//userService := service.NewUserService()

	for msg := range msgsCh {
		logging.Info("Message received", "msg", string(msg.Body))
	}
}
