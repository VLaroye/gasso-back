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
	From *Account `gorm:"foreignkey:ID"`
	To *Account `gorm:"foreignkey:ID"`
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
	var invoices []*Invoice

	result := ir.db.Find(&invoices)

	if result.Error != nil {
		ir.logger.Errorw("error fetching list invoices from db",
			"error", result.Error,
		)
		return nil, result.Error
	}

	ir.logger.Infow("list invoices fetched from db",
		"nb of invoices fetched", result.RowsAffected,
	)

	response := make([]*model.Invoice, len(invoices))

	for i, invoice := range invoices {
		response[i] = model.NewInvoice(
			invoice.ID,
			invoice.Label,
			invoice.Amount,
			invoice.ReceiptDate,
			invoice.DueDate,
			model.NewAccount(invoice.From.ID, invoice.From.Name),
			model.NewAccount(invoice.To.ID, invoice.To.Name),
		)
	}

	return response, nil
}

func (ir *invoiceRepository) FindByID(id string) (*model.Invoice, error) {
	return nil, nil
}

func (ir *invoiceRepository) Create(id, label string, amount int, receiptDate, dueDate time.Time, from, to string) error {
	// TODO: Think about receive format here for dates (string ? Time ? ...), implement call to db
	return nil
}

func (ir *invoiceRepository) Update(invoice *model.Invoice) error {
	return nil
}

func (ir *invoiceRepository) Delete(id string) error {
	return nil
}
