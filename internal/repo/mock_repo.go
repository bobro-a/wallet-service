package repo

import (
	"context"
	"errors"
	"wallet/internal/wallet"

	"github.com/google/uuid"
)

type MockWalletRepo struct {
	ChangeAmountFunc func(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error)
	GetAmountFunc    func(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error)
	CreateWalletFunc func(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error)
}

func (m *MockWalletRepo) ChangeAmount(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error) {
	if m.ChangeAmountFunc != nil {
		return m.ChangeAmountFunc(ctx, w)
	}
	return nil, errors.New("ChangeAmount not implemented")
}

func (m *MockWalletRepo) GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
	if m.GetAmountFunc != nil {
		return m.GetAmountFunc(ctx, uuid)
	}
	return nil, errors.New("GetAmount not implemented")
}

func (m *MockWalletRepo) CreateWallet(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error) {
	if m.CreateWalletFunc != nil {
		return m.CreateWalletFunc(ctx, req)
	}
	return nil, errors.New("CreateWallet not implemented")
}
