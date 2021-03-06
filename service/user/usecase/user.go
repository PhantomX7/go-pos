package usecase

import (
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/service/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/request_util"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type UserUsecase struct {
	userRepo user.UserRepository
}

func New(userRepo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (a *UserUsecase) Create(request request.UserCreateRequest) (*models.User, error) {
	userM := models.User{
		Username: request.Username,
		Password: request.Password,
		RoleId:   uint64(request.RoleId),
	}

	err := a.userRepo.Insert(&userM)
	if err != nil {
		return nil, err
	}
	return &userM, nil
}

func (a *UserUsecase) Update(userID uint64, request request.UserUpdateRequest) (*models.User, error) {
	userM, err := a.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// copy content of request into user model found by id
	_ = copier.Copy(&userM, &request)

	err = a.userRepo.Update(userM)
	if err != nil {
		return nil, err
	}
	return userM, nil
}

func (a *UserUsecase) Index(paginationConfig request_util.PaginationConfig) ([]models.User, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	users, err := a.userRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := a.userRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return users, meta, nil
}

func (a *UserUsecase) Show(userID uint64) (*models.User, error) {
	return a.userRepo.FindByID(userID)
}
