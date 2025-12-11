package service

import (
	"context"
	"errors"
	"testing"
	"wallet/internal/repo"
	"wallet/internal/wallet"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var testUUID = uuid.MustParse("b1b2c3d4-e5f6-7890-1234-567890abcdef")

func TestChangeAmount_DepositSuccess(t *testing.T) {
	initialAmount := decimal.NewFromInt(100)
	depositAmount := decimal.NewFromInt(50)
	expectedAmount := initialAmount.Add(depositAmount)

	mockRepo := &repo.MockWalletRepo{
		GetAmountFunc: func(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
			return &wallet.Wallet{
				WalletID: uuid,
				Amount:   initialAmount,
			}, nil
		},

		ChangeAmountFunc: func(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error) {
			w.Amount = expectedAmount
			return &w, nil
		},
	}

	svc := NewService(mockRepo)

	req := wallet.Wallet{
		WalletID:      testUUID,
		OperationType: "DEPOSIT",
		Amount:        depositAmount,
	}

	result, err := svc.ChangeAmount(context.Background(), req)

	if err != nil {
		t.Fatalf("Ожидалась ошибка: нет, получена: %v", err)
	}

	if !result.Amount.Equal(expectedAmount) {
		t.Errorf("Ожидаемый баланс: %s, получен: %s", expectedAmount.String(), result.Amount.String())
	}
}

func TestChangeAmount_WithdrawInsufficientBalance(t *testing.T) {

	withdrawAmount := decimal.NewFromInt(150)

	mockRepo := &repo.MockWalletRepo{
		ChangeAmountFunc: func(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error) {
			return nil, errors.New("Insufficient balance in database")
		},
	}

	svc := NewService(mockRepo)

	req := wallet.Wallet{
		WalletID:      testUUID,
		OperationType: "WITHDRAW",
		Amount:        withdrawAmount,
	}

	_, err := svc.ChangeAmount(context.Background(), req)

	if err == nil {
		t.Fatal("Ожидалась ошибка недостаточного баланса, но ошибок не получено")
	}
}

func TestGetAmount_WalletNotFound(t *testing.T) {
	mockRepo := &repo.MockWalletRepo{
		GetAmountFunc: func(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
			return nil, repo.ErrorNoWallet
		},
	}

	svc := NewService(mockRepo)

	_, err := svc.GetAmount(context.Background(), testUUID)

	if !errors.Is(err, repo.ErrorNoWallet) {
		t.Errorf("Ожидалась ошибка 'нет кошелька', получена другая ошибка: %v", err)
	}
}
