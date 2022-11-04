package repository

import (
	"context"
	"database/sql"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (r *Users) Create(ctx context.Context, user domain.WalletUser) (domain.WalletUser, error) {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO wallet (balance) values (0) returning id, created_at, updated_at").
		Scan(&user.Wallet.Id, &user.Wallet.CreatedAt, &user.Wallet.UpdatedAt)
	if err != nil {
		return user, err
	}
	err = r.db.QueryRowContext(ctx,
		"INSERT INTO wallet_user (first_name, last_name, email, wallet) values ($1, $2, $3, $4) "+
			"returning id, created_at, updated_at",
		user.FirstName, user.LastName, user.Email, user.Wallet.Id).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (u *Users) CheckByEmail(ctx context.Context, email string) (domain.WalletUser, error) {
	var user domain.WalletUser
	err := u.db.QueryRowContext(ctx, "SELECT id, email FROM wallet_user WHERE email=$1", email).
		Scan(&user.Id, &user.Email)
	if err != sql.ErrNoRows {
		return user, types.ErrAlreadyExists
	}
	return user, nil
}
