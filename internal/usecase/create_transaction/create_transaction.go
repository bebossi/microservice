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

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated events.EventInterface
}

func NewCreateTransactionUseCase(
	 Uow uow.UowInterface,
	 eventDispatcher events.EventDispatcherInterface,
	 transactionCreated events.EventInterface,
	 balanceUpdated events.EventInterface,
	) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow: Uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated: balanceUpdated,
	}
}


func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	balanceUpdated := &BalanceUpdatedOutputDTO{}
	fmt.Printf("Debug: Starting transaction creation with input: %+v\n", input)
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)
        fmt.Printf("Debug: Finding account from: %s\n", input.AccountIDFrom)

		accountFrom, err := accountRepository.FindByID(input.AccountIDFrom)
		if err != nil {
			fmt.Println("err find account from", err)
			return  err
		}
		fmt.Printf("Debug: Found account from: %+v\n", accountFrom)

		accountTo, err := accountRepository.FindByID(input.AccountIDTo)
		fmt.Printf("Debug: Finding account to: %s\n", input.AccountIDTo)

		if err != nil {
			fmt.Println("err find account to", err)
			return err
		}
		fmt.Printf("Debug: Found account to: %+v\n", accountTo)

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

		balanceUpdated.AccountIDFrom = input.AccountIDFrom
		balanceUpdated.AccountIDTo = input.AccountIDTo
		balanceUpdated.BalanceAccountIDFrom = accountFrom.Balance
		balanceUpdated.BalanceAccountIDTo = accountTo.Balance
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)
	uc.BalanceUpdated.SetPayload(balanceUpdated)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)

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

