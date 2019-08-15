package usecase

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/product"
	"github.com/PhantomX7/go-pos/product/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type ProductUsecase struct {
	productRepo product.ProductRepository
}

func NewProductUsecase(productRepo product.ProductRepository) product.ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

func (a *ProductUsecase) Create(request request.ProductCreateRequest) (models.Product, error) {
	productM := models.Product{
		Name:            request.Name,
		Pinyin:          request.Pinyin,
		Description:     request.Description,
		Stock:           request.Stock,
		Unit:            request.Unit,
		UnitAmount:      request.UnitAmount,
		CapitalPrice:    request.CapitalPrice,
		SellPriceCash:   request.SellPriceCash,
		SellPriceCredit: request.SellPriceCredit,
	}
	err := a.productRepo.Insert(&productM)
	if err != nil {
		return productM, err
	}
	return productM, nil
}

func (a *ProductUsecase) Update(productID int64, request request.ProductUpdateRequest) (models.Product, error) {
	productM, err := a.productRepo.FindByID(productID)
	if err != nil {
		return productM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&productM, &request)

	err = a.productRepo.Update(&productM)
	if err != nil {
		return productM, err
	}
	return productM, nil
}

func (a *ProductUsecase) Index(paginationConfig request.ProductPaginationConfig) ([]models.Product, response_util.PaginationMeta, error) {
	meta := response_util.PaginationMeta{
		Offset: paginationConfig.Offset(),
		Limit:  paginationConfig.Limit(),
		Total:  0,
	}
	products, err := a.productRepo.FindAll(paginationConfig)
	if err != nil {
		return nil, meta, err
	}
	total, err := a.productRepo.Count(paginationConfig)
	if err != nil {
		return nil, meta, err
	}
	meta.Total = total

	return products, meta, nil
}

func (a *ProductUsecase) Show(productID int64) (models.Product, error) {
	return a.productRepo.FindByID(productID)
}
