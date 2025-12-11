package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/internal/repo"
	"wallet/internal/service"
	"wallet/internal/wallet"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var testUUID = uuid.MustParse("b1b2c3d4-e5f6-7890-1234-567890abcdef")

func TestChangeAmount_Success(t *testing.T) {
	expectedAmount := decimal.NewFromInt(150)

	mockSvc := &service.MockWalletService{
		ChangeAmountFunc: func(_ context.Context, _ wallet.Wallet) (*wallet.Wallet, error) {
			return &wallet.Wallet{WalletID: testUUID, Amount: expectedAmount}, nil
		},
	}
	handler := NewHandler(mockSvc)

	requestBody := []byte(`{"wallet_id": "b1b2c3d4-e5f6-7890-1234-567890abcdef", "operation_type": "DEPOSIT", "amount": 50.00}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ChangeAmount(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v. Тело ответа: %s", http.StatusOK, status, rr.Body.String())
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Не удалось декодировать ответ JSON: %v", err)
	}
	if !resp.Amount.Equal(expectedAmount) {
		t.Errorf("Ожидаемая сумма %s, получена %s", expectedAmount.String(), resp.Amount.String())
	}
}

func TestChangeAmount_InvalidOperationType(t *testing.T) {
	handler := NewHandler(nil)

	requestBody := []byte(`{"wallet_id": "b1b2c3d4-e5f6-7890-1234-567890abcdef", "operation_type": "TRANSFER", "amount": 50.00}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(requestBody))
	rr := httptest.NewRecorder()

	handler.ChangeAmount(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusBadRequest, status)
	}
	expectedBody := "the operation_type field must have the value DEPOSIT or WITHDRAW\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Неверное сообщение об ошибке. Ожидалось: %s, Получено: %s", expectedBody, rr.Body.String())
	}
}

func TestGetAmount_Success(t *testing.T) {
	expectedAmount := decimal.NewFromInt(150)

	mockSvc := &service.MockWalletService{
		GetAmountFunc: func(_ context.Context, _ uuid.UUID) (*wallet.Wallet, error) {
			return &wallet.Wallet{WalletID: testUUID, Amount: expectedAmount}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+testUUID.String(), nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("wallet_uuid", testUUID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetAmount(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус %v, получен %v. Тело: %s", http.StatusOK, status, rr.Body.String())
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Не удалось декодировать ответ JSON: %v", err)
	}
	if !resp.Amount.Equal(expectedAmount) {
		t.Errorf("Ожидаемая сумма %s, получена %s", expectedAmount.String(), resp.Amount.String())
	}
}

func TestGetAmount_WalletNotFound(t *testing.T) {
	mockSvc := &service.MockWalletService{
		GetAmountFunc: func(_ context.Context, _ uuid.UUID) (*wallet.Wallet, error) {
			return nil, repo.ErrorNoWallet
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+testUUID.String(), nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("wallet_uuid", testUUID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetAmount(rr, req)

	if status := rr.Code; status != http.StatusOK { // Обратите внимание: ваш код возвращает 200 OK при ErrorNoWallet
		t.Errorf("Ожидался статус %v, получен %v", http.StatusOK, status)
	}
	expectedBody := "You need to deposit a positive amount to create a wallet\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Неверное сообщение об ошибке. Ожидалось: %s, Получено: %s", expectedBody, rr.Body.String())
	}
}

func TestGetAmount_InvalidUUID(t *testing.T) {
	handler := NewHandler(nil)

	invalidUUID := "not-a-valid-uuid"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+invalidUUID, nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("wallet_uuid", invalidUUID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetAmount(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидался статус %v, получен %v", http.StatusBadRequest, status)
	}
	expectedBody := "Invalid UUID format."
	if rr.Body.String() != expectedBody {
		t.Errorf("Неверное сообщение об ошибке. Ожидалось: %s, Получено: %s", expectedBody, rr.Body.String())
	}
}
