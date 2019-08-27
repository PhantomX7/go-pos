package entity

import "github.com/PhantomX7/go-pos/models"

type OrderInvoiceDetail struct {
	models.OrderInvoice
	OrderTransaction []OrderTransactionDetail `json:"order_transaction,omitempty"`
}
