package usecase

import (
	"errors"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
	uuid2 "github.com/google/uuid"
)

type AccountUsecase interface {
	ListAccounts() ([]*model.Account, error)
	GetAccountByID(id string) (*model.Account, error)
	CreateAccount(name string) error
	UpdateAccount(id, name string) error
	DeleteAccount(id string) error
}

type accountUsecase struct {
	repo    repository.AccountRepository
	service *service.AccountService
}

func NewAccountUsecase(repo repository.AccountRepository, service *service.AccountService) *accountUsecase {
	return &accountUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *accountUsecase) ListAccounts() ([]*model.Account, error) {
	accounts, err := u.repo.FindAll()

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (u *accountUsecase) GetAccountByID(id string) (*model.Account, error) {
	account, err := u.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (u *accountUsecase) CreateAccount(name string) error {
	if err := u.service.Duplicated(name); err != nil {
		return err
	}

	uuid := uuid2.New()

	err := u.repo.Create(uuid.String(), name)

	if err != nil {
		return err
	}

	return nil
}

func (u *accountUsecase) UpdateAccount(id, name string) error {
	account, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("account not found")
	}

	if err := u.service.Duplicated(name); err != nil {
		return err
	}

	err = u.repo.Update(id, name)

	if err != nil {
		return err
	}

	return nil
}

func (u *accountUsecase) DeleteAccount(id string) error {
	// Check if account is linked to invoices
	// It it's the case, return an error
	isLinkedToInvoices, err := u.service.IsLinkedToInvoices(id)
	if err != nil {
		return err
	}

	if isLinkedToInvoices {
		return errors.New("account is linked to invoices")
	}

	err = u.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
