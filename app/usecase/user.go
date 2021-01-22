package usecase

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
	uuid2 "github.com/google/uuid"
)

type UserUsecase interface {
	ListUsers() ([]*model.User, error)
	RegisterUser(email string) error
}

type userUsecase struct {
	repo    repository.UserRepository
	service *service.UserService
}

func NewUserUsecase(repo repository.UserRepository, service *service.UserService) *userUsecase {
	return &userUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *userUsecase) ListUsers() ([]*model.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userUsecase) RegisterUser(email string) error {
	uuid := uuid2.New()
	if err := u.service.Duplicated(email); err != nil {
		return err
	}

	user := model.NewUser(uuid.String(), email)

	if err := u.repo.Save(user); err != nil {
		return err
	}

	return nil
}
