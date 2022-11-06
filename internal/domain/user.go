package domain

import (
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" validate:"required,gte=2"`
	LastName  string    `json:"lastName" validate:"required,gte=2"`
	Email     string    `json:"email" validate:"required,email"`
	Wallet    Wallet    `json:"wallet"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WebtUser struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

type WebUserWalletBalance struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Balance   float64 `json:"balance"`
	Reserved  float64 `json:"reserved"`
}

type Wallet struct {
	ID           uuid.UUID            `json:"id"`
	Balance      float64              `json:"balance"`
	Reserved     float64              `json:"reserved"`
	Transactions []wallet.Transaction `json:"transactions,omitempty"`
	CreatedAt    time.Time            `json:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
}

func (User *User) ToWebUser() *WebtUser {
	return &WebtUser{
		ID:        User.ID,
		FirstName: User.FirstName,
		LastName:  User.LastName,
		Email:     User.Email,
	}
}
func (User *User) ToWebWalletUser() *WebUserWalletBalance {
	return &WebUserWalletBalance{
		FirstName: User.FirstName,
		LastName:  User.LastName,
		Email:     User.Email,
		Balance:   User.Wallet.Balance,
		Reserved:  User.Wallet.Reserved,
	}
}

type InputReportUserTnx struct {
	IDUser    uuid.UUID `json:"user_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Order     string    `json:"order"`
	SortField string    `json:"sort_field"`
}

type OutputReportUserTnx struct {
	Date       time.Time
	Commentary string
	Amount     float64
}
