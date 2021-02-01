package usecase

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
	uuid2 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Login(email, password string) error
	RegisterUser(email, password string) error
	ListUsers() ([]*model.User, error)
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

func (u *userUsecase) Login(email, password string) error {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) RegisterUser(email, password string) error {
	uuid := uuid2.New()
	if err := u.service.Duplicated(email); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.NewUser(uuid.String(), email, string(hashedPassword))

	if err := u.repo.Save(user); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) ListUsers() ([]*model.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
