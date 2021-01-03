package repository

import "github.com/VLaroye/gasso-back/app/domain/model"

type UserRepository interface {
	FindAll() ([]*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Save(user *model.User) error
}
