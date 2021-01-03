package repository

import "github.com/VLaroye/gasso-back/app/domain/model"

type AccountRepository interface {
	FindAll() ([]*model.Account, error)
	FindByName(name string) (*model.Account, error)
	FindById(id string) (*model.Account, error)
	Create(account *model.Account) error
	Update(account *model.Account) error
	Delete(id string) error
}
