package database

import (
	"database/sql"
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db *sql.DB
	client *entity.Client
	client2 *entity.Client
	accountFrom *entity.Account
	accountTo *entity.Account
	transactionDB *TransactionDB
}

func (s *TransactionDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id TEXT PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME, updated_at DATETIME)")
	db.Exec("CREATE TABLE accounts (id TEXT PRIMARY KEY, client_id TEXT, balance REAL, created_at DATETIME, updated_at DATETIME)")
	db.Exec("CREATE TABLE transactions (id TEXT PRIMARY KEY, account_id_from TEXT, account_id_to TEXT, amount REAL, created_at DATETIME)")

	client, err := entity.NewClient("John Doe", "john.doe@example.com")
	s.Nil(err)
	s.client = client

	client2, err := entity.NewClient("John Doe 2", "john.doe2@example.com")
	s.Nil(err)
	s.client2 = client2

	accountFrom := entity.NewAccount(client)
	accountFrom.Balance = 1000.0
	s.Nil(err)
	s.accountFrom = accountFrom

	accountTo := entity.NewAccount(client2)
	accountTo.Balance = 1000.0
	s.Nil(err)
	s.accountTo = accountTo

	transactionDB := NewTransactionDB(db)
	s.Nil(err)
	s.transactionDB = transactionDB
	
}

func (s *TransactionDBTestSuite) TearDownTest() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {	
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100.0)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)
}