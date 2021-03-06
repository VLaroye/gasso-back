package model

import "time"

type Invoice struct {
	id          string
	amount      int
	label       string
	receiptDate time.Time
	dueDate     time.Time
	from        string
	to          string
}

func NewInvoice(id, label string, amount int, receiptDate, dueDate time.Time, from, to string) *Invoice {
	return &Invoice{
		id:          id,
		amount:      amount,
		label:       label,
		receiptDate: receiptDate,
		dueDate:     dueDate,
		from:        from,
		to:          to,
	}
}

func (invoice *Invoice) GetId() string {
	return invoice.id
}

func (invoice *Invoice) GetAmount() int {
	return invoice.amount
}

func (invoice *Invoice) GetLabel() string {
	return invoice.label
}

func (invoice *Invoice) GetReceiptDate() time.Time {
	return invoice.receiptDate
}

func (invoice *Invoice) GetDueDate() time.Time {
	return invoice.dueDate
}

func (invoice *Invoice) GetTo() string {
	return invoice.to
}

func (invoice *Invoice) GetFrom() string {
	return invoice.from
}
