package repository

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"time"
)

type InvoiceRepository interface {
	List() ([]*model.Invoice, error)
	FindByID(id string) (*model.Invoice, error)
	Create(id, label string, amount int, receiptDate, dueDate time.Time, from, to string) error
	Update(invoice *model.Invoice) error
	Delete(id string) error
}
