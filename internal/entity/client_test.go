package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("John Doe", "john.doe@example.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "john.doe@example.com", client.Email)
	assert.NotNil(t, client.ID)
	assert.NotNil(t, client.CreatedAt)
	assert.NotNil(t, client.UpdatedAt)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("John Doe", "john.doe@example.com")
	err := client.Update("John Doe Updated", "john.doe.updated@example.com")
	assert.Nil(t, err)
	assert.Equal(t, "John Doe Updated", client.Name)
	assert.Equal(t, "john.doe.updated@example.com", client.Email)
}

func TestUpdateClientWhenArgsAreInvalid(t *testing.T) {
	client, _ := NewClient("John Doe", "john.doe@example.com")
	err := client.Update("", "")
	assert.Error(t, err, "name is required")
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("John Doe", "john.doe@example.com")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}

func TestAddAccountToClientWhenAccountBelongsToAnotherClient(t *testing.T) {
	client, _ := NewClient("John Doe", "john.doe@example.com")
	account := NewAccount(nil)
	err := client.AddAccount(account)
	assert.Error(t, err, "account does not belong to client")
}