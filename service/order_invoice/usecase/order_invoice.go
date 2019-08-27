package usecase

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/service/order_invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/order_transaction"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type OrderInvoiceUsecase struct {
	orderInvoiceRepo     order_invoice.OrderInvoiceRepository
	orderTransactionRepo order_transaction.OrderTransactionRepository
	productRepo          product.ProductRepository
}

func New(
	orderInvoiceRepo order_invoice.OrderInvoiceRepository,
	orderTransactionRepo order_transaction.OrderTransactionRepository,
	productRepo product.ProductRepository,
) order_invoice.OrderInvoiceUsecase {
	return &OrderInvoiceUsecase{
		orderInvoiceRepo:     orderInvoiceRepo,
		orderTransactionRepo: orderTransactionRepo,
		productRepo:          productRepo,
	}
}

func (a *OrderInvoiceUsecase) Create(request request.OrderInvoiceCreateRequest) (*models.OrderInvoice, error) {
	orderInvoiceM := models.OrderInvoice{
		Description:   request.Description,
		Date:          request.Date,
		PaymentStatus: request.PaymentStatus,
		TotalBuyPrice: 0,
	}
	err := a.orderInvoiceRepo.Insert(&orderInvoiceM, nil)
	if err != nil {
		return nil, err
	}
	return &orderInvoiceM, nil
}

func (a *OrderInvoiceUsecase) Update(orderInvoiceID uint64, request request.OrderInvoiceUpdateRequest) (*models.OrderInvoice, error) {
	orderInvoiceM, err := a.orderInvoiceRepo.FindByID(orderInvoiceID)
	if err != nil {
		return orderInvoiceM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(orderInvoiceM, &request)

	err = a.orderInvoiceRepo.Update(orderInvoiceM, nil)
	if err != nil {
		return orderInvoiceM, err
	}
	return orderInvoiceM, nil
}

func (a *OrderInvoiceUsecase) Delete(orderInvoiceID uint64) error {
	orderTransactions, err := a.orderTransactionRepo.FindByOrderInvoiceID(orderInvoiceID)
	if err != nil {
		return err
	}
	if len(orderTransactions) > 0 {
		err := errors.ErrUnprocessableEntity
		err.Message = "please delete all order transaction first"
		return err
	}
	err = a.orderInvoiceRepo.Delete(&models.OrderInvoice{ID: orderInvoiceID}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *OrderInvoiceUsecase) Index(paginationConfig request.OrderInvoicePaginationConfig) ([]entity.OrderInvoiceDetail, response_util.PaginationMeta, error) {
	res := []entity.OrderInvoiceDetail{}
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	orderInvoices, err := a.orderInvoiceRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	for _, inv := range orderInvoices {
		res = append(res, entity.OrderInvoiceDetail{
			OrderInvoice: inv,
		})
	}

	total, err := a.orderInvoiceRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return res, meta, nil
}

func (a *OrderInvoiceUsecase) Show(orderInvoiceID uint64) (*entity.OrderInvoiceDetail, error) {
	orderInvoiceM, err := a.orderInvoiceRepo.FindByID(orderInvoiceID)
	if err != nil {
		return nil, err
	}

	orderTransactions, err := a.orderTransactionRepo.FindByOrderInvoiceID(orderInvoiceID)
	if err != nil {
		return nil, err
	}

	var orderTransactionDetails []entity.OrderTransactionDetail
	for _, t := range orderTransactions {
		p, _ := a.productRepo.FindByID(t.ProductID)

		orderTransactionDetails = append(orderTransactionDetails, entity.OrderTransactionDetail{
			OrderTransaction: t,
			Product:          *p,
		})
	}

	return &entity.OrderInvoiceDetail{
		OrderInvoice:     *orderInvoiceM,
		OrderTransaction: orderTransactionDetails,
	}, err
}

func (a *OrderInvoiceUsecase) SyncOrderInvoice(orderInvoiceID uint64) error {
	orderTransactions, err := a.orderTransactionRepo.FindByOrderInvoiceID(orderInvoiceID)
	if err != nil {
		return err
	}

	totalBuyPrice := 0.0

	for _, t := range orderTransactions {
		totalBuyPrice += t.TotalBuyPrice
	}
	orderInvoiceM, err := a.orderInvoiceRepo.FindByID(orderInvoiceID)
	if err != nil {
		return err
	}

	orderInvoiceM.TotalBuyPrice = totalBuyPrice

	err = a.orderInvoiceRepo.Update(orderInvoiceM, nil)
	if err != nil {
		return err
	}

	return nil
}
