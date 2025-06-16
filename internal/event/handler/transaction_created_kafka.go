package handler

import (
	"fmt"
	"sync"

	"github.com/bebossi/microservice/pkg/events"
	"github.com/bebossi/microservice/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *TransactionCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup)  {
	defer wg.Done()

	h.Kafka.Publish(message, nil, "transactions")
	fmt.Println("Transaction created", message)
}

