package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) invoice.InvoiceRepository {
	return &InvoiceRepository{
		db: db,
	}
}

func (i *InvoiceRepository) Insert(invoice *models.Invoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Create(invoice).Error
	if err != nil {
		log.Println("error-insert-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *InvoiceRepository) Update(invoice *models.Invoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Save(invoice).Error
	if err != nil {
		log.Println("error-update-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *InvoiceRepository) Delete(invoice *models.Invoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(invoice).Error
	if err != nil {
		log.Println("error-delete-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *InvoiceRepository) FindAll(config request_util.PaginationConfig) ([]models.Invoice, error) {
	var results []models.Invoice

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
		log.Println("error-find-invoice:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (i *InvoiceRepository) FindByID(invoiceID uint64) (*models.Invoice, error) {
	model := models.Invoice{}

	err := i.db.Where("id = ?", invoiceID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-invoice-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (i *InvoiceRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := i.db.Model(&models.Invoice{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-invoice:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
