package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"golang.org/x/text/message"
)

type RabbitClient struct {
	conn *amqp.Connection
	channel *amqp.Channel
	url string
}

// Abrindo conexão com o RabbitMQ
func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
		}

		return conn, channel, nil
}

// Criar um novo client para ele abrir a conexão
func NewRabbitClient(connectionURL string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connectionURL)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn: conn,
		channel: channel,
		url: connectionURL,
	}, nil
}

// Responsável por consumir novas mensagens. Vai ser passado uma fila e ele vai ler a mensagem
func (client *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	// Declarar uma exchange
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declarar uma fila
	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Fazer o bind (toda vez que uma mensagem for enviada para a minha exchange essa mensagem vai ser roteada para determinada fila)
	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := client.channel.Consume(
		queueName,
		"goapp",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue: %v", err)
	}

	return msgs, nil
}

func (client *RabbitClient) PublishMessage(exchange, routingKey, queueName, message []byte) error {
	// É necessário declarar a exchange, fila e key para garantir que não tenha problema (vai ser uma fila de confirmação)
	// A fila que vai ler vai a ser a da aplicação django a aplicação go vai publicar

	// Declarar uma exchange
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declarar uma fila
	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)

	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	err = client.channel.Publish(
		exchange, routingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body: message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish queue: %v", err)
	}
	return nil
}

// Fechar a conexão e o canal
func (client *RabbitClient) Close() {
	client.channel.Close()
	client.conn.Close()
}