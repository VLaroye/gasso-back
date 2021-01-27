package service

import (
	"fmt"
	"github.com/VLaroye/gasso-back/app/domain/repository"
)

type AccountService struct {
	repo repository.AccountRepository
	invoiceRepo repository.InvoiceRepository
}

func NewAccountService(repo repository.AccountRepository, invoiceRepo repository.InvoiceRepository) *AccountService {
	return &AccountService{
		repo: repo,
		invoiceRepo: invoiceRepo,
	}
}

func (s *AccountService) Duplicated(name string) error {
	account, err := s.repo.FindByName(name)

	if account != nil {
		return fmt.Errorf("%s already exists", name)
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) IsLinkedToInvoices(accountId string) (bool, error) {
	invoices, err := s.invoiceRepo.ListByAccount(accountId)
	if err != nil {
		return false, err
	}

	return invoices != nil, nil
}
