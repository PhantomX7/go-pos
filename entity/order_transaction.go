package entity

import "github.com/PhantomX7/go-pos/models"

type OrderTransactionDetail struct {
	models.OrderTransaction
	Product models.Product `json:"product"`
}
