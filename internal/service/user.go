package service

import (
	"context"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	CheckUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckWalletByEmail(ctx context.Context, email string) (domain.User, error)
	CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)
}

type Users struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func (u *Users) Create(ctx context.Context, user domain.User) (domain.User, error) {
	newUserWallet, err := u.repo.Create(ctx, user)
	if err != nil {
		return newUserWallet, err
	}
	return newUserWallet, err
}

func (u *Users) GetById(ctx context.Context, id int) (domain.User, error) {
	return domain.User{}, nil
}

func (u *Users) CheckUserByEmail(ctx context.Context, email string) (domain.User, error) {
	checkUser, err := u.repo.CheckUserByEmail(ctx, email)
	if err != nil {
		return checkUser, err
	}

	return domain.User{}, nil
}

func (u *Users) CheckWalletByEmail(ctx context.Context, email string) (domain.User, error) {
	checkWalletUser, err := u.repo.CheckWalletByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return checkWalletUser, nil
}

func (u *Users) CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error) {
	userWallet, err := u.repo.CreateWallet(ctx, input)
	if err != nil {
		return domain.User{}, err
	}
	return userWallet, nil
}
