package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"wallet/internal/repo"
	"wallet/internal/service"
	"wallet/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperationHandler struct {
	service service.WalletService
}

func NewHandler(service service.WalletService) *OperationHandler {
	return &OperationHandler{
		service: service,
	}
}

type Response struct {
	WalletId uuid.UUID
	Amount   decimal.Decimal
}

func (o *OperationHandler) ChangeAmount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("use POST for this request"))
		return
	}
	var wallet wallet.Wallet
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&wallet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := fmt.Sprintf("Failed to decode JSON: %s\n", err.Error())
		w.Write([]byte(errorMessage))
		return
	}
	if wallet.OperationType != "DEPOSIT" && wallet.OperationType != "WITHDRAW" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("the operation_type field must have the value DEPOSIT or WITHDRAW\n"))
		return
	}
	answer, err := o.service.ChangeAmount(r.Context(), wallet)
	if errors.Is(err, repo.ErrorNoWallet) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You need to deposit a positive amount to create a wallet\n"))
		return
	} else if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	response := Response{
		WalletId: answer.WalletID,
		Amount:   answer.Amount,
	}
	jsonData, _ := json.Marshal(response)
	w.Write(jsonData)
}

func (o *OperationHandler) GetAmount(w http.ResponseWriter, r *http.Request) {
	walletUUIDStr := chi.URLParam(r, "wallet_uuid")
	if walletUUIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: Wallet UUID is missing in the path.\n"))
		return
	}
	walletUUID, err := uuid.Parse(walletUUIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid UUID format."))
		return
	}
	answer, err := o.service.GetAmount(r.Context(), walletUUID)
	if errors.Is(err, repo.ErrorNoWallet) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You need to deposit a positive amount to create a wallet\n"))
		return
	} else if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	response := Response{
		WalletId: answer.WalletID,
		Amount:   answer.Amount,
	}
	jsonData, _ := json.Marshal(response)
	w.Write(jsonData)
}
