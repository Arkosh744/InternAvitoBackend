package wallet

import "github.com/google/uuid"

type InputDeposit struct {
	IDUser    uuid.UUID `json:"id_user"`
	IDWallet  uuid.UUID `json:"id_wallet"`
	EmailUser string    `json:"email"`
	Amount    float64   `json:"amount" validate:"required,gte=0"`
}
