package mysql_test

import (
	"testing"
	"time"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/user"
	userRepo "github.com/PhantomX7/go-pos/service/user/repository/mysql"
	"github.com/PhantomX7/go-pos/utils/errors"

	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
)

const databaseName = "pos"

type MysqlTestSuite struct {
	suite.Suite
	repository user.UserRepository
}

func TestUserMysql(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}

func (suite *MysqlTestSuite) SetupSuite() {
	suite.repository = suite.initializeRepository()
}

func (suite *MysqlTestSuite) initializeRepository() user.UserRepository {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, err := gorm.Open(mocket.DriverName, databaseName)

	if err != nil {
		panic(err)
	}

	return userRepo.New(db)
}

func (suite *MysqlTestSuite) TestInsert() {
	m := suite.initializeUser()

	suite.Run("successful creation", func() {
		mocket.Catcher.Reset().NewMock().WithError(nil)

		err := suite.repository.Insert(m)
		suite.Nil(err)
	})

	suite.Run("with error response", func() {
		mocket.Catcher.Reset().NewMock().WithExecException()

		err := suite.repository.Insert(m)
		suite.NotNil(err)
	})
}

func (suite *MysqlTestSuite) TestUpdate() {
	m := suite.initializeUser()

	suite.Run("successful update", func() {
		mocket.Catcher.Reset().NewMock().WithError(nil)

		err := suite.repository.Update(m)
		suite.Nil(err)
	})

	suite.Run("with error response", func() {
		mocket.Catcher.Reset().NewMock().WithExecException()

		err := suite.repository.Update(m)
		suite.NotNil(err)
	})
}

func (suite *MysqlTestSuite) TestFindByUserID() {
	suite.Run("user found", func() {
		reply := []map[string]interface{}{{"id": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		userM, err := suite.repository.FindByID(1)
		suite.Nil(err)
		suite.Equal(userM.ID, uint(1))

	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		_, err := suite.repository.FindByID(10001)

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{{}} // empty response
		mocket.Catcher.Reset().NewMock().WithArgs(1000).WithReply(reply)

		_, err := suite.repository.FindByID(10000)
		suite.Assert().Equal(errors.ErrNotFound, err)
	})
}

func (suite *MysqlTestSuite) TestFindByUsername() {
	suite.Run("user found", func() {
		reply := []map[string]interface{}{{"username": "test"}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		_, err := suite.repository.FindByUsername("test")
		suite.Nil(err)
	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		_, err := suite.repository.FindByUsername("test")

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{{}} // empty response
		mocket.Catcher.Reset().NewMock().WithArgs("test").WithReply(reply)

		_, err := suite.repository.FindByUsername("test2")
		suite.Assert().Equal(errors.ErrNotFound, err)
	})
}

func (suite *MysqlTestSuite) initializeUser() *models.User {
	return &models.User{
		ID:        1,
		Username:  "test",
		Password:  "test",
		RoleId:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
