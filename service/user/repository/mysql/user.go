package mysql

import (
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (a *UserRepository) Insert(user *models.User) error {
	err := encryptUserPassword(user)
	if err != nil {
		return err
	}
	err = a.db.Create(user).Error
	if err != nil {
		log.Println("error-insert-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (a *UserRepository) Update(user *models.User) error {
	err := encryptUserPassword(user)
	if err != nil {
		return err
	}
	err = a.db.Save(user).Error
	if err != nil {
		log.Println("error-update-user:", err)
		return errors.ErrUnprocessableEntity
	}
	return nil
}

func (a *UserRepository) FindAll(config request_util.PaginationConfig) ([]models.User, response_util.PaginationMeta, error) {
	return nil, response_util.PaginationMeta{}, nil
}

func (a *UserRepository) FindByID(userID uint64) (models.User, error) {
	model := models.User{}

	err := a.db.Where("id = ?", userID).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-user-by-id:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
}

func (a *UserRepository) FindByUsername(username string) (models.User, error) {
	model := models.User{}

	err := a.db.Where("username = ?", username).First(&model).Error

	if gorm.IsRecordNotFoundError(err) {
		return model, errors.ErrNotFound
	}

	if err != nil {
		log.Println("error-find-user-by-username:", err)
		return model, errors.ErrUnprocessableEntity
	}

	return model, nil
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
