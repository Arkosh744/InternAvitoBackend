package service

import (
	"context"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain"
	"github.com/Arkosh744/InternAvitoBackend/internal/domain/wallet"
	"github.com/google/uuid"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.InputUser) (domain.User, error)
	GetUserBalance(ctx context.Context, user domain.User) (domain.User, error)
	CheckUserByEmail(ctx context.Context, email string) (domain.User, error)
	CheckWalletByUserID(ctx context.Context, uuid uuid.UUID) (domain.User, error)
	CheckWalletByEmail(ctx context.Context, user string) (domain.User, error)

	CreateWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)
	DepositWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error)

	CheckAndDoTransfer(ctx context.Context, input wallet.InputTransferUsers) (domain.User, error)
	BuyServiceUser(ctx context.Context, input wallet.InputBuyServiceUser) (wallet.OutPendingOrder, error)
	ManageOrder(ctx context.Context, input wallet.InputOrderManager) (wallet.OutOrderManager, error)

	ReportMonth(ctx context.Context, input wallet.InputReportMonth) ([]wallet.ReportMonth, error)
	ReportForUser(ctx context.Context, input domain.InputReportUserTnx) ([]domain.OutputReportUserTnx, error)
}

type Users struct {
	repo UsersRepository
}

func NewUsersService(repo UsersRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func (u *Users) Create(ctx context.Context, user domain.InputUser) (domain.User, error) {
	newUserWallet, err := u.repo.Create(ctx, user)
	if err != nil {
		return newUserWallet, err
	}
	return newUserWallet, err
}

func (u *Users) GetUserBalance(ctx context.Context, user domain.User) (domain.User, error) {
	userBalance, err := u.repo.GetUserBalance(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return userBalance, nil
}

func (u *Users) CheckUserByEmail(ctx context.Context, email string) (domain.User, error) {
	checkUser, err := u.repo.CheckUserByEmail(ctx, email)
	if err != nil {
		return checkUser, err
	}

	return domain.User{}, nil
}

func (u *Users) CheckWalletByUserID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	checkWalletUser, err := u.repo.CheckWalletByUserID(ctx, uuid)
	if err != nil {
		return domain.User{}, err
	}

	return checkWalletUser, nil
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

func (u *Users) DepositWallet(ctx context.Context, input wallet.InputDeposit) (domain.User, error) {
	userWallet, err := u.repo.DepositWallet(ctx, input)
	if err != nil {
		return domain.User{}, err
	}
	if err != nil {
		return domain.User{}, err
	}
	return userWallet, nil
}

func (u *Users) CheckAndDoTransfer(ctx context.Context, input wallet.InputTransferUsers) (domain.User, error) {
	checkTransfer, err := u.repo.CheckAndDoTransfer(ctx, input)
	if err != nil {
		return domain.User{}, err
	}
	return checkTransfer, nil
}

func (u *Users) BuyServiceUser(ctx context.Context, input wallet.InputBuyServiceUser) (wallet.OutPendingOrder, error) {
	buyServiceUser, err := u.repo.BuyServiceUser(ctx, input)
	if err != nil {
		return wallet.OutPendingOrder{}, err
	}
	return buyServiceUser, nil
}

func (u *Users) ManageOrder(ctx context.Context, input wallet.InputOrderManager) (wallet.OutOrderManager, error) {
	manageOrder, err := u.repo.ManageOrder(ctx, input)
	if err != nil {
		return manageOrder, err
	}
	return manageOrder, nil
}

func (u *Users) ReportForUser(ctx context.Context, input domain.InputReportUserTnx) ([]domain.OutputReportUserTnx, error) {
	reportData, err := u.repo.ReportForUser(ctx, input)
	if err != nil {
		return reportData, err
	}
	return reportData, nil
}

func (u *Users) ReportMonth(ctx context.Context, input wallet.InputReportMonth) ([]wallet.ReportMonth, error) {
	reportData, err := u.repo.ReportMonth(ctx, input)
	if err != nil {
		return reportData, err
	}
	return reportData, nil
}
