package types

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("datamodel conflict")
	ErrForbidden           = errors.New("forbidden access")
	ErrBadRequest          = errors.New("bad request")
	ErrPartialOk           = errors.New("partial okay")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrGone                = errors.New("resource gone")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrNotAllowed          = errors.New("operation not allowed")
	ErrBusy                = errors.New("resource is busy")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrAlreadyExists       = errors.New("already exists")
	ErrorInternal          = errors.New("internal error")
	ErrNotEnoughData       = errors.New("not enough user data provided")
	ErrTooMuchData         = errors.New("too much user data provided, " +
		"only one of the following is allowed: wallet, user_id, email")
	ErrNoWallet          = errors.New("user does not exist or does not have a wallet")
	ErrSameUser          = errors.New("user cannot send money to himself")
	ErrInsufficientFunds = errors.New("insufficient funds to complete the transaction")
	ErrUserFromNotFound  = errors.New("sender does not exist")
	ErrUserToNotFound    = errors.New("you want to send money to a user that does not exist")
	ErrUserBuyer         = errors.New("buyer does not exist")
	ErrServiceNotFound   = errors.New("service does not exist")
	ErrOrderNotFound     = errors.New("cannot find order by id with given user id")
	ErrOrderCompleted    = errors.New("order is already completed")
)

type HTTPError struct {
	Code    int
	Message string
}

func HTTPCode(err error) int {
	code := 0
	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}
	return code
}
