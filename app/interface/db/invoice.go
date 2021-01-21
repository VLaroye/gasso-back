package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	ID string
	Amount int
	Label string
	ReceiptDate time.Time
	DueDate time.Time
	From *Account
	To *Account
}

func NewInvoice(id, label string, amount int, receiptDate, dueDate time.Time, from, to *Account) *Invoice {
	return &Invoice{
		ID:          id,
		Amount:      amount,
		Label:       label,
		ReceiptDate: receiptDate,
		DueDate:     dueDate,
		From:        from,
		To:          to,
	}
}

type invoiceRepository struct {
	db *gorm.DB
	logger *zap.SugaredLogger
}

func NewInvoiceRepository(db *gorm.DB, logger *zap.SugaredLogger) *invoiceRepository {
	return &invoiceRepository{
		db: db,
		logger: logger,
	}
}

func (ir *invoiceRepository) List() ([]*model.Invoice, error) {
	return nil, nil
}

func (ir *invoiceRepository) FindByID(id string) (*model.Invoice, error) {
	return nil, nil
}

func (ir *invoiceRepository) Create(invoice *model.Invoice) error {
	return nil
}

func (ir *invoiceRepository) Update(invoice *model.Invoice) error {
	return nil
}

func (ir *invoiceRepository) Delete(id string) error {
	return nil
}
