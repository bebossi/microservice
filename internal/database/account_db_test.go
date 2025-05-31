package database

import (
	"database/sql"
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db *sql.DB
	accountDB *AccountDB
	client *entity.Client
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id TEXT PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME, updated_at DATETIME)")
	db.Exec("CREATE TABLE accounts (id TEXT PRIMARY KEY, client_id TEXT, balance REAL, created_at DATETIME, updated_at DATETIME)")
	s.accountDB = NewAccountDB(db)
	s.client, err = entity.NewClient("John Doe", "john.doe@example.com")
	s.Nil(err)
}


func (s *AccountDBTestSuite) TearDownTest() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite (t *testing.T){
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestFindByClientID() {
    _, err := s.db.Exec("INSERT INTO clients (id, name, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt, s.client.UpdatedAt)
    s.NoError(err)
    
    _, err = s.db.Exec("INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        "account-id-1", s.client.ID, 100.0, s.client.CreatedAt, s.client.UpdatedAt)
    s.NoError(err)
    
    account, err := s.accountDB.FindByClientID(s.client.ID)
    s.NoError(err)
    s.NotNil(account)
    s.Equal(s.client.ID, account.Client.ID)
    s.Equal(100.0, account.Balance)
}

func (s *AccountDBTestSuite) TestSave() {

}
