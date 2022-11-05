package wallet

import "github.com/google/uuid"

type InputDeposit struct {
	IDUser    uuid.UUID `json:"id_user"`
	IDWallet  uuid.UUID `json:"id_wallet"`
	EmailUser string    `json:"email"`
	Amount    float64   `json:"amount" validate:"required,gte=0"`
}

type InputTransferUsers struct {
	FromID uuid.UUID `json:"from_id" validate:"required"`
	ToID   uuid.UUID `json:"to_id" validate:"required"`
	Amount float64   `json:"amount" validate:"required,gte=0"`
}

type InputBuyServiceUser struct {
	IDUser      uuid.UUID `json:"id_user"`
	ServiceName string    `json:"service_name"`
	Cost        float64   `json:"cost" validate:"required,gte=0"`
}

type OutPendingOrder struct {
	ID          uuid.UUID `json:"id_order"`
	Cost        float64   `json:"cost"`
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	Txn         uuid.UUID `json:"txn"`
}
