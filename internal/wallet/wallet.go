package wallet

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	WalletID      uuid.UUID       `json:"wallet_id" db:"id"`
	OperationType string          `json:"operation_type"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
}
