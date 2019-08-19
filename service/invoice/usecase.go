package invoice

import (
	"github.com/PhantomX7/go-pos/entity"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type InvoiceUsecase interface {
	Create(request request.InvoiceCreateRequest) (models.Invoice, error)
	Update(invoiceID uint64, request request.InvoiceUpdateRequest) (models.Invoice, error)
	Delete(invoiceID uint64) error
	Index(paginationConfig request.InvoicePaginationConfig) ([]entity.InvoiceDetail, response_util.PaginationMeta, error)
	Show(invoiceID uint64) (models.Invoice, error)
}

func UpdateInvoice(
	invoiceID uint64,
	transactionRepo transaction.TransactionRepository,
	invoiceRepo InvoiceRepository,
) {
	transactions, err := transactionRepo.FindByInvoiceID(invoiceID)
	if err != nil {
		return
	}
	totalCapital := 0.0;
	totalSellPrice := 0.0;
	totalProfit := 0.0;
	for _, t := range transactions {
		totalCapital += t.Amount * t.CapitalPrice
		totalSellPrice += t.TotalSellPrice
		totalProfit += (t.SellPrice - t.CapitalPrice) * t.Amount
	}
	invoiceM, err := invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return
	}

	invoiceM.TotalCapital = totalCapital
	invoiceM.TotalSellPrice = totalSellPrice
	invoiceM.TotalProfit = totalProfit

	err = invoiceRepo.Update(&invoiceM, nil)
	if err != nil {
		return
	}
}
