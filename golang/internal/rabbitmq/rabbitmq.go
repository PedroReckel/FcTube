package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitClient encapsula a conexão e o canal do RabbitMQ
type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// Abrindo conexão com o RabbitMQ
func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	// Conecta ao RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Abre um canal na conexão
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close() // Garante que a conexão seja fechada se o canal não abrir
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
		conn:    conn,
		channel: channel,
		url:     connectionURL,
	}, nil
}

// Responsável por consumir novas mensagens. Vai ser passado uma fila e ele vai ler a mensagem
func (client *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	// Declarar uma exchange
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,  // Durable
		false, // Auto-delete corrigido
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Declarar uma fila
	queue, err := client.channel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete corrigido
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Fazer o bind (toda vez que uma mensagem for enviada para a minha exchange essa mensagem vai ser roteada para determinada fila)
	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	// Consumir mensagens
	msgs, err := client.channel.Consume(
		queueName,
		"goapp", // Consumer Tag
		false,   // Auto-ack
		false,   // Exclusive
		false,   // No-local
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue: %v", err)
	}

	return msgs, nil
}

// Publicar mensagens na fila
func (client *RabbitClient) PublishMessage(exchange, routingKey, queueName string, message []byte) error {
	// Declarar uma exchange
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,  // Durable
		false, // Auto-delete corrigido
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	// Publicar mensagem
	err = client.channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

// Fechar a conexão e o canal
func (client *RabbitClient) Close() {
	if client.channel != nil {
		_ = client.channel.Close() // Garante que erros não interrompam a execução
	}
	if client.conn != nil {
		_ = client.conn.Close()
	}
}
