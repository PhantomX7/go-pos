package response

import "github.com/PhantomX7/go-pos/models"

type InvoiceResponse []InvoiceDetail

type InvoiceDetail struct {
	models.Invoice
	Customer *models.Customer `json:"customer,omitempty"`
}
