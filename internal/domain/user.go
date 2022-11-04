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

type WebtUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type WebUserWalletBalance struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	Reserved  float64 `json:"reserved"`
}

type Wallet struct {
	Id           uuid.UUID     `json:"id"`
	Balance      float64       `json:"balance"`
	Reserved     float64       `json:"reserved"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

func (WalletUser *WalletUser) ToWebUser() *WebtUser {
	return &WebtUser{
		FirstName: WalletUser.FirstName,
		LastName:  WalletUser.LastName,
		Email:     WalletUser.Email,
	}
}
func (WalletUser *WalletUser) ToWebWalletUser() *WebUserWalletBalance {
	return &WebUserWalletBalance{
		FirstName: WalletUser.FirstName,
		LastName:  WalletUser.LastName,
		Email:     WalletUser.Email,
		Balance:   WalletUser.Wallet.Balance,
		Reserved:  WalletUser.Wallet.Reserved,
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
