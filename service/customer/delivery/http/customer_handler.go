package http

import (
	"github.com/PhantomX7/go-pos/service/customer"
	"github.com/PhantomX7/go-pos/service/customer/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerUsecase customer.CustomerUsecase
}

func NewCustomerHandler(customerUC customer.CustomerUsecase) server.Handler {
	return &CustomerHandler{
		customerUsecase: customerUC,
	}
}

func (h *CustomerHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	customerRoute := r.Group("/customer", m.AuthHandle())
	{
		customerRoute.GET("/", h.Index)
		customerRoute.GET("/:id", h.Show)
		customerRoute.POST("/", h.Create)
		customerRoute.PUT("/:id", h.Update)
		customerRoute.DELETE("/:id", h.Delete)
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req request.CustomerCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	customerModel, err := h.customerUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, customerModel)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	var req request.CustomerUpdateRequest
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	customerModel, err := h.customerUsecase.Update(customerID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, customerModel)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	err = h.customerUsecase.Delete(customerID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

func (h *CustomerHandler) Index(c *gin.Context) {
	customers, customerPagination, err := h.customerUsecase.Index(request.NewCustomerPaginationConfig(c.Request.URL.Query()))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, response_util.IndexResponse{
		Data: customers,
		Meta: customerPagination,
	})
}

func (h *CustomerHandler) Show(c *gin.Context) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	customerModel, err := h.customerUsecase.Show(customerID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, customerModel)
}
