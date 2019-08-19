package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (p *ProductRepository) Insert(product *models.Product, tx *gorm.DB) error {
	var db = p.db
	if tx != nil {
		db = tx
	}
	err := db.Create(product).Error
	if err != nil {
		log.Println("error-insert-product:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (p *ProductRepository) Update(product *models.Product, tx *gorm.DB) error {
	var db = p.db
	if tx != nil {
		db = tx
	}
	err := db.Save(product).Error
	if err != nil {
		log.Println("error-update-product:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (p *ProductRepository) FindAll(config request_util.PaginationConfig) ([]models.Product, error) {
	var results []models.Product

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := p.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-product:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (p *ProductRepository) FindByID(productID uint64) (models.Product, error) {
	model := models.Product{}

	err := p.db.Where("id = ?", productID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-product-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (p *ProductRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := p.db.Model(&models.Product{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-product:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}
