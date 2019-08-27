package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type OrderInvoiceRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) order_invoice.OrderInvoiceRepository {
	return &OrderInvoiceRepository{
		db: db,
	}
}

func (i *OrderInvoiceRepository) Insert(orderInvoice *models.OrderInvoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Create(orderInvoice).Error
	if err != nil {
		log.Println("error-insert-order-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *OrderInvoiceRepository) Update(orderInvoice *models.OrderInvoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Save(orderInvoice).Error
	if err != nil {
		log.Println("error-update-order-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *OrderInvoiceRepository) Delete(orderInvoice *models.OrderInvoice, tx *gorm.DB) error {
	var db = i.db
	if tx != nil {
		db = tx
	}
	err := db.Delete(orderInvoice).Error
	if err != nil {
		log.Println("error-delete-order-invoice:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (i *OrderInvoiceRepository) FindAll(config request_util.PaginationConfig) ([]models.OrderInvoice, error) {
	results := []models.OrderInvoice{}

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	log.Println(sc.Query, sc.Args)
	err := i.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-order-invoice:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (i *OrderInvoiceRepository) FindByID(orderInvoiceID uint64) (*models.OrderInvoice, error) {
	model := models.OrderInvoice{}

	err := i.db.Where("id = ?", orderInvoiceID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-order-invoice-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (i *OrderInvoiceRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := i.db.Model(&models.OrderInvoice{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-order-invoice:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
