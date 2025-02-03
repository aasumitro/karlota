package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQConnection() Option {
	return func(cfg *Config) {
		log.Println("Trying to open rabbitmq connection pool . . . .")
		pubConn, err := amqp.Dial(cfg.RabbitMQDsnURL)
		if err != nil {
			log.Fatalf("RABBITMQ_PUBLISHER_ERROR: %s", err.Error())
		}
		cfg.Infra.RabbitMQPublisher = pubConn
		subConn, err := amqp.Dial(cfg.RabbitMQDsnURL)
		if err != nil {
			log.Fatalf("RABBITMQ_SUBSCRIBER_ERROR: %s", err.Error())
		}
		cfg.Infra.RabbitMQConsumer = subConn
		log.Println("RabbitMQ connection pool created . . . .")
	}
}
