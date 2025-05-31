package gateway

import "github.com/bebossi/microservice/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
