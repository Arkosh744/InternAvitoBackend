package service

import (
	"context"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.WalletUser) (domain.WalletUser, error)
	CheckByEmail(ctx context.Context, email string) (domain.WalletUser, error)
}

type Users struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func (u *Users) Create(ctx context.Context, user domain.WalletUser) (domain.WalletUser, error) {
	newUserWallet, err := u.repo.Create(ctx, user)
	if err != nil {
		return newUserWallet, err
	}
	return newUserWallet, err
}

func (u *Users) GetById(ctx context.Context, id int) (domain.WalletUser, error) {
	return domain.WalletUser{}, nil
}

func (u *Users) CheckByEmail(ctx context.Context, email string) (domain.WalletUser, error) {
	checkUser, err := u.repo.CheckByEmail(ctx, email)
	if err != nil {
		return checkUser, err
	}

	return domain.WalletUser{}, nil
}
