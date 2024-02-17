package rabbitmq

import (
	"github.com/esmailemami/chess/shared/logging"
	"github.com/esmailemami/chess/shared/message-brokers/rabbitmq"
	"github.com/spf13/viper"
)

// this rabbitmq connection is only used for consuming messages from producers
var amqp *rabbitmq.RabbitMQ

func InitializeConsumerConnection() {
	var (
		username = viper.GetString("rabbitmq.username")
		password = viper.GetString("rabbitmq.password")
		address  = viper.GetString("rabbitmq.address")
	)

	amqpConn, err := rabbitmq.New(username, password, address)

	if err != nil {
		logging.FatalE("rabbit MQ connnection failed", err)
	}

	amqp = amqpConn

	// initialize exchanges and queues
	initializeMediaRabbitMQ()
	initializeUserRabbitMQ()

	logging.Info("rabbit MQ connected")
}
