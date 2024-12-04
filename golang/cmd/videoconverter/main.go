package main

import (
	"database/sql"
	"fmt"
	"imersaofc/internal/converter"
	"imersaofc/internal/rabbitmq"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "user")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvOrDefault("POSTGRES_DB", "converter")
	host := getEnvOrDefault("POSTGRES_HOST", "postgres")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")	

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s", user, password, dbname, host, sslmode)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", slog.String("error", err.Error()))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Failed to ping PostgreSQL", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Connected to PostgreSQL successfully")
	return db, nil
}

// Pegar valores default de variaveis de ambiente
func getEnvOrDefault(key, defaultValue string) string { // Vair ler a variavel de ambiente e se ele não existir a gnt pega um valor padrão
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	db, err := connectPostgres()
	if err != nil {
		panic(err)
	}
	rabbitMQURL := getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	rabbitmqClient, err := rabbitmq.NewRabbitClient(rabbitMQURL)
	if err != nil {
		panic(err)
	}
	defer rabbitmqClient.Close()

	convertionExch := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	queueName := getEnvOrDefault("CONVERSION_QUEUE", "video_conversion_queue")
	convesionKey := getEnvOrDefault("CONVERSION_KEY", "conversion")
	confirmationKey := getEnvOrDefault("CONFIMATION_KEY", "finish-conversion")
	confirmationQueue := getEnvOrDefault("CONFIMATION_QUEUE", "video_confirmation_queue")

	vc := converter.NewVideoConverter(rabbitmqClient, db)
	// vc.Handle([]byte(`{"video_id": 1, "path": "/media/uploads/1"}`))

	msgs, err := rabbitmqClient.ConsumeMessages(convertionExch, convesionKey, queueName)
	if err != nil {
		slog.Error("failed to consume messages", slog.String("error", err.Error()))
	}

	// Fica lendo as mensagens que chegam
	for d := range msgs {
		// Gerando uma nova thread (go rotina)
		go func(delivery amqp.Delivery)  {
			vc.Handle(delivery, convertionExch, confirmationKey, confirmationQueue)
		}(d)
	}

}