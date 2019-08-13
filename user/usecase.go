package user

import (
	"github.com/PhantomX7/go-pos/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/user/delivery/http/response"
	"github.com/PhantomX7/go-pos/models"
)

type UserUsecase interface {
	Create(user models.User) (models.User, error)
	Update(userID int64, user request.UserUpdateRequest) (models.User, error)
	Index(paginationConfig request.PaginationConfig) ([]models.User, response.UserPaginationMeta, error)
	Show(userID int64) (models.User, error)
}
