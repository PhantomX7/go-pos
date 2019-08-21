package usecase

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/customer"
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type InvoiceUsecase struct {
	invoiceRepo     invoice.InvoiceRepository
	customerRepo    customer.CustomerRepository
	transactionRepo transaction.TransactionRepository
	productRepo     product.ProductRepository
}

func NewInvoiceUsecase(
	invoiceRepo invoice.InvoiceRepository,
	customerRepo customer.CustomerRepository,
	transactionRepo transaction.TransactionRepository,
	productRepo product.ProductRepository,
) invoice.InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo:     invoiceRepo,
		customerRepo:    customerRepo,
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

func (a *InvoiceUsecase) Create(request request.InvoiceCreateRequest) (models.Invoice, error) {
	invoiceM := models.Invoice{
		CustomerId:     request.CustomerId,
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
		return invoiceM, err
	}
	return invoiceM, nil
}

func (a *InvoiceUsecase) Update(invoiceID uint64, request request.InvoiceUpdateRequest) (models.Invoice, error) {
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return invoiceM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&invoiceM, &request)

	err = a.invoiceRepo.Update(&invoiceM, nil)
	if err != nil {
		return invoiceM, err
	}
	return invoiceM, nil
}

func (a *InvoiceUsecase) Delete(invoiceID uint64) error {
	err := a.invoiceRepo.Delete(&models.Invoice{ID: invoiceID}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *InvoiceUsecase) Index(paginationConfig request.InvoicePaginationConfig) ([]entity.InvoiceDetail, response_util.PaginationMeta, error) {
	var res []entity.InvoiceDetail
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
		c, _ := a.customerRepo.FindByID(inv.CustomerId)
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

func (a *InvoiceUsecase) Show(invoiceID uint64) (entity.InvoiceDetail, error) {
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return entity.InvoiceDetail{}, err
	}

	customerM, err := a.customerRepo.FindByID(invoiceM.CustomerId)
	if err != nil {
		return entity.InvoiceDetail{}, err
	}

	transactions, err := a.transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return entity.InvoiceDetail{}, err
	}

	transactionDetails := []entity.TransactionDetail{}
	for _, t := range transactions {
		p, _ := a.productRepo.FindByID(t.ProductId)
		transactionDetails = append(transactionDetails, entity.TransactionDetail{
			Transaction: t,
			Product:     p,
		})
	}

	return entity.InvoiceDetail{
		Invoice:      invoiceM,
		Customer:     customerM,
		Transactions: transactionDetails,
	}, err
}

func (a *InvoiceUsecase) SyncInvoice(invoiceID uint64) error {
	transactions, err := a.transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return err
	}
	totalCapital := 0.0;
	totalSellPrice := 0.0;
	totalProfit := 0.0;
	for _, t := range transactions {
		totalCapital += t.Amount * t.CapitalPrice
		totalSellPrice += t.TotalSellPrice
		totalProfit += (t.SellPrice - t.CapitalPrice) * t.Amount
	}
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return err
	}

	invoiceM.TotalCapital = totalCapital
	invoiceM.TotalSellPrice = totalSellPrice
	invoiceM.TotalProfit = totalProfit

	err = a.invoiceRepo.Update(&invoiceM, nil)
	if err != nil {
		return err
	}

	return nil
}
