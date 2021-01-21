package service

import "github.com/VLaroye/gasso-back/app/domain/repository"

type InvoiceService struct {
	repo repository.InvoiceRepository
}

func NewInvoiceService(repo repository.InvoiceRepository) *InvoiceService {
	return &InvoiceService{
		repo: repo,
	}
}