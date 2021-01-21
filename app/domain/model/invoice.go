package model

import "time"

type Invoice struct {
	id string
	amount int
	label string
	receiptDate time.Time
	dueDate time.Time
	from *Account
	to *Account
}

func NewInvoice(id, label string, amount int, receiptDate, dueDate time.Time, from, to *Account) *Invoice {
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