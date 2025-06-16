package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bebossi/microservice/internal/database"
	"github.com/bebossi/microservice/internal/event"
	"github.com/bebossi/microservice/internal/event/handler"
	"github.com/bebossi/microservice/internal/usecase/create_account"
	"github.com/bebossi/microservice/internal/usecase/create_client"
	"github.com/bebossi/microservice/internal/usecase/create_transaction"
	"github.com/bebossi/microservice/internal/web"
	"github.com/bebossi/microservice/internal/web/webserver"
	"github.com/bebossi/microservice/pkg/events"
	"github.com/bebossi/microservice/pkg/kafka"
	"github.com/bebossi/microservice/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	_ "github.com/go-sql-driver/mysql"
)


func main() {
    db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
        "root",           
        "root",           
        "mysql",      
        "3306",          
        "wallet"))        
    if err != nil {
        panic(err)
    }
    defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	transactionCreated := event.NewTransactionCreated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	// transactionDb := database.NewTransactionDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		uow,
		eventDispatcher,
		transactionCreated,
	)

	webServer := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is starting on port 8080...")

	webServer.Start()
}
