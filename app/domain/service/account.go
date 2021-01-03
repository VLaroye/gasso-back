package service

import (
	"fmt"

	"github.com/VLaroye/gasso-back/app/domain/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
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
