package usecase

import (
	"github.com/VLaroye/gasso-back/app/domain/repository"
	"github.com/VLaroye/gasso-back/app/domain/service"
)

type InvoiceUsecase interface{}

type invoiceUsecase struct {
	service *service.InvoiceService
	repo repository.InvoiceRepository
}

