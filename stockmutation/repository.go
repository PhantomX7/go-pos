package stockmutation

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type StockMutationRepository interface {
	Insert(stockMutation *models.StockMutation) error
	FindAll(config request_util.PaginationConfig) ([]models.StockMutation, error)
	FindByID(stockMutationID int64) (models.StockMutation, error)
	Count(config request_util.PaginationConfig) (int, error)
}
