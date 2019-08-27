package usecase

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/service/order_transaction"
	"github.com/PhantomX7/go-pos/service/order_transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/stock_mutation"
	"github.com/PhantomX7/go-pos/utils/database"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type OrderTransactionUsecase struct {
	orderTransactionRepo order_transaction.OrderTransactionRepository
	stockMutationRepo    stock_mutation.StockMutationRepository
	orderInvoiceRepo     order_invoice.OrderInvoiceRepository
	productRepo          product.ProductRepository
}

func New(
	orderTransactionRepo order_transaction.OrderTransactionRepository,
	stockMutationRepo stock_mutation.StockMutationRepository,
	orderInvoiceRepo order_invoice.OrderInvoiceRepository,
	productRepo product.ProductRepository,
) order_transaction.OrderTransactionUsecase {
	return &OrderTransactionUsecase{
		orderTransactionRepo: orderTransactionRepo,
		stockMutationRepo:    stockMutationRepo,
		orderInvoiceRepo:     orderInvoiceRepo,
		productRepo:          productRepo,
	}
}

func (t *OrderTransactionUsecase) Create(request request.OrderTransactionCreateRequest) (*models.OrderTransaction, error) {
	var orderTransactionM models.OrderTransaction

	// init orderTransaction
	tx := database.BeginTransactions()

	stockMutationM := models.StockMutation{
		ProductID: request.ProductId,
		Amount:    request.Amount,
		Type:      models.StockMutationIN,
	}

	// create stock mutation
	err := t.stockMutationRepo.Insert(&stockMutationM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productM, err := t.productRepo.FindByID(request.ProductId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// add the stock
	productM.Stock += request.Amount

	err = t.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	orderInvoiceM, err := t.orderInvoiceRepo.FindByID(request.OrderInvoiceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	orderTransactionM = models.OrderTransaction{
		OrderInvoiceID:  request.OrderInvoiceId,
		ProductID:       request.ProductId,
		StockMutationID: stockMutationM.ID,
		BuyPrice:        request.BuyPrice,
		Amount:          request.Amount,
		TotalBuyPrice:   request.BuyPrice * request.Amount,
		Date:            orderInvoiceM.Date,
	}

	// insert the orderTransaction
	err = t.orderTransactionRepo.Insert(&orderTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &orderTransactionM, nil
}

func (t *OrderTransactionUsecase) Update(orderTransactionID uint64, request request.OrderTransactionUpdateRequest) (*models.OrderTransaction, error) {
	orderTransactionM, err := t.orderTransactionRepo.FindByID(orderTransactionID)
	if err != nil {
		return orderTransactionM, err
	}

	// init orderTransaction
	tx := database.BeginTransactions()

	if request.Amount != orderTransactionM.Amount {
		// delete previous stock mutation
		err = t.stockMutationRepo.Delete(&models.StockMutation{ID: orderTransactionM.StockMutationID}, tx)
		if err != nil {
			tx.Rollback()
			return orderTransactionM, err
		}

		stockMutationM := models.StockMutation{
			ProductID: orderTransactionM.ProductID,
			Amount:    request.Amount,
			Type:      models.StockMutationIN,
		}

		// create new stock mutation
		err = t.stockMutationRepo.Insert(&stockMutationM, tx)
		if err != nil {
			tx.Rollback()
			return orderTransactionM, err
		}

		// re-update product stock
		productM, err := t.productRepo.FindByID(orderTransactionM.ProductID)
		if err != nil {
			tx.Rollback()
			return orderTransactionM, err
		}

		productM.Stock -= orderTransactionM.Amount - request.Amount

		err = t.productRepo.Update(productM, tx)
		if err != nil {
			tx.Rollback()
			return orderTransactionM, err
		}


		// update stock mutation id
		orderTransactionM.StockMutationID = stockMutationM.ID
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&orderTransactionM, &request)

	err = t.orderTransactionRepo.Update(orderTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return orderTransactionM, err
	}

	tx.Commit()

	return orderTransactionM, nil
}

func (t *OrderTransactionUsecase) Delete(orderTransactionID uint64) error {
	orderTransactionM, err := t.orderTransactionRepo.FindByID(orderTransactionID)
	if err != nil {
		return err
	}

	// init orderTransaction
	tx := database.BeginTransactions()

	// delete previous stock mutation
	err = t.stockMutationRepo.Delete(&models.StockMutation{ID: orderTransactionM.StockMutationID}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// re-update product stock
	productM, err := t.productRepo.FindByID(orderTransactionM.ProductID)
	if err != nil {
		tx.Rollback()
		return err
	}

	productM.Stock -= orderTransactionM.Amount

	err = t.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete orderTransaction
	err = t.orderTransactionRepo.Delete(orderTransactionM, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (t *OrderTransactionUsecase) Index(paginationConfig request.OrderTransactionPaginationConfig) ([]models.OrderTransaction, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	orderTransactions, err := t.orderTransactionRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := t.orderTransactionRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return orderTransactions, meta, nil
}

func (t *OrderTransactionUsecase) Show(orderTransactionID uint64) (*models.OrderTransaction, error) {
	return t.orderTransactionRepo.FindByID(orderTransactionID)
}
