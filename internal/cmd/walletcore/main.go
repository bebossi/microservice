package main

import (
	"database/sql"
	"fmt"

	"github.com/bebossi/microservice/internal/database"
	"github.com/bebossi/microservice/internal/event"
	"github.com/bebossi/microservice/internal/usecase/create_account"
	"github.com/bebossi/microservice/internal/usecase/create_client"
	"github.com/bebossi/microservice/internal/usecase/create_transaction"
	"github.com/bebossi/microservice/internal/web"
	"github.com/bebossi/microservice/internal/web/webserver"
	"github.com/bebossi/microservice/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
        "root",           
        "root",           
        "localhost",      
        "3306",          
        "wallet"))        
    if err != nil {
        panic(err)
    }
    defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreated := event.NewTransactionCreated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		transactionDb,
		accountDb,
		eventDispatcher,
		transactionCreated,
	)

	webServer := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webServer.Start()
}
