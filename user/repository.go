package user

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/user/delivery/http/response"
)

type UserRepository interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	FindAll(config request.PaginationConfig) ([]models.User, response.UserPaginationMeta, error)
	FindByID(userID int64) (models.User, error)
	FindByUsername(username string) (models.User, error)
}
