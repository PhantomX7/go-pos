package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) transaction.TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (t *TransactionRepository) Insert(transaction *models.Transaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Create(transaction).Error
	if err != nil {
		log.Println("error-insert-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *TransactionRepository) Update(transaction *models.Transaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Save(transaction).Error
	if err != nil {
		log.Println("error-update-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *TransactionRepository) Delete(transaction *models.Transaction, tx *gorm.DB) error {
	var db = t.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(transaction).Error
	if err != nil {
		log.Println("error-delete-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (t *TransactionRepository) FindAll(config request_util.PaginationConfig) ([]models.Transaction, error) {
	var results []models.Transaction

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
		log.Println("error-find-transaction:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (t *TransactionRepository) FindByID(transactionID uint64) (models.Transaction, error) {
	model := models.Transaction{}

	err := t.db.Where("id = ?", transactionID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-transaction-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (t *TransactionRepository) FindByInvoiceID(invoiceID uint64) ([]models.Transaction, error) {
	model := []models.Transaction{}

	err := t.db.Where("invoice_id = ?", invoiceID).Find(&model).Error

	if err != nil {
		log.Println("error-find-transaction-by-invoice-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (t *TransactionRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := t.db.Model(&models.Transaction{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-transaction:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
