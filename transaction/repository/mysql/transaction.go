package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/transaction"
	"github.com/PhantomX7/go-pos/models"
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

func (i *TransactionRepository) Insert(transaction *models.Transaction) error {
	err := i.db.Create(transaction).Error
	if err != nil {
		log.Println("error-insert-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *TransactionRepository) Update(transaction *models.Transaction) error {
	err := i.db.Save(transaction).Error
	if err != nil {
		log.Println("error-update-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *TransactionRepository) Delete(transaction *models.Transaction) error {
	err := i.db.Delete(transaction).Error
	if err != nil {
		log.Println("error-delete-transaction:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *TransactionRepository) FindAll(config request_util.PaginationConfig) ([]models.Transaction, error) {
	var results []models.Transaction

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
		log.Println("error-find-transaction:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (i *TransactionRepository) FindByID(transactionID int64) (models.Transaction, error) {
	model := models.Transaction{}

	err := i.db.Where("id = ?", transactionID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-transaction-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (i *TransactionRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := i.db.Model(&models.Transaction{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-transaction:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
