package mysql

import (
	"github.com/PhantomX7/go-pos/service/stock_adjustment"
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type StockAdjustmentRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) stock_adjustment.StockAdjustmentRepository {
	return &StockAdjustmentRepository{
		db: db,
	}
}

func (i *StockAdjustmentRepository) Insert(stockAdjustment *models.StockAdjustment, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Create(stockAdjustment).Error
	if err != nil {
		log.Println("error-insert-stockAdjustment:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *StockAdjustmentRepository) Delete(stockAdjustment *models.StockAdjustment, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(stockAdjustment).Error
	if err != nil {
		log.Println("error-delete-stockAdjustment:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *StockAdjustmentRepository) FindAll(config request_util.PaginationConfig) ([]models.StockAdjustment, error) {
	var results []models.StockAdjustment

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
		log.Println("error-find-stockAdjustment:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (i *StockAdjustmentRepository) FindByID(stockAdjustmentID uint64) (*models.StockAdjustment, error) {
	model := models.StockAdjustment{}

	err := i.db.Where("id = ?", stockAdjustmentID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-stock-adjustment-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (i *StockAdjustmentRepository) FindByProductID(productID uint64) (*models.StockAdjustment, error) {
	model := models.StockAdjustment{}

	err := i.db.Where("product_id = ?", productID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-stock-adjustment-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (i *StockAdjustmentRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := i.db.Model(&models.StockAdjustment{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-stock-adjustment:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
