package service

import (
	"context"
	"errors"
	"fmt"
	"wallet/internal/repo"
	"wallet/internal/wallet"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletService interface {
	GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error)
	ChangeAmount(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error)
}

type service struct {
	repo repo.WalletRepository
}

func NewService(repo repo.WalletRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
	w, err := s.repo.GetAmount(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *service) ChangeAmount(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error) {
	w, err := s.repo.GetAmount(ctx, req.WalletID)
	if errors.Is(err, repo.ErrorNoWallet) && req.OperationType != "WITHDRAW" && req.Amount.GreaterThan(decimal.NewFromInt(0)) {
		newValue, e := s.repo.CreateWallet(ctx, req)
		if e != nil {
			return nil, e
		}
		w = newValue
	} else if err != nil {
		return nil, err
	}
	if req.OperationType == "WITHDRAW" && req.Amount.GreaterThan(w.Amount) {
		return nil, fmt.Errorf("the operation cannot be completed, too little money in the wallet")
	}
	return s.repo.ChangeAmount(ctx, req)
}
