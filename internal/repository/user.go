package repository

import (
	"context"
	"database/sql"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/Arkosh744/InternAvitoBackend/pkg/lib/types"
	"github.com/google/uuid"
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

func (u *Users) CheckWalletByUserID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	var user domain.User
	err := u.db.QueryRowContext(ctx, "SELECT id, wallet FROM users WHERE id=$1", uuid).
		Scan(&user.ID, &user.Wallet.ID)
	if err != nil {
		return user, err
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

func (u *Users) DepositWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error) {
	var user domain.User
	txn, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return user, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	err = txn.QueryRowContext(ctx,
		"UPDATE wallets SET balance=wallets.balance+$1 FROM users INNER JOIN wallets w on w.id = users.wallet WHERE wallets.id=$2 RETURNING wallets.id, wallets.balance, email",
		input.Amount, input.IDWallet).
		Scan(&user.Wallet.ID, &user.Wallet.Balance, &user.Email)
	if err != nil {
		return user, err
	}
	err = txn.QueryRowContext(ctx,
		"INSERT INTO transactions (wallet_id, amount, status, commentary) values ($1, $2, $3, $4)",
		user.Wallet.ID, input.Amount, "approved", "Deposit").Scan()
	return user, txn.Commit()
}

func (u *Users) CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error) {
	var user domain.User

	txn, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return user, err
	}
	defer func() {
		_ = txn.Rollback()
	}()
	err = txn.QueryRowContext(ctx, "INSERT INTO wallets (balance, reserved) values (0, 0) returning id").
		Scan(&user.Wallet.ID)
	if err != nil {
		return user, err
	}
	if input.IDUser != uuid.Nil {
		err = txn.QueryRowContext(ctx,
			"UPDATE users SET wallet=$1 WHERE id=$2 returning id", user.Wallet.ID, input.IDUser).Scan(&user.ID)
		if err != nil {
			return user, err
		}
	} else if input.EmailUser != "" {
		err = txn.QueryRowContext(ctx,
			"UPDATE users SET wallet=$1 WHERE email=$2 returning id", user.Wallet.ID, input.EmailUser).Scan(&user.ID)
		if err != nil {
			return user, err
		}
	}
	return user, txn.Commit()
}

func (u *Users) GetUserBalance(ctx context.Context, user domain.User) (domain.User, error) {
	err := u.db.QueryRowContext(ctx, "SELECT email, balance FROM wallets INNER JOIN users u on wallets.id = u.wallet WHERE u.id=$1", user.ID).
		Scan(&user.Email, &user.Wallet.Balance)
	if err == sql.ErrNoRows {
		return user, types.ErrNoWallet
	} else if err != nil {
		return user, err
	}
	return user, nil
}
