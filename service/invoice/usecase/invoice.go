package usecase

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/customer"
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/return_transaction"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type InvoiceUsecase struct {
	invoiceRepo           invoice.InvoiceRepository
	customerRepo          customer.CustomerRepository
	transactionRepo       transaction.TransactionRepository
	productRepo           product.ProductRepository
	returnTransactionRepo return_transaction.ReturnTransactionRepository
}

func New(
	invoiceRepo invoice.InvoiceRepository,
	customerRepo customer.CustomerRepository,
	transactionRepo transaction.TransactionRepository,
	productRepo product.ProductRepository,
	returnTransactionRepo return_transaction.ReturnTransactionRepository,
) invoice.InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo:           invoiceRepo,
		customerRepo:          customerRepo,
		transactionRepo:       transactionRepo,
		productRepo:           productRepo,
		returnTransactionRepo: returnTransactionRepo,
	}
}

func (a *InvoiceUsecase) Create(request request.InvoiceCreateRequest) (*models.Invoice, error) {
	invoiceM := models.Invoice{
		CustomerID:     request.CustomerId,
		Description:    request.Description,
		Date:           request.Date,
		PaymentStatus:  request.PaymentStatus,
		PaymentType:    request.PaymentType,
		TotalCapital:   0,
		TotalSellPrice: 0,
		TotalProfit:    0,
	}
	err := a.invoiceRepo.Insert(&invoiceM, nil)
	if err != nil {
		return nil, err
	}
	return &invoiceM, nil
}

func (a *InvoiceUsecase) Update(invoiceID uint64, request request.InvoiceUpdateRequest) (*models.Invoice, error) {
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return invoiceM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(invoiceM, &request)

	err = a.invoiceRepo.Update(invoiceM, nil)
	if err != nil {
		return invoiceM, err
	}
	return invoiceM, nil
}

func (a *InvoiceUsecase) Delete(invoiceID uint64) error {
	transactions, err := a.transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return err
	}
	if len(transactions) > 0 {
		err := errors.ErrUnprocessableEntity
		err.Message = "please delete all transaction first"
		return err
	}
	err = a.invoiceRepo.Delete(&models.Invoice{ID: invoiceID}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *InvoiceUsecase) Index(paginationConfig request.InvoicePaginationConfig) ([]entity.InvoiceDetail, response_util.PaginationMeta, error) {
	res := []entity.InvoiceDetail{}
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	invoices, err := a.invoiceRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	for _, inv := range invoices {
		c, _ := a.customerRepo.FindByID(inv.CustomerID)
		res = append(res, entity.InvoiceDetail{
			Invoice:  inv,
			Customer: c,
		})
	}

	total, err := a.invoiceRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return res, meta, nil
}

func (a *InvoiceUsecase) Show(invoiceID uint64) (*entity.InvoiceDetail, error) {
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}

	customerM, err := a.customerRepo.FindByID(invoiceM.CustomerID)
	if err != nil {
		return nil, err
	}

	transactions, err := a.transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return nil, err
	}

	var transactionDetails []entity.TransactionDetail
	for _, t := range transactions {
		p, _ := a.productRepo.FindByID(t.ProductID)
		r, _ := a.returnTransactionRepo.FindByTransactionID(t.ID)

		transactionDetails = append(transactionDetails, entity.TransactionDetail{
			Transaction:       t,
			Product:           *p,
			ReturnTransaction: r,
		})
	}

	return &entity.InvoiceDetail{
		Invoice:      *invoiceM,
		Customer:     customerM,
		Transactions: transactionDetails,
	}, err
}

func (a *InvoiceUsecase) SyncInvoice(invoiceID uint64) error {
	transactions, err := a.transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return err
	}

	totalCapital := 0.0
	totalSellPrice := 0.0
	totalProfit := 0.0

	for _, t := range transactions {
		amount := t.Amount
		r, _ := a.returnTransactionRepo.FindByTransactionID(t.ID)
		if r != nil {
			amount -= r.Amount
		}

		totalCapital += amount * t.CapitalPrice
		totalSellPrice += amount * t.SellPrice
		totalProfit += (t.SellPrice - t.CapitalPrice) * amount
	}
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return err
	}

	invoiceM.TotalCapital = totalCapital
	invoiceM.TotalSellPrice = totalSellPrice
	invoiceM.TotalProfit = totalProfit

	err = a.invoiceRepo.Update(invoiceM, nil)
	if err != nil {
		return err
	}

	return nil
}
