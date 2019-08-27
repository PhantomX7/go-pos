package usecase_test

import (
	"errors"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/service/user/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/user/mocks"
	"github.com/PhantomX7/go-pos/service/user/usecase"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	repository *mocks.UserRepository
	usecase    user.UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.repository = new(mocks.UserRepository)
	suite.usecase = usecase.New(suite.repository)
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (suite *UserUsecaseTestSuite) initializeUser() *models.User {
	return &models.User{}
}

func (suite *UserUsecaseTestSuite) TestUpdate() {
	m := suite.initializeUser()

	suite.Run("update success", func() {
		suite.repository.On("FindByID", uint64(1)).Return(m, nil).Once()
		suite.repository.On("Update", m).Return(nil).Once()
		userM, err := suite.usecase.Update(1, request.UserUpdateRequest{Username: "test"})

		suite.Nil(err)
		suite.NotNil(userM)

	})

	suite.Run("update failed when querying user", func() {
		expectedErr := errors.New("unknown error")
		suite.repository.On("FindByID", uint64(1)).Return(nil, expectedErr).Once()
		userM, err := suite.usecase.Update(1, request.UserUpdateRequest{Username: "test"})

		suite.NotNil(err)
		suite.Equal(expectedErr, err)
		suite.Nil(userM)
	})

	suite.Run("update failed when updating", func() {
		expectedErr := errors.New("unknown error")
		suite.repository.On("FindByID", uint64(1)).Return(m, nil).Once()
		suite.repository.On("Update", m).Return(expectedErr).Once()
		userM, err := suite.usecase.Update(1, request.UserUpdateRequest{Username: "test"})

		suite.NotNil(err)
		suite.Equal(expectedErr, err)
		suite.Nil(userM)
	})

	//expectedErr := errors.New("unknown error")
	//
	//suite.repository.On("FindByID", uint64(1)).Return(&models.User{ID: 1}, expectedErr)
	//
	//user, err := suite.usecase.Update(1, request.UserUpdateRequest{Username: "test"})
	//
	//suite.Equal(expectedErr, err)
}
