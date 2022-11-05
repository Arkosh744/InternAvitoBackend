package types

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("datamodel conflict")
	ErrForbidden           = errors.New("forbidden access")
	ErrNeedMore            = errors.New("need more input")
	ErrBadRequest          = errors.New("bad request")
	ErrPartialOk           = errors.New("partial okay")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrGone                = errors.New("resource gone")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrNotAllowed          = errors.New("operation not allowed")
	ErrBusy                = errors.New("resource is busy")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotEnoughData       = errors.New("not enough user data provided")
	ErrTooMuchData         = errors.New("too much user data provided, " +
		"only one of the following is allowed: wallet, user_id, email")
	ErrNoWallet = errors.New("user does not exist or does not have a wallet")
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
