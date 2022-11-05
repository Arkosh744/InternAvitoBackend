package repository

import (
	"context"
	"database/sql"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (u *Users) Create(ctx context.Context, user domain.User) (domain.User, error) {
	err := u.db.QueryRowContext(ctx,
		"INSERT INTO users (first_name, last_name, email) values ($1, $2, $3) returning id",
		user.FirstName, user.LastName, user.Email).Scan(&user.ID)
	return user, err
}

func (u *Users) CheckUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := u.db.QueryRowContext(ctx, "SELECT id, email FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email)
	if err != sql.ErrNoRows {
		return user, types.ErrAlreadyExists
	}
	return user, nil
}

func (u *Users) CheckWalletByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := u.db.QueryRowContext(ctx, "SELECT email, wallet FROM users WHERE email=$1", email).
		Scan(&user.Email, &user.Wallet.ID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *Users) CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error) {
	var user domain.User
	err := u.db.QueryRowContext(ctx,
		"INSERT INTO wallets (balance, reserved) values (0, 0) returning id").Scan(&user.Wallet.ID)
	if err != nil {
		return user, err
	}
	err = u.db.QueryRowContext(ctx,
		"UPDATE users SET wallet=$1 WHERE email=$2 returning id, email",
		user.Wallet.ID, input.EmailUser).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}
