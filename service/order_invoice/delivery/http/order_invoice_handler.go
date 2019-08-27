package http

import (
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/service/order_invoice/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/utils/response_util"
	"github.com/gin-gonic/gin"
)

type OrderInvoiceHandler struct {
	orderInvoiceUsecase order_invoice.OrderInvoiceUsecase
}

func New(orderInvoiceUC order_invoice.OrderInvoiceUsecase) server.Handler {
	return &OrderInvoiceHandler{
		orderInvoiceUsecase: orderInvoiceUC,
	}
}

func (h *OrderInvoiceHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	orderInvoiceRoute := r.Group("/order_invoice", m.AuthHandle())
	{
		orderInvoiceRoute.GET("/", h.Index)
		orderInvoiceRoute.GET("/:id", h.Show)
		orderInvoiceRoute.POST("/", h.Create)
		orderInvoiceRoute.PUT("/:id", h.Update)
		orderInvoiceRoute.DELETE("/:id", h.Delete)
	}
}

func (h *OrderInvoiceHandler) Create(c *gin.Context) {
	var req request.OrderInvoiceCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	orderInvoiceModel, err := h.orderInvoiceUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderInvoiceModel)
}

func (h *OrderInvoiceHandler) Update(c *gin.Context) {
	var req request.OrderInvoiceUpdateRequest
	orderInvoiceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	orderInvoiceModel, err := h.orderInvoiceUsecase.Update(orderInvoiceID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderInvoiceModel)
}

func (h *OrderInvoiceHandler) Delete(c *gin.Context) {
	orderInvoiceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	err = h.orderInvoiceUsecase.Delete(orderInvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderInvoiceHandler) Index(c *gin.Context) {
	orderInvoices, orderInvoicePagination, err := h.orderInvoiceUsecase.Index(request.NewOrderInvoicePaginationConfig(c.Request.URL.Query()))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, response_util.IndexResponse{
		Data: orderInvoices,
		Meta: orderInvoicePagination,
	})
}

func (h *OrderInvoiceHandler) Show(c *gin.Context) {
	orderInvoiceID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	orderInvoiceModel, err := h.orderInvoiceUsecase.Show(orderInvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderInvoiceModel)
}
