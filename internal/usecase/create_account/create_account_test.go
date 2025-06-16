package create_account

import (
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(clientID string) (*entity.Account, error) {
	args := m.Called(clientID)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Get(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "john.doe@example.com")
	clientMock := &ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)

	accountMock := &AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)

	useCase := NewCreateAccountUseCase(accountMock, clientMock)
	inputDto := &CreateAccountInputDTO{
		ClientID: client.ID,
	}
	outputDto, err := useCase.Execute(*inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, outputDto)
	assert.NotEmpty(t, outputDto.ID)
	clientMock.AssertExpectations(t)
	accountMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
