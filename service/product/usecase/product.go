package usecase

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/product/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/stock_adjustment"
	"github.com/PhantomX7/go-pos/utils/database"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/jinzhu/copier"
)

// apply business logic here

type ProductUsecase struct {
	productRepo         product.ProductRepository
	stockAdjustmentRepo stock_adjustment.StockAdjustmentRepository
}

func New(
	productRepo product.ProductRepository,
	stockAdjustmentRepo stock_adjustment.StockAdjustmentRepository,
) product.ProductUsecase {
	return &ProductUsecase{
		productRepo:         productRepo,
		stockAdjustmentRepo: stockAdjustmentRepo,
	}
}

func (a *ProductUsecase) Create(request request.ProductCreateRequest) (*models.Product, error) {
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
	tx := database.BeginTransactions()

	err := a.productRepo.Insert(&productM, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = a.stockAdjustmentRepo.Insert(&models.StockAdjustment{
		ProductID:   productM.ID,
		Description: "new product",
		Amount:      productM.Stock,
	}, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &productM, nil
}

func (a *ProductUsecase) Update(productID uint64, request request.ProductUpdateRequest) (*models.Product, error) {
	productM, err := a.productRepo.FindByID(productID)
	if err != nil {
		return productM, err
	}

	// copy content of request into request model found by id
	_ = copier.Copy(&productM, &request)

	tx := database.BeginTransactions()

	if request.Stock != nil {
		if *request.Stock != productM.Stock {
			err = a.stockAdjustmentRepo.Insert(&models.StockAdjustment{
				ProductID:   productM.ID,
				Description: request.StockDescription,
				Amount:      productM.Stock - *request.Stock,
			}, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// update product
	err = a.productRepo.Update(productM, tx)
	if err != nil {
		tx.Rollback()
		return productM, err
	}
	tx.Commit()

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

func (a *ProductUsecase) Show(productID uint64) (*models.Product, error) {
	return a.productRepo.FindByID(productID)
}
