package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Invoice struct {
	ID          string
	Amount      int
	Label       string
	ReceiptDate time.Time
	DueDate     time.Time
	FromId string
	From        *Account `gorm:"foreignkey:FromId"`
	ToId string
	To          *Account `gorm:"foreignkey:ToId"`
}

func NewInvoice(id, label string, amount int, receiptDate, dueDate time.Time, fromId, toId string) *Invoice {
	return &Invoice{
		ID:          id,
		Amount:      amount,
		Label:       label,
		ReceiptDate: receiptDate,
		DueDate:     dueDate,
		FromId: fromId,
		ToId: toId,
	}
}

type invoiceRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewInvoiceRepository(db *gorm.DB, logger *zap.SugaredLogger) *invoiceRepository {
	return &invoiceRepository{
		db:     db,
		logger: logger,
	}
}

func (ir *invoiceRepository) List() ([]*model.Invoice, error) {
	var invoices []*Invoice

	result := ir.db.Preload(clause.Associations).Find(&invoices)

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
	var invoice Invoice
	result := ir.db.Where("id = ?", id).Find(&invoice)

	if result.Error != nil {
		ir.logger.Errorw("find invoice by id failed",
			"id", id,
		)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		ir.logger.Infow("find invoice by id failed, invoice not found",
			"id", id,
		)
		return nil, nil
	}

	ir.logger.Infow("invoice fetched by id",
		"id", id,
	)

	return model.NewInvoice(
		invoice.ID,
		invoice.Label,
		invoice.Amount,
		invoice.ReceiptDate,
		invoice.DueDate,
		model.NewAccount(invoice.From.ID, invoice.From.Name),
		model.NewAccount(invoice.To.ID, invoice.To.Name),
	), nil
}

func (ir *invoiceRepository) Create(id, label string, amount int, receiptDate, dueDate time.Time, from, to *model.Account) error {
	invoiceToInsert := NewInvoice(
		id,
		label,
		amount,
		receiptDate,
		dueDate,
		from.GetId(),
		to.GetId(),
	)

	result := ir.db.Create(&invoiceToInsert)

	if result.Error != nil {
		ir.logger.Errorw("create invoice failed",
			"invoiceId", invoiceToInsert.ID,
			"error", result.Error,
		)
		return result.Error
	}

	ir.logger.Infow("invoice created", "invoice", invoiceToInsert)

	return nil
}

func (ir *invoiceRepository) Update(invoice *model.Invoice) error {
	return nil
}

func (ir *invoiceRepository) Delete(id string) error {
	return nil
}
