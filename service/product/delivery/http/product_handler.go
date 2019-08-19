package http

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/service/product"
	"github.com/PhantomX7/go-pos/service/product/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase product.ProductUsecase
}

func NewProductHandler(productUC product.ProductUsecase) server.Handler {
	return &ProductHandler{
		productUsecase: productUC,
	}
}

func (h *ProductHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	productRoute := r.Group("/product", m.AuthHandle())
	{
		productRoute.GET("/", h.Index)
		productRoute.GET("/:id", h.Show)
		productRoute.POST("/", h.Create)
		productRoute.PUT("/:id", h.Update)
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req request.ProductCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	productModel, err := h.productUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, productModel)
}

func (h *ProductHandler) Update(c *gin.Context) {
	var req request.ProductUpdateRequest
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrUnprocessableEntity).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	productModel, err := h.productUsecase.Update(productID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, productModel)
}

func (h *ProductHandler) Index(c *gin.Context) {
	products, productPagination, err := h.productUsecase.Index(request.NewProductPaginationConfig(c.Request.URL.Query()))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.JSON(http.StatusOK, response_util.IndexResponse{
		Data: products,
		Meta: productPagination,
	})
}

func (h *ProductHandler) Show(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrUnprocessableEntity).SetType(gin.ErrorTypePublic)
	}

	productModel, err := h.productUsecase.Show(productID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, productModel)
}
