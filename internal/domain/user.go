package domain

import "time"

type User struct {
	Id        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Wallet    Wallet    `json:"wallet"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Wallet struct {
	Id           int64         `json:"id"`
	Balance      float64       `json:"balance"`
	Reserved     float64       `json:"reserved"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type Transaction struct {
	Id        int64             `json:"id"`
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
