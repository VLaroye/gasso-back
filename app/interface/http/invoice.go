package http

import (
	"net/http"
	"time"

	"github.com/VLaroye/gasso-back/app/domain/model"
	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterInvoiceHandlers(router *mux.Router, service *invoiceService) {
	router.HandleFunc("/invoices", service.List).Methods("GET")
}

type Invoice struct {
	ID string `json:"id"`
	Amount int `json:"amount"`
	Label string `json:"label"`
	ReceiptDate time.Time `json:"receipt_date"`
	DueDate time.Time `json:"due_date"`
	From *Account `json:"from"`
	To *Account `json:"to"`
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

type invoiceService struct {
	usecase usecase.InvoiceUsecase
}

func NewInvoiceService(usecase usecase.InvoiceUsecase) *invoiceService {
	return &invoiceService{usecase: usecase}
}

func toInvoices(invoices []*model.Invoice) []*Invoice {
	result := make([]*Invoice, len(invoices))

	for i, invoice := range invoices {
		result[i] = NewInvoice(
			invoice.GetId(),
			invoice.GetLabel(),
			invoice.GetAmount(),
			invoice.GetReceiptDate(),
			invoice.GetDueDate(),
			NewAccount(invoice.GetFrom().GetId(), invoice.GetFrom().GetName()),
			NewAccount(invoice.GetTo().GetId(), invoice.GetTo().GetName()),
		)
	}
	return result
}

func (service *invoiceService) List(w http.ResponseWriter, r *http.Request) {
	type invoiceResponse struct {
		Invoices []*Invoice `json:"accounts"`
	}

	invoices, err := service.usecase.List()

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := invoiceResponse{
		Invoices: toInvoices(invoices),
	}

	respondJSON(w, response)
	return
}

