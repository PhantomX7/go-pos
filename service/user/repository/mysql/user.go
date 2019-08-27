package mysql

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Insert(user *models.User) error {
	err := encryptUserPassword(user)
	if err != nil {
		return err
	}
	err = u.db.Create(user).Error
	if err != nil {
		log.Println("error-insert-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (u *UserRepository) Update(user *models.User) error {
	err := encryptUserPassword(user)
	if err != nil {
		return err
	}
	err = u.db.Save(user).Error
	if err != nil {
		log.Println("error-update-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (u *UserRepository) FindAll(config request_util.PaginationConfig) ([]models.User, error) {
	var results []models.User

	//default order
	order := "id"
	orderConfig := config.Order()
	if orderConfig != "" {
		order = orderConfig
	}
	sc := config.SearchClause()
	err := u.db.Order(order).
		Limit(config.Limit()).
		Offset(config.Offset()).
		Where(sc.Query, sc.Args...).
		Find(&results).Error
	if err != nil {
		log.Println("error-find-user:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return results, nil
}

func (u *UserRepository) FindByID(userID uint64) (*models.User, error) {
	model := models.User{}

	err := u.db.Where("id = ?", userID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-user-by-id:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (u *UserRepository) FindByUsername(username string) (*models.User, error) {
	model := models.User{}

	err := u.db.Where("username = ?", username).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-user-by-username:", err)
		return nil, errors.ErrUnprocessableEntity
	}

	return &model, nil
}

func (u *UserRepository) Count(config request_util.PaginationConfig) (int, error) {
	var count int

	sc := config.SearchClause()
	err := u.db.Model(&models.User{}).Where(sc.Query, sc.Args...).
		Count(&count).Error
	if err != nil {
		log.Println("error-count-user:", err)
		return 0, errors.ErrUnprocessableEntity
	}

	return count, nil
}

func encryptUserPassword(user *models.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error-encrypting-password:", err)
		return errors.ErrUnprocessableEntity
	}
	user.Password = string(password)
	return nil
}
