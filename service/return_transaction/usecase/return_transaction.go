package usecase

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/return_transaction"
	"github.com/PhantomX7/go-pos/service/return_transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/stock_mutation"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/utils/database"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/jinzhu/copier"
)

// apply business logic here

type ReturnTransactionUsecase struct {
	returnTransactionRepo return_transaction.ReturnTransactionRepository
	transactionRepo       transaction.TransactionRepository
	stockMutationRepo     stock_mutation.StockMutationRepository
	productRepo           product.ProductRepository
}

func New(
	returnTransactionRepo return_transaction.ReturnTransactionRepository,
	transactionRepo transaction.TransactionRepository,
	stockMutationRepo stock_mutation.StockMutationRepository,
	productRepo product.ProductRepository,
) return_transaction.ReturnTransactionUsecase {
	return &ReturnTransactionUsecase{
		returnTransactionRepo: returnTransactionRepo,
		transactionRepo:       transactionRepo,
		stockMutationRepo:     stockMutationRepo,
		productRepo:           productRepo,
	}
}

func (t *ReturnTransactionUsecase) Create(request request.ReturnTransactionCreateRequest) (*models.ReturnTransaction, error) {
	var returnTransactionM models.ReturnTransaction

	//get recorded transaction and validate amount
	transactionM, err := t.transactionRepo.FindByID(request.TransactionID)
	if err != nil {
		return nil, err
	}
	if request.Amount > transactionM.Amount {
		err := errors.ErrUnprocessableEntity
		err.Message = map[string]string{"amount": "invalid amount"}
		return nil, err
	}

	// init returnTransaction
	tx := database.BeginTransactions()

	stockMutationM := models.StockMutation{
		ProductID: transactionM.ProductID,
		Amount:    request.Amount,
		Type:      models.StockMutationIN,
	}

	// create stock mutation
	err = t.stockMutationRepo.Insert(&stockMutationM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productM, err := t.productRepo.FindByID(transactionM.ProductID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// increase the stock
	productM.Stock += request.Amount

	err = t.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	returnTransactionM = models.ReturnTransaction{
		TransactionID:   request.TransactionID,
		StockMutationID: stockMutationM.ID,
		InvoiceID:       transactionM.InvoiceID,
		Amount:          request.Amount,
	}

	// insert the returnTransaction
	err = t.returnTransactionRepo.Insert(&returnTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &returnTransactionM, nil
}

func (t *ReturnTransactionUsecase) Update(returnTransactionID uint64, request request.ReturnTransactionUpdateRequest) (*models.ReturnTransaction, error) {
	returnTransactionM, err := t.returnTransactionRepo.FindByID(returnTransactionID)
	if err != nil {
		return returnTransactionM, err
	}

	//get recorded transaction and validate amount
	transactionM, err := t.transactionRepo.FindByID(returnTransactionM.TransactionID)
	if err != nil {
		return returnTransactionM, err
	}
	if request.Amount > transactionM.Amount {
		err := errors.ErrUnprocessableEntity
		err.Message = map[string]string{"amount": "invalid amount"}
		return returnTransactionM, err
	}

	// init returnTransaction
	tx := database.BeginTransactions()

	// delete previous stock mutation
	err = t.stockMutationRepo.Delete(&models.StockMutation{ID: returnTransactionM.StockMutationID}, tx)
	if err != nil {
		tx.Rollback()
		return returnTransactionM, err
	}

	stockMutationM := models.StockMutation{
		ProductID: transactionM.ProductID,
		Amount:    request.Amount,
		Type:      models.StockMutationIN,
	}

	// create new stock mutation
	err = t.stockMutationRepo.Insert(&stockMutationM, tx)
	if err != nil {
		tx.Rollback()
		return returnTransactionM, err
	}

	// re-update product stock
	productM, err := t.productRepo.FindByID(transactionM.ProductID)
	if err != nil {
		tx.Rollback()
		return returnTransactionM, err
	}

	productM.Stock -= returnTransactionM.Amount - request.Amount

	err = t.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return returnTransactionM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&returnTransactionM, &request)
	// update stock mutation id
	returnTransactionM.StockMutationID = stockMutationM.ID

	err = t.returnTransactionRepo.Update(returnTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return returnTransactionM, err
	}

	tx.Commit()

	return returnTransactionM, nil
}

func (t *ReturnTransactionUsecase) Delete(returnTransactionID uint64) error {
	returnTransactionM, err := t.returnTransactionRepo.FindByID(returnTransactionID)
	if err != nil {
		return err
	}

	//get recorded transaction
	transactionM, err := t.transactionRepo.FindByID(returnTransactionM.TransactionID)
	if err != nil {
		return err
	}

	// init returnTransaction
	tx := database.BeginTransactions()

	// delete previous stock mutation
	err = t.stockMutationRepo.Delete(&models.StockMutation{ID: returnTransactionM.StockMutationID}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// re-update product stock
	productM, err := t.productRepo.FindByID(transactionM.ProductID)
	if err != nil {
		tx.Rollback()
		return err
	}

	productM.Stock -= returnTransactionM.Amount

	err = t.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete return transaction
	err = t.returnTransactionRepo.Delete(returnTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (t *ReturnTransactionUsecase) Show(returnTransactionID uint64) (*models.ReturnTransaction, error) {
	return t.returnTransactionRepo.FindByID(returnTransactionID)
}
