package mysql_test

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/service/user/delivery/http/request"
	userRepo "github.com/PhantomX7/go-pos/service/user/repository/mysql"
	"github.com/PhantomX7/go-pos/utils/errors"
	"testing"

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
		suite.Assert().Equal(uint64(1), userM.ID)

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

		userM, err := suite.repository.FindByID(10000)
		suite.Assert().Equal(errors.ErrNotFound, err)
		suite.Assert().Nil(userM)
	})
}

func (suite *MysqlTestSuite) TestFindByUsername() {
	suite.Run("user found", func() {
		reply := []map[string]interface{}{{"username": "test"}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		userM, err := suite.repository.FindByUsername("test")
		suite.Nil(err)
		suite.Assert().Equal("test", userM.Username)
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

		userM, err := suite.repository.FindByUsername("test2")
		suite.Assert().Equal(errors.ErrNotFound, err)
		suite.Assert().Nil(userM)
	})
}

func (suite *MysqlTestSuite) TestFindAll() {
	config := request.NewUserPaginationConfig(
		map[string][]string{
			"id": []string{"1"},
		})
	suite.Run("1 user found", func() {

		reply := []map[string]interface{}{{"id": 1, "username": "test"}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		users, err := suite.repository.FindAll(config)
		suite.Nil(err)
		suite.Assert().Equal(1, len(users))
		suite.Assert().Equal(uint64(1), users[0].ID)
	})

	suite.Run("multiple user found", func() {
		reply := []map[string]interface{}{
			{"id": 1, "username": "test"},
			{"id": 2, "username": "test"},
			{"id": 3, "username": "test"},
		}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		users, err := suite.repository.FindAll(config)
		suite.Nil(err)
		suite.Assert().Equal(3, len(users))
	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		users, err := suite.repository.FindAll(config)

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
		suite.Assert().Nil(users)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{} // empty response
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		users, err := suite.repository.FindAll(config)
		suite.Assert().Nil(err)
		suite.Assert().Equal([]models.User{}, users)
	})
}

func (suite *MysqlTestSuite) TestCount() {
	config := request.NewUserPaginationConfig(
		map[string][]string{
			"id": []string{"1"},
		})
	suite.Run("1 user found", func() {

		reply := []map[string]interface{}{{"count": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		total, err := suite.repository.Count(config)
		suite.Nil(err)
		suite.Assert().Equal(1, total)
	})

	suite.Run("3 user found", func() {
		reply := []map[string]interface{}{{"count": 3}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		total, err := suite.repository.Count(config)
		suite.Nil(err)
		suite.Assert().Equal(3, total)
	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		total, err := suite.repository.Count(config)

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
		suite.Assert().Equal(0, total)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{{"count": 0}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		total, err := suite.repository.Count(config)
		suite.Assert().Nil(err)
		suite.Assert().Equal(0, total)
	})
}

func (suite *MysqlTestSuite) initializeUser() *models.User {
	return &models.User{}
}
