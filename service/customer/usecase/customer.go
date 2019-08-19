package usecase

import (
	"github.com/PhantomX7/go-pos/service/customer"
	"github.com/PhantomX7/go-pos/service/customer/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type CustomerUsecase struct {
	customerRepo customer.CustomerRepository
}

func NewCustomerUsecase(customerRepo customer.CustomerRepository) customer.CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
	}
}

func (a *CustomerUsecase) Create(request request.CustomerCreateRequest) (models.Customer, error) {
	customerM := models.Customer{
		Name:    request.Name,
		Address: request.Address,
		Phone:   request.Phone,
	}
	err := a.customerRepo.Insert(&customerM)
	if err != nil {
		return customerM, err
	}
	return customerM, nil
}

func (a *CustomerUsecase) Update(customerID int64, request request.CustomerUpdateRequest) (models.Customer, error) {
	customerM, err := a.customerRepo.FindByID(customerID)
	if err != nil {
		return customerM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&customerM, &request)

	err = a.customerRepo.Update(&customerM)
	if err != nil {
		return customerM, err
	}
	return customerM, nil
}

func (a *CustomerUsecase) Delete(customerID int64) error {
	err := a.customerRepo.Delete(&models.Customer{ID: customerID})
	if err != nil {
		return err
	}

	return nil
}

func (a *CustomerUsecase) Index(paginationConfig request.CustomerPaginationConfig) ([]models.Customer, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}

	customers, err := a.customerRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	total, err := a.customerRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}

	meta.Total = total

	return customers, meta, nil
}

func (a *CustomerUsecase) Show(customerID int64) (models.Customer, error) {
	return a.customerRepo.FindByID(customerID)
}
