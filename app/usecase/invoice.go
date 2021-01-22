package usecase

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
)

type InvoiceUsecase interface{
	List() ([]*model.Invoice, error)
}

type invoiceUsecase struct {
	service *service.InvoiceService
	repo repository.InvoiceRepository
}

func NewInvoiceUsecase(service *service.InvoiceService, repo repository.InvoiceRepository) *invoiceUsecase {
	return &invoiceUsecase{
		service: service,
		repo:    repo,
	}
}

func (u *invoiceUsecase) List() ([]*model.Invoice, error) {
	invoices, err := u.repo.List()
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (u *invoiceUsecase) Create(label string, amount int, receiptDate, dueDate string, from, to string) error {
	// TODO: 'Real' implementation of this function
	return nil
}

