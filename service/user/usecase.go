package user

import (
	"github.com/PhantomX7/go-pos/service/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type UserUsecase interface {
	Create(request request.UserCreateRequest) (*models.User, error)
	Update(userID uint64, request request.UserUpdateRequest) (*models.User, error)
	Index(paginationConfig request_util.PaginationConfig) ([]models.User, response_util.PaginationMeta, error)
	Show(userID uint64) (*models.User, error)
}
