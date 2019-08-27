package user

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
)

type UserRepository interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	FindAll(config request_util.PaginationConfig) ([]models.User, error)
	FindByID(userID uint64) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Count(config request_util.PaginationConfig) (int, error)
}
