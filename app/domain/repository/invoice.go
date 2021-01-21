package repository

import "github.com/VLaroye/gasso-back/app/domain/model"

type InvoiceRepository interface {
	FindByID(id string) (*model.Invoice)
}
