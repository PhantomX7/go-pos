package stockmutation

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type StockMutationRepository interface {
	Insert(stockMutation *models.StockMutation, tx *gorm.DB) error
	FindAll(config request_util.PaginationConfig) ([]models.StockMutation, error)
	FindByID(stockMutationID uint64) (models.StockMutation, error)
	Count(config request_util.PaginationConfig) (int, error)
}
