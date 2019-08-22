package request

// request related struct

type ReturnTransactionCreateRequest struct {
	TransactionID uint64  `form:"transaction_id" binding:"required,exist=transactions.id,unique=return_transactions.transaction_id"`
	Amount        float64 `form:"amount" binding:"required"`
}

type ReturnTransactionUpdateRequest struct {
	Amount float64 `form:"amount"`
}
