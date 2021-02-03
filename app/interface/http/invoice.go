package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VLaroye/gasso-back/app/usecase"
	"github.com/gorilla/mux"
)

func RegisterInvoiceHandlers(router *mux.Router, service *invoiceService) {
	router.HandleFunc("/invoices", AuthenticationNeeded(service.List)).Methods("GET")
	router.HandleFunc("/invoices", AuthenticationNeeded(service.Create)).Methods("POST")
}

type Invoice struct {
	ID          string    `json:"id"`
	Amount      int       `json:"amount"`
	Label       string    `json:"label"`
	ReceiptDate time.Time `json:"receipt_date"`
	DueDate     time.Time `json:"due_date"`
	From        *Account  `json:"from"`
	To          *Account  `json:"to"`
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
	usecase        usecase.InvoiceUsecase
	accountUsecase usecase.AccountUsecase
}

func NewInvoiceService(usecase usecase.InvoiceUsecase, accountUsecase usecase.AccountUsecase) *invoiceService {
	return &invoiceService{usecase: usecase, accountUsecase: accountUsecase}
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

	response := make([]*Invoice, len(invoices))

	for i, invoice := range invoices {
		from, err := service.accountUsecase.GetAccountByID(invoice.GetFrom())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		to, err := service.accountUsecase.GetAccountByID(invoice.GetTo())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if from == nil || to == nil {
			respondError(w, http.StatusInternalServerError, "invoice's 'to' or 'from' field is invalid")
			return
		}

		response[i] = NewInvoice(
			invoice.GetId(),
			invoice.GetLabel(),
			invoice.GetAmount(),
			invoice.GetReceiptDate(),
			invoice.GetDueDate(),
			NewAccount(from.GetId(), from.GetName()),
			NewAccount(to.GetId(), to.GetName()),
		)
	}

	respondJSON(w, response)
	return
}

func (service *invoiceService) Create(w http.ResponseWriter, r *http.Request) {
	type invoiceRequest struct {
		Label       string `json:"label"`
		Amount      int    `json:"amount"`
		ReceiptDate string `json:"receipt_date"`
		DueDate     string `json:"due_date"`
		From        string `json:"from"`
		To          string `json:"to"'`
	}

	var request invoiceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if accounts exist
	// TODO: Should be done on invoice usecase/service ?
	from, err := service.accountUsecase.GetAccountByID(request.From)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if from == nil {
		respondError(w, http.StatusBadRequest, "can't find 'from' account")
		return
	}

	to, err := service.accountUsecase.GetAccountByID(request.To)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if to == nil {
		respondError(w, http.StatusBadRequest, "can't find 'to' account")
		return
	}

	// Check if dates are valid
	receiptDate, err := time.Parse(time.RFC3339, request.ReceiptDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid 'receipt_date'")
		return
	}
	dueDate, err := time.Parse(time.RFC3339, request.DueDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid 'due_date'")
		return
	}

	// Create invoice
	if err := service.usecase.Create(
		request.Label,
		request.Amount,
		receiptDate,
		dueDate,
		from.GetId(),
		to.GetId(),
	); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
