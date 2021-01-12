package repository

import "github.com/VLaroye/gasso-back/app/domain/model"

type AccountRepository interface {
	FindAll() ([]*model.Account, error)
	FindByName(name string) (*model.Account, error)
	FindByID(id string) (*model.Account, error)
	Create(id, name string) error
	Update(id, name string) error
	Delete(id string) error
}
