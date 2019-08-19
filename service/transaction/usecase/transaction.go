package usecase

import (
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/service/transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type TransactionUsecase struct {
	transactionRepo transaction.TransactionRepository
}

func NewTransactionUsecase(transactionRepo transaction.TransactionRepository) transaction.TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo: transactionRepo,
	}
}

func (a *TransactionUsecase) Create(request request.TransactionCreateRequest) (models.Transaction, error) {
	transactionM := models.Transaction{

		TotalSellPrice: 0,
	}
	err := a.transactionRepo.Insert(&transactionM)
	if err != nil {
		return transactionM, err
	}
	return transactionM, nil
}

func (a *TransactionUsecase) Update(transactionID int64, request _case.TransactionUpdateRequest) (models.Transaction, error) {
	transactionM, err := a.transactionRepo.FindByID(transactionID)
	if err != nil {
		return transactionM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&transactionM, &request)

	err = a.transactionRepo.Update(&transactionM)
	if err != nil {
		return transactionM, err
	}
	return transactionM, nil
}

func (a *TransactionUsecase) Delete(transactionID int64) error {
	err := a.transactionRepo.Delete(&models.Transaction{ID: transactionID})
	if err != nil {
		return err
	}

	return nil
}

func (a *TransactionUsecase) Index(paginationConfig request.TransactionPaginationConfig) ([]models.Transaction, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	transactions, err := a.transactionRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := a.transactionRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return transactions, meta, nil
}

func (a *TransactionUsecase) Show(transactionID int64) (models.Transaction, error) {
	return a.transactionRepo.FindByID(transactionID)
}
