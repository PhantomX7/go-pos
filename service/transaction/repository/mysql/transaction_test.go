package mysql_test

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/service/transaction/delivery/http/request"
	transactionRepo "github.com/PhantomX7/go-pos/service/transaction/repository/mysql"
	"github.com/PhantomX7/go-pos/utils/errors"
	"testing"

	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
)

const databaseName = "pos"

type MysqlTestSuite struct {
	suite.Suite
	repository transaction.TransactionRepository
}

func TestTransactionMysql(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}

func (suite *MysqlTestSuite) SetupSuite() {
	suite.repository = suite.initializeRepository()
}

func (suite *MysqlTestSuite) initializeRepository() transaction.TransactionRepository {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, err := gorm.Open(mocket.DriverName, databaseName)

	if err != nil {
		panic(err)
	}

	return transactionRepo.New(db)
}

func (suite *MysqlTestSuite) TestInsert() {
	m := suite.initializeTransaction()

	suite.Run("successful creation", func() {
		mocket.Catcher.Reset().NewMock().WithError(nil)

		err := suite.repository.Insert(m, nil)
		suite.Nil(err)
	})

	suite.Run("with error response", func() {
		mocket.Catcher.Reset().NewMock().WithExecException()

		err := suite.repository.Insert(m, nil)
		suite.NotNil(err)
	})
}

func (suite *MysqlTestSuite) TestUpdate() {
	m := suite.initializeTransaction()

	suite.Run("successful update", func() {
		mocket.Catcher.Reset().NewMock().WithError(nil)

		err := suite.repository.Update(m, nil)
		suite.Nil(err)
	})

	suite.Run("with error response", func() {
		mocket.Catcher.Reset().NewMock().WithExecException()

		err := suite.repository.Update(m, nil)
		suite.NotNil(err)
	})
}

func (suite *MysqlTestSuite) TestDelete() {
	m := suite.initializeTransaction()

	suite.Run("successful delete", func() {
		mocket.Catcher.Reset().NewMock().WithError(nil)

		err := suite.repository.Delete(m, nil)
		suite.Nil(err)
	})

	suite.Run("with error response", func() {
		mocket.Catcher.Reset().NewMock().WithExecException()

		err := suite.repository.Delete(m, nil)
		suite.NotNil(err)
	})
}

func (suite *MysqlTestSuite) TestFindByTransactionID() {
	suite.Run("transaction found", func() {
		reply := []map[string]interface{}{{"id": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		transactionM, err := suite.repository.FindByID(1)
		suite.Nil(err)
		suite.Assert().Equal(uint64(1), transactionM.ID)

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

		transactionM, err := suite.repository.FindByID(10000)
		suite.Assert().Equal(errors.ErrNotFound, err)
		suite.Assert().Nil(transactionM)
	})
}

func (suite *MysqlTestSuite) TestFindByInvoiceID() {
	suite.Run("transactions found", func() {
		reply := []map[string]interface{}{{"id": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		transactions, err := suite.repository.FindByInvoiceID(1)
		suite.Nil(err)
		suite.Assert().Equal(1, len(transactions))
	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		_, err := suite.repository.FindByInvoiceID(1)

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{{}} // empty response
		mocket.Catcher.Reset().NewMock().WithArgs("test").WithReply(reply)

		transactions, err := suite.repository.FindByInvoiceID(1)
		suite.Assert().Nil(err)
		suite.Assert().Equal([]models.Transaction{}, transactions)
	})
}

func (suite *MysqlTestSuite) TestFindAll() {
	config := request.NewTransactionPaginationConfig(
		map[string][]string{
			"id": []string{"1"},
		})
	suite.Run("1 transaction found", func() {

		reply := []map[string]interface{}{{"id": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		transactions, err := suite.repository.FindAll(config)
		suite.Nil(err)
		suite.Assert().Equal(1, len(transactions))
		suite.Assert().Equal(uint64(1), transactions[0].ID)
	})

	suite.Run("multiple transaction found", func() {
		reply := []map[string]interface{}{
			{"id": 1},
			{"id": 2},
			{"id": 3},
		}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		transactions, err := suite.repository.FindAll(config)
		suite.Nil(err)
		suite.Assert().Equal(3, len(transactions))
	})

	suite.Run("with unknown database error", func() {
		mocket.Catcher.Reset().NewMock().WithQueryException()
		transactions, err := suite.repository.FindAll(config)

		suite.Assert().NotNil(err)
		suite.Assert().Equal(errors.ErrUnprocessableEntity, err)
		suite.Assert().Nil(transactions)
	})

	suite.Run("Record not found", func() {
		reply := []map[string]interface{}{} // empty response
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		transactions, err := suite.repository.FindAll(config)
		suite.Assert().Nil(err)
		suite.Assert().Equal([]models.Transaction{}, transactions)
	})
}

func (suite *MysqlTestSuite) TestCount() {
	config := request.NewTransactionPaginationConfig(
		map[string][]string{
			"id": []string{"1"},
		})
	suite.Run("1 transaction found", func() {

		reply := []map[string]interface{}{{"count": 1}}
		mocket.Catcher.Reset().NewMock().WithReply(reply)

		total, err := suite.repository.Count(config)
		suite.Nil(err)
		suite.Assert().Equal(1, total)
	})

	suite.Run("3 transaction found", func() {
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

func (suite *MysqlTestSuite) initializeTransaction() *models.Transaction {
	return &models.Transaction{}
}
