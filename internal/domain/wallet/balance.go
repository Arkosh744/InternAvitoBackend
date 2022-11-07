package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID         uuid.UUID         `json:"id"`
	WalletId   int64             `json:"walletId"`
	Amount     float64           `json:"cost"`
	Status     TransactionStatus `json:"status"`
	Commentary string            `json:"commentary,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type TransactionStatus struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}
