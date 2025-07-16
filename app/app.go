package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/handlers"
	"github.com/jrovieri/gobanking/service"
)

func Start() {
	router := mux.NewRouter()
	dbClient := getDBClient()

	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)

	accountHandler := handlers.AccountHandler{Service: service.NewAccountService(accountRepositoryDb)}
	customerHandler := handlers.CustomerHandler{Service: service.NewCustomerService(customerRepositoryDb)}

	router.HandleFunc("/customers", customerHandler.GetAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", customerHandler.GetCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", accountHandler.NewAccount).Methods(http.MethodPost)
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", accountHandler.MakeTransaction).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe("localhost:8000", router))

}

func getDBClient() *sqlx.DB {
	client, err := sqlx.Open("mysql", "gobanking:wks32x23@tcp(127.0.0.1:3306)/gobanking")
	if err != nil {
		panic(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatalln(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
