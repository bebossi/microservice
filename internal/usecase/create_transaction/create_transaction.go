package create_transaction

import (
	"fmt"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/bebossi/microservice/internal/gateway"
	"github.com/bebossi/microservice/pkg/events"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID string `json:"id"`
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	transactionGateway gateway.TransactionGateway,
	 accountGateway gateway.AccountGateway, 
	 eventDispatcher events.EventDispatcherInterface,
	 transactionCreated events.EventInterface,
	) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}


func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	accountFrom, err := uc.AccountGateway.FindByClientID(input.AccountIDFrom)
	if err != nil {
		fmt.Println("err find account from", err)
		return nil, err
	}
	accountTo, err := uc.AccountGateway.FindByClientID(input.AccountIDTo)
	if err != nil {
		fmt.Println("err find account to", err)
		return nil, err
	}
	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		fmt.Println("err create transaction", err)
		return nil, err
	}
	err = uc.TransactionGateway.Create(transaction)
	if err != nil {
		fmt.Println("err create transaction", err)
		return nil, err
	}
	output := &CreateTransactionOutputDTO{
		ID: transaction.ID,
	}

	uc.TransactionCreated.SetPayload(transaction)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	return output, nil
}
