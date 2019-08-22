package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/return_transaction"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type ReturnTransactionRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) return_transaction.ReturnTransactionRepository {
	return &ReturnTransactionRepository{
		db: db,
	}
}

func (t *ReturnTransactionRepository) Insert(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Create(returnTransaction).Error
	if err != nil {
		log.Println("error-insert-return-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *ReturnTransactionRepository) Update(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Save(returnTransaction).Error
	if err != nil {
		log.Println("error-update-return-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *ReturnTransactionRepository) Delete(returnTransaction *models.ReturnTransaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(returnTransaction).Error
	if err != nil {
		log.Println("error-delete-return-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *ReturnTransactionRepository) FindAll(config request_util.PaginationConfig) ([]models.ReturnTransaction, error) {
	var results []models.ReturnTransaction

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
		log.Println("error-find-return-transaction:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (t *ReturnTransactionRepository) FindByID(returnTransactionID uint64) (*models.ReturnTransaction, error) {
	model := models.ReturnTransaction{}

	err := t.db.Where("id = ?", returnTransactionID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-return-transaction-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (t *ReturnTransactionRepository) FindByTransactionID(transactionID uint64) (*models.ReturnTransaction, error) {
	model := models.ReturnTransaction{}

	err := t.db.Where("transaction_id = ?", transactionID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-return-transaction-by-transaction-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (t *ReturnTransactionRepository) FindByInvoiceID(invoiceID uint64) ([]models.ReturnTransaction, error) {
	var model []models.ReturnTransaction

	err := t.db.Where("invoice_id = ?", invoiceID).First(&model).Error

	if err != nil {
		log.Println("error-find-return-transaction-by-invoice-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (t *ReturnTransactionRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := t.db.Model(&models.ReturnTransaction{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-returnTransaction:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
