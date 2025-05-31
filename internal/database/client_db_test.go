package database

import (
	"database/sql"
	"testing"

	"github.com/bebossi/microservice/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id TEXT PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME, updated_at DATETIME)")
	s.clientDB = NewClientDB(db)

}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "john.doe@example.com")
	s.clientDB.Save(client)

	cleintDB, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, cleintDB.ID)
	s.Equal(client.Name, cleintDB.Name)
	s.Equal(client.Email, cleintDB.Email)
}
func (s *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "1",
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}
	err := s.clientDB.Save(client)
	s.Nil(err)
}

