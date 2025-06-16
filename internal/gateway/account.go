package gateway

import "github.com/bebossi/microservice/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	UpdateBalance(account *entity.Account) error
	FindByID(clientID string) (*entity.Account, error)
}
