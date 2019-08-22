package mysql

import (
	"log"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/role"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) role.RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Insert(role *models.Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		log.Println("error-insert-role:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (r *RoleRepository) Update(role *models.Role) error {
	err := r.db.Save(role).Error
	if err != nil {
		log.Println("error-update-role:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (r *RoleRepository) FindAll(config request_util.PaginationConfig) ([]models.Role, error) {
	var results []models.Role

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := r.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-role:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (r *RoleRepository) FindByID(roleID uint64) (*models.Role, error) {
	model := models.Role{}

	err := r.db.Where("id = ?", roleID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-role-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (r *RoleRepository) FindByName(roleName string) (*models.Role, error) {
	model := models.Role{}

	err := r.db.Where("name = ?", roleName).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-role-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}
