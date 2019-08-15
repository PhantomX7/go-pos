package invoice

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request"
)

type StockMutationRepository interface {
	Insert(invoice *models.StockMutation) error
	Update(invoice *models.StockMutation) error
	Delete(invoice *models.StockMutation) error
	FindAll(config request.PaginationConfig) ([]models.StockMutation, error)
	FindByID(invoiceID int64) (models.StockMutation, error)
	Count(config request.PaginationConfig) (int, error)
}
