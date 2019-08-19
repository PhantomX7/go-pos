package product

import (
	"github.com/PhantomX7/go-pos/service/product/delivery/http/request"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/utils/response_util"
)

type ProductUsecase interface {
	Create(request request.ProductCreateRequest) (models.Product, error)
	Update(productID uint64, request request.ProductUpdateRequest) (models.Product, error)
	Index(paginationConfig request.ProductPaginationConfig) ([]models.Product, response_util.PaginationMeta, error)
	Show(productID uint64) (models.Product, error)
}
