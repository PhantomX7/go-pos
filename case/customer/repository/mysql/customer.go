package mysql

import (
	"github.com/PhantomX7/go-pos/case/customer"
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) customer.CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (c *CustomerRepository) Insert(customer *models.Customer) error {
	err := c.db.Create(customer).Error
	if err != nil {
		log.Println("error-insert-customer:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (c *CustomerRepository) Update(customer *models.Customer) error {
	err := c.db.Save(customer).Error
	if err != nil {
		log.Println("error-update-customer:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (c *CustomerRepository) Delete(customer *models.Customer) error {
	err := c.db.Delete(customer).Error
	if err != nil {
		log.Println("error-delete-customer:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (c *CustomerRepository) FindAll(config request_util.PaginationConfig) ([]models.Customer, error) {
	var results []models.Customer

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := c.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-customer:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (c *CustomerRepository) FindByID(customerID int64) (models.Customer, error) {
	model := models.Customer{}

	err := c.db.Where("id = ?", customerID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-customer-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (c *CustomerRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := c.db.Model(&models.Customer{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-customer:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
