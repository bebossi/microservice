package create_client

import (
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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


func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)

	useCase := NewCreateClientUseCase(m)

	output, err := useCase.Execute(CreateClientInputDTO{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john.doe@example.com", output.Email)
	m.AssertExpectations(t)
}
