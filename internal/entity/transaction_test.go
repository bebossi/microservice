package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("John Doe", "john.doe@example.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane Doe", "jane.doe@example.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 900, account1.Balance)
	assert.Equal(t, 1100, account2.Balance)
}

func TestCreateTransactionWithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("John Doe", "john.doe@example.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane Doe", "jane.doe@example.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 10000)
	assert.Nil(t, transaction)
	assert.Error(t, err, "account from does not have enough balance")
	assert.Equal(t, 1000, account1.Balance)
	assert.Equal(t, 1000, account2.Balance)
	assert.Error(t, err, "amount must be greater than 0")
}

func TestCreateTransactionWithAmountLessThanZero(t *testing.T) {
	
}
