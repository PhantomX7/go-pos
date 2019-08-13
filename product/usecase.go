package product

import (
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/product/delivery/http/request"
	"github.com/PhantomX7/go-pos/product/delivery/http/response"
)

type ProductUsecase interface {
	Create(request request.ProductCreateRequest) (models.Product, error)
	Update(productID int64, request request.ProductUpdateRequest) (models.Product, error)
	Index(paginationConfig request.ProductPaginationConfig) ([]models.Product, response.ProductPaginationMeta, error)
	Show(productID int64) (models.Product, error)
}
