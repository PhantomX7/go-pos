package usecase

import (
	"github.com/PhantomX7/go-pos/user"
	"github.com/PhantomX7/go-pos/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/response"
	"github.com/PhantomX7/go-pos/models"
	"github.com/jinzhu/copier"
)

// apply business logic here

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (a *UserUsecase) Create(user models.User) (models.User, error) {
	err := a.userRepo.Insert(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (a *UserUsecase) Update(userID int64, user request.UserUpdateRequest) (models.User, error) {
	userM, err := a.userRepo.FindByID(userID)
	if err != nil {
		return userM, err
	}

	// copy content of request into user model found by id
	_ = copier.Copy(&userM, &user)

	err = a.userRepo.Update(&userM)
	if err != nil {
		return userM, err
	}
	return userM, nil
}

func (a *UserUsecase) Index(paginationConfig request.PaginationConfig) ([]models.User, response.PaginationMeta, error) {
	return nil, response.PaginationMeta{}, nil
}

func (a *UserUsecase) Show(userID int64) (models.User, error) {
	return a.userRepo.FindByID(userID)
}
