package wallet

import "github.com/google/uuid"

type InputDeposit struct {
	IDUser    uuid.UUID `json:"id_user"`
	IDWallet  uuid.UUID `json:"id_wallet"`
	EmailUser string    `json:"email" validate:"email"`
	Amount    float64   `json:"amount" validate:"required,gte=0"`
}

type InputTransferUsers struct {
	FromID uuid.UUID `json:"from_id" validate:"required"`
	ToID   uuid.UUID `json:"to_id" validate:"required"`
	Amount float64   `json:"amount" validate:"required,gte=0"`
}
