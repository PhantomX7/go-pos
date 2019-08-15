package usecase

import (
	"github.com/PhantomX7/go-pos/invoice"
	"github.com/PhantomX7/go-pos/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response"
	"github.com/jinzhu/copier"
)

// apply business logic here

type InvoiceUsecase struct {
	invoiceRepo invoice.InvoiceRepository
}

func NewInvoiceUsecase(invoiceRepo invoice.InvoiceRepository) invoice.InvoiceUsecase {
	return &InvoiceUsecase{
		invoiceRepo: invoiceRepo,
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

func (a *InvoiceUsecase) Index(paginationConfig request.InvoicePaginationConfig) ([]models.Invoice, response.PaginationMeta, error) {
	meta := response.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	invoices, err := a.invoiceRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := a.invoiceRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return invoices, meta, nil
}

func (a *InvoiceUsecase) Show(invoiceID int64) (models.Invoice, error) {
	return a.invoiceRepo.FindByID(invoiceID)
}
