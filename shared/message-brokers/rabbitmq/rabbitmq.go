package rabbitmq

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

const (
	Direct  = "direct"
	Fanout  = "fanout"
	Topic   = "topic"
	Headers = "headers"
)

// RabbitMQ represents a simple RabbitMQ connection.
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// New creates a new RabbitMQ instance and initializes the connection and channel.
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

// Close closes the RabbitMQ connection and channel.
func (rmq *RabbitMQ) Close() {
	if rmq.ch != nil {
		rmq.ch.Close()
	}
	if rmq.conn != nil {
		rmq.conn.Close()
	}
}

// DeclareExchange declares an exchange with the given name and type.
func (rmq *RabbitMQ) DeclareExchange(exchangeName, exchangeType string) error {
	err := rmq.ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // exchange type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}
	return nil
}

// DeclareQueue declares a queue with the given name.
func (rmq *RabbitMQ) DeclareQueue(queueName string, durable, autoDelete, exclusive bool) (amqp.Queue, error) {
	q, err := rmq.ch.QueueDeclare(
		queueName,  // queue name
		durable,    // durable
		autoDelete, // auto-delete
		exclusive,  // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare queue: %v", err)
	}
	return q, nil
}

// BindQueueToExchange binds a queue to an exchange with a specific routing key.
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

// PublishMessage publishes a message to the specified exchange with a routing key.
func (rmq *RabbitMQ) PublishMessage(ctx context.Context, exchange, routingKey string, data any) error {

	bts, err := json.Marshal(&data)

	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	err = rmq.ch.PublishWithContext(ctx,
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bts),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

// ConsumeMessages consumes messages from the specified queue.
func (rmq *RabbitMQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := rmq.ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
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

// ConsumeMessagesFromDirect consumes messages from a queue bound to a direct exchange with a routing key.
func (rmq *RabbitMQ) ConsumeMessagesFromExchange(queueName, exchangeName, routingKey string) (<-chan amqp.Delivery, error) {
	// Declare the queue
	queue, err := rmq.DeclareQueue(queueName, true, false, false)
	if err != nil {
		return nil, err
	}

	// Bind the queue to the topic exchange with the routing pattern
	err = rmq.BindQueueToExchange(queue.Name, exchangeName, routingKey)
	if err != nil {
		return nil, err
	}

	// Consume messages from the queue
	msgs, err := rmq.ConsumeMessages(queue.Name)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
