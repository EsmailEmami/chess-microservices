package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

const (
	Direct  = "direct"
	Fanout  = "fanout"
	Topic   = "topic"
	Headers = "headers"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(username, password, address string) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s", username, password, address)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &RabbitMQ{conn: conn, ch: ch}, nil
}

func (rmq *RabbitMQ) Close() {
	if rmq.ch != nil {
		rmq.ch.Close()
	}
	if rmq.conn != nil {
		rmq.conn.Close()
	}
}

func (rmq *RabbitMQ) DeclareExchange(exchangeName, exchangeType string, durable, autoDelete bool) error {
	err := rmq.ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // exchange type
		durable,      // durable
		autoDelete,   // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}
	return nil
}

func (rmq *RabbitMQ) DeclareQueue(queueName string, durable, autoDelete, exclusive bool) (amqp.Queue, error) {
	q, err := rmq.ch.QueueDeclare(
		queueName,  // queue name
		durable,    // durable, persisted the queue when the broker is restarted
		autoDelete, // auto-delete, will be delete when the project is shutted down
		exclusive,  // exclusive, will make the queue only avaiable for the creator connection
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue: %v", err)
	}
	return q, nil
}

func (rmq *RabbitMQ) BindQueueToExchange(queueName, exchangeName, routingKey string) error {
	err := rmq.ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %v", err)
	}
	return nil
}

const (
	DeliveryModeTransient  uint8 = 0
	DeliveryModePersistent uint8 = 2
)

type Message struct {
	Body any

	Headers amqp.Table

	// Properties
	DeliveryMode  uint8  // Transient (0 or 1) or Persistent (2)
	CorrelationId string // correlation identifier
	ReplyTo       string // address to to reply to (ex: RPC)
	MessageId     string // message identifier
	Type          string // message type name
	UserId        string // creating user id - ex: "guest"
	AppId         string // creating application id

}

func (rmq *RabbitMQ) PublishMessage(ctx context.Context, exchange, routingKey string, msg *Message) error {

	bts, err := json.Marshal(&msg.Body)

	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = rmq.ch.PublishWithContext(ctx,
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(bts),
			DeliveryMode:  msg.DeliveryMode,
			Timestamp:     time.Now(),
			CorrelationId: msg.CorrelationId,
			ReplyTo:       msg.ReplyTo,
			MessageId:     msg.MessageId,
			Type:          msg.Type,
			UserId:        msg.UserId,
			AppId:         msg.AppId,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (rmq *RabbitMQ) ConsumeMessages(queueName string, autoAck bool) (<-chan amqp.Delivery, error) {
	msgs, err := rmq.ch.Consume(
		queueName, // queue
		"",        // consumer
		autoAck,   // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %v", err)
	}
	return msgs, nil
}
