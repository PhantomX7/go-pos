package http

import (
	"github.com/PhantomX7/go-pos/service/order_invoice"
	"github.com/PhantomX7/go-pos/service/order_transaction"
	"github.com/PhantomX7/go-pos/service/order_transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/gin-gonic/gin"
)

type OrderTransactionHandler struct {
	orderTransactionUsecase order_transaction.OrderTransactionUsecase
	orderInvoiceUsecase     order_invoice.OrderInvoiceUsecase
}

func New(orderTransactionUC order_transaction.OrderTransactionUsecase, orderInvoiceUC order_invoice.OrderInvoiceUsecase) server.Handler {
	return &OrderTransactionHandler{
		orderTransactionUsecase: orderTransactionUC,
		orderInvoiceUsecase:     orderInvoiceUC,
	}
}

func (h *OrderTransactionHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	orderTransactionRoute := r.Group("/order_transaction", m.AuthHandle())
	{
		//orderTransactionRoute.GET("/", h.Index)
		orderTransactionRoute.GET("/:id", h.Show)
		orderTransactionRoute.POST("/", h.Create)
		orderTransactionRoute.PUT("/:id", h.Update)
		orderTransactionRoute.DELETE("/:id", h.Delete)
	}
}

func (h *OrderTransactionHandler) Create(c *gin.Context) {
	var req request.OrderTransactionCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	// create orderTransaction
	orderTransactionModel, err := h.orderTransactionUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync order invoice
	err = h.orderInvoiceUsecase.SyncOrderInvoice(orderTransactionModel.OrderInvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderTransactionModel)
}

func (h *OrderTransactionHandler) Update(c *gin.Context) {
	var req request.OrderTransactionUpdateRequest
	orderTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	orderTransactionModel, err := h.orderTransactionUsecase.Update(orderTransactionID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync order invoice
	err = h.orderInvoiceUsecase.SyncOrderInvoice(orderTransactionModel.OrderInvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderTransactionModel)
}

func (h *OrderTransactionHandler) Delete(c *gin.Context) {
	orderTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	orderTransactionModel, err := h.orderTransactionUsecase.Show(orderTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	err = h.orderTransactionUsecase.Delete(orderTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync order invoice
	err = h.orderInvoiceUsecase.SyncOrderInvoice(orderTransactionModel.OrderInvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

//func (h *OrderTransactionHandler) Index(c *gin.Context) {
//	orderTransactions, orderTransactionPagination, err := h.orderTransactionUsecase.Index(request.NewOrderTransactionPaginationConfig(c.Request.URL.Query()))
//	if err != nil {
//		_ = c.Error(err).SetType(gin.ErrorTypePublic)
//		return
//	}
//
//	c.JSON(http.StatusOK, map[string]interface{}{
//		"data": orderTransactions,
//		"meta": orderTransactionPagination,
//	})
//}

func (h *OrderTransactionHandler) Show(c *gin.Context) {
	orderTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	orderTransactionModel, err := h.orderTransactionUsecase.Show(orderTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, orderTransactionModel)
}
