package domain

import (
	"github.com/google/uuid"
	"time"
)

type WalletUser struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" validate:"required,gte=2"`
	LastName  string    `json:"lastName" validate:"required,gte=2"`
	Email     string    `json:"email" validate:"required,email"`
	Wallet    Wallet    `json:"wallet"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WebWalletUser struct {
	FirstName string    `json:"firstName" validate:"required,gte=2"`
	LastName  string    `json:"lastName" validate:"required,gte=2"`
	Email     string    `json:"email" validate:"required,email"`
	Wallet    WebWallet `json:"wallet"`
}

type Wallet struct {
	Id           uuid.UUID     `json:"id"`
	Balance      float64       `json:"balance"`
	Reserved     float64       `json:"reserved"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type WebWallet struct {
	Balance  float64 `json:"balance"`
	Reserved float64 `json:"reserved"`
}

func (WalletUser *WalletUser) ToWeb() *WebWalletUser {
	return &WebWalletUser{
		FirstName: WalletUser.FirstName,
		LastName:  WalletUser.LastName,
		Email:     WalletUser.Email,
		Wallet: WebWallet{
			Balance:  WalletUser.Wallet.Balance,
			Reserved: WalletUser.Wallet.Reserved,
		},
	}
}

type Transaction struct {
	Id         uuid.UUID         `json:"id"`
	WalletId   int64             `json:"walletId"`
	Amount     float64           `json:"amount"`
	Status     TransactionStatus `json:"status"`
	Commentary string            `json:"commentary,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type TransactionStatus struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}
