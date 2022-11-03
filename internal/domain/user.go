package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email" validate:"required"`
	Wallet    Wallet    `json:"wallet"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Wallet struct {
	Id           uuid.UUID     `json:"id"`
	Balance      float64       `json:"balance"`
	Reserved     float64       `json:"reserved"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type Transaction struct {
	Id        uuid.UUID         `json:"id"`
	WalletId  int64             `json:"walletId"`
	Amount    float64           `json:"amount"`
	Status    TransactionStatus `json:"status"`
	Comment   string            `json:"comment,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type TransactionStatus struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}
