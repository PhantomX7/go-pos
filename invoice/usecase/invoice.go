package usecase

import (
	"github.com/PhantomX7/go-pos/customer"
	"github.com/PhantomX7/go-pos/invoice"
	"github.com/PhantomX7/go-pos/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/invoice/delivery/http/response"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type InvoiceUsecase struct {
	invoiceRepo  invoice.InvoiceRepository
	customerRepo customer.CustomerRepository
}

func NewInvoiceUsecase(invoiceRepo invoice.InvoiceRepository, customerRepo customer.CustomerRepository) invoice.InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo:  invoiceRepo,
		customerRepo: customerRepo,
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
	err := a.invoiceRepo.Insert(&invoiceM)
	if err != nil {
		return invoiceM, err
	}
	return invoiceM, nil
}

func (a *InvoiceUsecase) Update(invoiceID int64, request request.InvoiceUpdateRequest) (models.Invoice, error) {
	invoiceM, err := a.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return invoiceM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&invoiceM, &request)

	err = a.invoiceRepo.Update(&invoiceM)
	if err != nil {
		return invoiceM, err
	}
	return invoiceM, nil
}

func (a *InvoiceUsecase) Delete(invoiceID int64) error {
	err := a.invoiceRepo.Delete(&models.Invoice{ID: invoiceID})
	if err != nil {
		return err
	}

	return nil
}

func (a *InvoiceUsecase) Index(paginationConfig request.InvoicePaginationConfig) (response.InvoiceResponse, response_util.PaginationMeta, error) {
	var res []response.InvoiceDetail
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
		res = append(res, response.InvoiceDetail{Invoice: inv, Customer: &c})
	}

	total, err := a.invoiceRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return res, meta, nil
}

func (a *InvoiceUsecase) Show(invoiceID int64) (models.Invoice, error) {
	return a.invoiceRepo.FindByID(invoiceID)
}
