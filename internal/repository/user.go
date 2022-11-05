package repository

import (
	"context"
	"database/sql"
	"fmt"
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
	fmt.Println(input)
	fmt.Println(user)
	err = txn.QueryRowContext(ctx,
		"UPDATE wallets w1 SET balance=w1.balance+$1 FROM wallets w2 INNER JOIN users on users.wallet = w2.id WHERE users.wallet=$2 AND w1.id =w2.id RETURNING w1.id, w1.balance, users.email",
		input.Amount, input.IDWallet).
		Scan(&user.Wallet.ID, &user.Wallet.Balance, &user.Email)
	fmt.Println(user)
	if err != nil {
		return user, err
	}
	_, err = txn.ExecContext(ctx,
		"INSERT INTO transactions (wallet_id, amount, status, commentary) values ($1, $2, $3, $4)",
		user.Wallet.ID, input.Amount, "approved", "Deposit")
	if err != nil {
		return user, err
	}
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

func (u *Users) CheckAndDoTransfer(ctx context.Context, input wallet.InputTransferUsers) (domain.User, error) {
	var toUser domain.User
	var fromUser domain.User

	txn, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return toUser, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	err = txn.QueryRowContext(ctx, "SELECT wallets.id, balance FROM wallets INNER JOIN users u on wallets.id=u.wallet WHERE u.id=$1", input.FromID).
		Scan(&fromUser.Wallet.ID, &fromUser.Wallet.Balance)
	if err == sql.ErrNoRows {
		return toUser, types.ErrUserFromNotFound
	} else if fromUser.Wallet.Balance < input.Amount {
		return toUser, types.ErrInsufficientFunds
	}
	err = txn.QueryRowContext(ctx, "SELECT email, wallet FROM users WHERE id=$1", input.ToID).Scan(&toUser.Email, &toUser.Wallet.ID)
	if err == sql.ErrNoRows {
		return toUser, types.ErrUserToNotFound
	} else if toUser.Wallet.ID == uuid.Nil {
		// Можно вызвать функцию
		// u.CreateWallet(ctx, wallet.InputDeposit{IDUser: input.ToID, Amount: input.Amount}),
		// но у нее внутри своя транзакция и поэтому думаю, что тогда не получится откатить текущую транзакцию
		err = txn.QueryRowContext(ctx, "INSERT INTO wallets (balance, reserved) values (0, 0) returning id").
			Scan(&toUser.Wallet.ID)
		if err != nil {
			return toUser, err
		}
		_, err = txn.ExecContext(ctx, "UPDATE users SET wallet=$1 WHERE id=$2", toUser.Wallet.ID, input.ToID)
		if err != nil {
			return toUser, err
		}
	}

	err = txn.QueryRowContext(ctx,
		"UPDATE wallets SET balance=wallets.balance-$1 WHERE id=$2 RETURNING balance",
		input.Amount, fromUser.Wallet.ID).Scan(&fromUser.Wallet.Balance)
	if err != nil {
		return toUser, err
	}
	err = txn.QueryRowContext(ctx,
		"UPDATE wallets SET balance=wallets.balance+$1 WHERE id=$2 RETURNING balance",
		input.Amount, toUser.Wallet.ID).Scan(&toUser.Wallet.Balance)

	_, err = txn.ExecContext(ctx, "INSERT INTO transactions (wallet_id, amount, status, commentary) values ($1, $2, $3, $4)",
		fromUser.Wallet.ID, -input.Amount, "approved", "payment send to user")
	_, err = txn.ExecContext(ctx, "INSERT INTO transactions (wallet_id, amount, status, commentary) values ($1, $2, $3, $4)",
		toUser.Wallet.ID, input.Amount, "approved", "payment received from user")
	if err != nil {
		return toUser, err
	}

	return toUser, txn.Commit()
}

func (u *Users) BuyServiceUser(ctx context.Context, input wallet.InputBuyServiceUser) (wallet.OutPendingOrder, error) {
	var user domain.User
	var order wallet.OutPendingOrder

	txn, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return order, err
	}
	defer func() {
		_ = txn.Rollback()
	}()

	err = txn.QueryRowContext(ctx,
		"SELECT wallets.id, balance FROM wallets INNER JOIN users u on wallets.id=u.wallet WHERE u.id=$1",
		input.IDUser).Scan(&user.Wallet.ID, &user.Wallet.Balance)
	if err == sql.ErrNoRows {
		return order, types.ErrUserBuyer
	} else if user.Wallet.Balance < input.Cost {
		return order, types.ErrInsufficientFunds
	}
	err = txn.QueryRowContext(ctx,
		"UPDATE wallets SET balance=wallets.balance-$1, reserved=wallets.reserved+$1 WHERE id=$2 RETURNING balance",
		input.Cost, user.Wallet.ID).Scan(&user.Wallet.Balance)
	if err != nil {
		return order, err
	}
	err = txn.QueryRowContext(ctx, "INSERT INTO transactions (wallet_id, amount, status, commentary) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Wallet.ID, -input.Cost, "pending", fmt.Sprintf("payment for %s", input.ServiceName)).Scan(&order.Txn)

	err = txn.QueryRowContext(ctx,
		"INSERT INTO orders (user_id, service, price, status, transaction_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, status",
		input.IDUser, input.ServiceName, input.Cost, "created", order.Txn).Scan(&order.ID, &order.Status)

	if err != nil {
		if err.Error() == `pq: insert or update on table "orders" violates foreign key constraint "service_fk"` {
			return wallet.OutPendingOrder{}, types.ErrServiceNotFound
		} else {
			return wallet.OutPendingOrder{}, err
		}
	}
	order.Cost = input.Cost
	order.ServiceName = input.ServiceName

	if err != nil {
		return wallet.OutPendingOrder{}, err
	}

	return order, txn.Commit()
}
