package usecase

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
	uuid2 "github.com/google/uuid"
	"time"
)

type InvoiceUsecase interface {
	List() ([]*model.Invoice, error)
	FindById(id string) (*model.Invoice, error)
	Create(label string, amount int, receiptDate, dueDate time.Time, from, to string) error
}

type invoiceUsecase struct {
	service *service.InvoiceService
	repo    repository.InvoiceRepository
}

func NewInvoiceUsecase(repo repository.InvoiceRepository, service *service.InvoiceService) *invoiceUsecase {
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

func (u *invoiceUsecase) FindById(id string) (*model.Invoice, error) {
	invoice, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (u *invoiceUsecase) Create(label string, amount int, receiptDate, dueDate time.Time, from, to string) error {
	uuid := uuid2.New()

	err := u.repo.Create(uuid.String(), label, amount , receiptDate, dueDate, from, to)
	if err != nil {
		return err
	}

	return nil
}
