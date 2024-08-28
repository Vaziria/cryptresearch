package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (cfg *RabbitConfig) GetUri() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
}

func (cfg *RabbitConfig) CreateConnection() *amqp091.Connection {
	conn, err := amqp091.Dial(cfg.GetUri())
	if err != nil {
		log.Panicln(err)
	}

	return conn
}

func GetRabbitConfig() *RabbitConfig {
	getKey := func(name string, defval string) string {
		value := os.Getenv(name)
		if value == "" {
			return defval
		}
		return value
	}

	return &RabbitConfig{
		Username: getKey("RABBITMQ_USER", "guest"),
		Password: getKey("RABBITMQ_PASS", "guest"),
		Host:     getKey("RABBITMQ_HOST", "localhost"),
		Port:     getKey("RABBITMQ_HOST", "5672"),
	}
}
