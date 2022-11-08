package wallet

import "github.com/google/uuid"

type InputDeposit struct {
	IDUser    uuid.UUID `json:"id_user,omitempty" default:""`
	IDWallet  uuid.UUID `json:"id_wallet,omitempty" default:""`
	EmailUser string    `json:"email,omitempty" validate:"omitempty,email" default:""`
	Amount    float64   `json:"amount" validate:"required,gte=0"`
}

type InputTransferUsers struct {
	FromID uuid.UUID `json:"from_id" validate:"required"`
	ToID   uuid.UUID `json:"to_id" validate:"required"`
	Amount float64   `json:"amount" validate:"required,gte=0"`
}

type InputBuyServiceUser struct {
	IDUser      uuid.UUID `json:"id_user" validate:"required"`
	ServiceName string    `json:"service_name" validate:"required"`
	Cost        float64   `json:"cost" validate:"required,gte=0"`
}

type OutPendingOrder struct {
	ID          uuid.UUID `json:"id_order"`
	Cost        float64   `json:"cost"`
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	Txn         uuid.UUID `json:"txn"`
}

type InputOrderManager struct {
	IDOrder uuid.UUID `json:"id_order" validate:"required"`
	IDUser  uuid.UUID `json:"id_user" validate:"required"`
	Status  string    `json:"status,omitempty"`
}

type OutOrderManager struct {
	IDOrder     uuid.UUID `json:"id_order"`
	Cost        float64   `json:"cost"`
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	TxnSeller   uuid.UUID `json:"txn_seller"`
	TxnBuyer    uuid.UUID `json:"txn_buyer"`
}

type InputReportMonth struct {
	Year  int `json:"year" validate:"required,gte=2020"`
	Month int `json:"month" validate:"required,gte=1,lte=12"`
}

type ReportMonth struct {
	Amount int
	Text   string
}
