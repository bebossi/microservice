package create_transaction

import (
	"context"
	"fmt"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/bebossi/microservice/internal/gateway"
	"github.com/bebossi/microservice/pkg/events"
	"github.com/bebossi/microservice/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID string `json:"id"`
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	 Uow uow.UowInterface,
	 eventDispatcher events.EventDispatcherInterface,
	 transactionCreated events.EventInterface,
	) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow: Uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}


func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.FindByClientID(input.AccountIDFrom)
		if err != nil {
			fmt.Println("err find account from", err)
			return  err
		}
		accountTo, err := accountRepository.FindByClientID(input.AccountIDTo)
		if err != nil {
			fmt.Println("err find account to", err)
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			fmt.Println("err create transaction", err)
			return err
		}
		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			fmt.Println("err update balance", err)
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			fmt.Println("err update balance", err)
			return err
		}
	
		err = transactionRepository.Create(transaction)
		if err != nil {
			fmt.Println("err create transaction", err)
			return err
		}
		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	return output, nil
}


func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}

