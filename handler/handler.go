package handler

import "github.com/confluentinc/confluent-kafka-go/kafka"

type Handler interface {
	Topic() string
	HandleMessage(message *kafka.Message) error
}
