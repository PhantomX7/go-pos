package entity

import "github.com/PhantomX7/go-pos/models"

type InvoiceDetail struct {
	models.Invoice
	Customer     *models.Customer     `json:"customer,omitempty"`
	Transactions []models.Transaction `json:"transactions"`
}
