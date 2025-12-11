package service

import (
	"context"
	"errors"
	"wallet/internal/wallet"

	"github.com/google/uuid"
)

type MockWalletService struct {
	ChangeAmountFunc func(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error)
	GetAmountFunc    func(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error)
}

func (m *MockWalletService) ChangeAmount(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error) {
	if m.ChangeAmountFunc != nil {
		return m.ChangeAmountFunc(ctx, w)
	}
	return nil, errors.New("ChangeAmount not implemented")
}

func (m *MockWalletService) GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
	if m.GetAmountFunc != nil {
		return m.GetAmountFunc(ctx, uuid)
	}
	return nil, errors.New("GetAmount not implemented")
}
