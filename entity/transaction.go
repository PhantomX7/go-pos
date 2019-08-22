package entity

import "github.com/PhantomX7/go-pos/models"

type TransactionDetail struct {
	models.Transaction
	Product           models.Product            `json:"product"`
	ReturnTransaction *models.ReturnTransaction `json:"return_transaction,omitempty"`
}
