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
	From string `json:"from"`
	To string `json:"to"`
}

func NewInvoice(id, label string, amount int, receiptDate, dueDate time.Time, from, to string) *Invoice {
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
			invoice.GetFrom().GetId(),
			invoice.GetTo().GetId(),
		)
	}
	return result
}

func (service *invoiceService) List(w http.ResponseWriter, r *http.Request) {
	type invoiceResponse struct {
		Invoices []*Invoice `json:"invoices"`
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

func (service *invoiceService) Create(w http.ResponseWriter, r *http.Request) {
	type invoiceRequest struct {
		Label string `json:"label"`
		Amount int `json:"amount"`
		ReceiptDate string `json:"receipt_date"`
		DueDate string `json:"due_date"`
		From string `json:"from"`
		To string `json:"to"'`
	}

	// TODO: Call usecase function, handle errors, respond
}

