package services

import amqp "github.com/rabbitmq/amqp091-go"

var (
	GlobalChannel *amqp.Channel
)
