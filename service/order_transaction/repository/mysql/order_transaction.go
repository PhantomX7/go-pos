package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_transaction"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type OrderTransactionRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) order_transaction.OrderTransactionRepository {
	return &OrderTransactionRepository{
		db: db,
	}
}

func (t *OrderTransactionRepository) Insert(orderTransaction *models.OrderTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Create(orderTransaction).Error
	if err != nil {
		log.Println("error-insert-order-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *OrderTransactionRepository) Update(orderTransaction *models.OrderTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Save(orderTransaction).Error
	if err != nil {
		log.Println("error-update-order-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *OrderTransactionRepository) Delete(orderTransaction *models.OrderTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(orderTransaction).Error
	if err != nil {
		log.Println("error-delete-order-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *OrderTransactionRepository) FindAll(config request_util.PaginationConfig) ([]models.OrderTransaction, error) {
	var results []models.OrderTransaction

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := t.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-order-transaction:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (t *OrderTransactionRepository) FindByID(orderTransactionID uint64) (*models.OrderTransaction, error) {
	model := models.OrderTransaction{}

	err := t.db.Where("id = ?", orderTransactionID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-order-transaction-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (t *OrderTransactionRepository) FindByOrderInvoiceID(orderInvoiceID uint64) ([]models.OrderTransaction, error) {
	model := []models.OrderTransaction{}

	err := t.db.Where("order_invoice_id = ?", orderInvoiceID).Find(&model).Error

	if err != nil {
		log.Println("error-find-order-transaction-by-order-invoice-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (t *OrderTransactionRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := t.db.Model(&models.OrderTransaction{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-order-transaction:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
