package create_transaction

import (
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/bebossi/microservice/internal/event"
	"github.com/bebossi/microservice/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

// UpdateBalance implements gateway.AccountGateway.
func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

// Get implements gateway.AccountGateway.
func (m *AccountGatewayMock) Get(id string) (*entity.Account, error) {
	panic("unimplemented")
}

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByClientID(clientID string) (*entity.Account, error) {
	args := m.Called(clientID)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "john.doe@example.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "jane.doe@example.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockAccount := &AccountGatewayMock{}
	mockAccount.On("FindByClientID", client1.ID).Return(account1, nil)
	mockAccount.On("FindByClientID", client2.ID).Return(account2, nil)

	mockTransaction := &TransactionGatewayMock{}
	mockTransaction.On("Create", mock.Anything).Return(nil)

	inputDto := &CreateTransactionInputDTO{
		AccountIDFrom: client1.ID,
		AccountIDTo:   client2.ID,
		Amount:        100,
	}

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	useCase := NewCreateTransactionUseCase(mockTransaction, mockAccount, dispatcher, event)
	outputDto, err := useCase.Execute(inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, outputDto)
	assert.NotEmpty(t, outputDto.ID)
	mockTransaction.AssertExpectations(t)
	mockAccount.AssertExpectations(t)
}
