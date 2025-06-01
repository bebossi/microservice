package gateway

import "github.com/bebossi/microservice/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	UpdateBalance(account *entity.Account) error
	FindByClientID(clientID string) (*entity.Account, error)
}
