package repository

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
)

type InvoiceRepository interface {
	List() ([]*model.Invoice, error)
	FindByID(id string) (*model.Invoice, error)
	Create(invoice *model.Invoice) error
	Update(invoice *model.Invoice) error
	Delete(id string) error
}
