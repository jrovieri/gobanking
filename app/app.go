package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/service"
)

func Start() {

	router := mux.NewRouter()

	customerHandler := CustomerHandler{service: service.NewCustomerService(domain.NewCustomerRepositoryDb())}
	router.HandleFunc("/customers", customerHandler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", customerHandler.getCustomer).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe("localhost:8000", router))

}
