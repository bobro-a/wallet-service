package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wallet/internal/wallet"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

var ErrorNoWallet = fmt.Errorf("there is no such wallet")

type WalletRepository interface {
	ChangeAmount(ctx context.Context, w wallet.Wallet) (*wallet.Wallet, error)
	GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error)
	CreateWallet(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error)
}

type walletRepo struct {
	db     *sqlx.DB
	nameDB string
}

func NewWalletRepo(db *sqlx.DB, nameDB string) *walletRepo {
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 2)
	return &walletRepo{
		db:     db,
		nameDB: nameDB,
	}
}

func (wr *walletRepo) ChangeAmount(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error) {
	tx, err := wr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	var w wallet.Wallet
	selectQuery := fmt.Sprintf("SELECT id, amount FROM %s WHERE id=$1 FOR UPDATE", wr.nameDB)
	err = tx.GetContext(ctx, &w, selectQuery, req.WalletID)

	if err != nil {
		return nil, err
	}

	switch req.OperationType {
	case "DEPOSIT":
		w.Amount = w.Amount.Add(req.Amount)
	case "WITHDRAW":
		w.Amount = w.Amount.Sub(req.Amount)
	default:
		return nil, errors.New("invalid operation type")
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET amount=$1 WHERE id=$2", wr.nameDB)

	_, err = tx.ExecContext(ctx, updateQuery, w.Amount, w.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet amount: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &w, nil
}

func (wr *walletRepo) GetAmount(ctx context.Context, uuid uuid.UUID) (*wallet.Wallet, error) {
	var w wallet.Wallet
	query := fmt.Sprintf("SELECT id, amount FROM %s WHERE id=$1;", wr.nameDB)

	err := wr.db.QueryRowContext(ctx, query, uuid).Scan(&w.WalletID, &w.Amount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrorNoWallet
	} else if err != nil {
		return nil, err
	}
	return &w, nil
}

func (wr *walletRepo) CreateWallet(ctx context.Context, req wallet.Wallet) (*wallet.Wallet, error) {
	w := &wallet.Wallet{
		WalletID:      req.WalletID,
		OperationType: "",
		Amount:        decimal.NewFromInt(0),
	}
	query := fmt.Sprintf("INSERT INTO %s (id, amount) VALUES ($1, 0)", wr.nameDB)

	_, err := wr.db.ExecContext(ctx, query, w.WalletID)
	if err != nil {
		return nil, fmt.Errorf("couldn't create a new wallet")
	}
	return w, nil
}
