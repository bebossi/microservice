package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bebossi/microservice/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{CreateTransactionUseCase: createTransactionUseCase}
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println("err create transaction", err)
		w.WriteHeader(http.StatusBadRequest)	
		return
	}
	output, err := h.CreateTransactionUseCase.Execute(dto)
	if err != nil {
		fmt.Println("err create transaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		fmt.Println("err create transaction", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
