package usecase

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/stockmutation"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/service/transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/database"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type TransactionUsecase struct {
	transactionRepo   transaction.TransactionRepository
	stockMutationRepo stockmutation.StockMutationRepository
	invoiceRepo       invoice.InvoiceRepository
	productRepo       product.ProductRepository
}

func NewTransactionUsecase(
	transactionRepo transaction.TransactionRepository,
	stockMutationRepo stockmutation.StockMutationRepository,
	invoiceRepo invoice.InvoiceRepository,
	productRepo product.ProductRepository,
) transaction.TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo:   transactionRepo,
		stockMutationRepo: stockMutationRepo,
		invoiceRepo:       invoiceRepo,
		productRepo:       productRepo,
	}
}

func (t *TransactionUsecase) Create(request request.TransactionCreateRequest) (models.Transaction, error) {
	var transactionM models.Transaction

	// init transaction
	tx := database.BeginTransactions()

	stockMutationM := models.StockMutation{
		ProductID: request.ProductId,
		Amount:    request.Amount,
		Type:      models.StockMutationOUT,
	}

	// create stock mutation
	err := t.stockMutationRepo.Insert(&stockMutationM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	productM, err := t.productRepo.FindByID(request.ProductId)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}
	// reduce the stock
	productM.Stock -= request.Amount
	if productM.Stock < 0 {
		err := errors.ErrUnprocessableEntity
		err.Message = map[string]string{"amount": "invalid amount"}
		return transactionM, err
	}

	err = t.productRepo.Update(&productM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	invoiceM, err := t.invoiceRepo.FindByID(request.InvoiceId)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	transactionM = models.Transaction{
		InvoiceId:       request.InvoiceId,
		ProductId:       request.ProductId,
		StockMutationId: stockMutationM.ID,
		CapitalPrice:    request.CapitalPrice,
		SellPrice:       request.SellPrice,
		Amount:          request.Amount,
		Profit:          (request.SellPrice * request.Amount) - (request.CapitalPrice * request.Amount),
		TotalSellPrice:  request.SellPrice * request.Amount,
		Date:            invoiceM.Date,
	}

	// insert the transaction
	err = t.transactionRepo.Insert(&transactionM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	tx.Commit()

	// update invoice data
	invoice.UpdateInvoice(transactionM.InvoiceId, t.transactionRepo, t.invoiceRepo)

	return transactionM, nil
}

func (t *TransactionUsecase) Update(transactionID uint64, request request.TransactionUpdateRequest) (models.Transaction, error) {
	transactionM, err := t.transactionRepo.FindByID(transactionID)
	if err != nil {
		return transactionM, err
	}

	// init transaction
	tx := database.BeginTransactions()

	// delete previous stock mutation
	err = t.stockMutationRepo.Delete(&models.StockMutation{ID: transactionM.ID}, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	stockMutationM := models.StockMutation{
		ProductID: transactionM.ProductId,
		Amount:    request.Amount,
		Type:      models.StockMutationOUT,
	}

	// create new stock mutation
	err = t.stockMutationRepo.Insert(&stockMutationM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	// re-update product stock
	productM, err := t.productRepo.FindByID(transactionM.ProductId)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	productM.Stock += transactionM.Amount - request.Amount
	if productM.Stock < 0 {
		err := errors.ErrUnprocessableEntity
		err.Message = map[string]string{"amount": "invalid amount"}
		return transactionM, err
	}

	err = t.productRepo.Update(&productM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&transactionM, &request)

	err = t.transactionRepo.Update(&transactionM, tx)
	if err != nil {
		tx.Rollback()
		return transactionM, err
	}

	tx.Commit()

	// update invoice data
	invoice.UpdateInvoice(transactionM.InvoiceId, t.transactionRepo, t.invoiceRepo)

	return transactionM, nil
}

func (t *TransactionUsecase) Delete(transactionID uint64) error {
	transactionM, err := t.transactionRepo.FindByID(transactionID)
	if err != nil {
		return err
	}

	// init transaction
	tx := database.BeginTransactions()

	// delete previous stock mutation
	err = t.stockMutationRepo.Delete(&models.StockMutation{ID: transactionM.ID}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// re-update product stock
	productM, err := t.productRepo.FindByID(transactionM.ProductId)
	if err != nil {
		tx.Rollback()
		return err
	}

	productM.Stock += transactionM.Amount

	err = t.productRepo.Update(&productM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete transaction
	err = t.transactionRepo.Delete(&transactionM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// update invoice data
	invoice.UpdateInvoice(transactionM.InvoiceId, t.transactionRepo, t.invoiceRepo)

	return nil
}

func (t *TransactionUsecase) Index(paginationConfig request.TransactionPaginationConfig) ([]models.Transaction, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	transactions, err := t.transactionRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := t.transactionRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return transactions, meta, nil
}

func (t *TransactionUsecase) Show(transactionID uint64) (models.Transaction, error) {
	return t.transactionRepo.FindByID(transactionID)
}
