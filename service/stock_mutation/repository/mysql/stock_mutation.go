package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/stock_mutation"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type StockMutationRepository struct {
	db *gorm.DB
}

func NewStockMutationRepository(db *gorm.DB) stock_mutation.StockMutationRepository {
	return &StockMutationRepository{
		db: db,
	}
}

func (i *StockMutationRepository) Insert(stockMutation *models.StockMutation, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Create(stockMutation).Error
	if err != nil {
		log.Println("error-insert-stockMutation:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *StockMutationRepository) Delete(stockMutation *models.StockMutation, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(stockMutation).Error
	if err != nil {
		log.Println("error-delete-stockMutation:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *StockMutationRepository) FindAll(config request_util.PaginationConfig) ([]models.StockMutation, error) {
	var results []models.StockMutation

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := i.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-stockMutation:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (i *StockMutationRepository) FindByID(stockMutationID uint64) (models.StockMutation, error) {
	model := models.StockMutation{}

	err := i.db.Where("id = ?", stockMutationID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-stock-mutation-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (i *StockMutationRepository) FindByProductID(productID uint64) (models.StockMutation, error) {
	model := models.StockMutation{}

	err := i.db.Where("product_id = ?", productID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-stock-mutation-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (i *StockMutationRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := i.db.Model(&models.StockMutation{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-stock-mutation:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
