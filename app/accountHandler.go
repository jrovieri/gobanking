package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jrovieri/gobanking/dto"
	"github.com/jrovieri/gobanking/service"
)

type AccountHandler struct {
	service service.AccountService
}

func (h AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	customerId := vars["customer_id"]

	var request dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = customerId
		account, err := h.service.NewAccount(request)
		if err != nil {
			writeResponse(w, err.Code, err.Message)
		} else {
			writeResponse(w, http.StatusCreated, account)
		}
	}
}
