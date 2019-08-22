package http

import (
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/return_transaction"
	"github.com/PhantomX7/go-pos/service/return_transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/gin-gonic/gin"
)

type ReturnTransactionHandler struct {
	returnTransactionUsecase return_transaction.ReturnTransactionUsecase
	invoiceUsecase           invoice.InvoiceUsecase
}

func New(returnTransactionUC return_transaction.ReturnTransactionUsecase, invoiceUC invoice.InvoiceUsecase) server.Handler {
	return &ReturnTransactionHandler{
		returnTransactionUsecase: returnTransactionUC,
		invoiceUsecase:           invoiceUC,
	}
}

func (h *ReturnTransactionHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	returnTransactionRoute := r.Group("/returntransaction", m.AuthHandle())
	{
		//returnTransactionRoute.GET("/", h.Index)
		returnTransactionRoute.GET("/:id", h.Show)
		returnTransactionRoute.POST("/", h.Create)
		returnTransactionRoute.PUT("/:id", h.Update)
		returnTransactionRoute.DELETE("/:id", h.Delete)
	}
}

func (h *ReturnTransactionHandler) Create(c *gin.Context) {
	var req request.ReturnTransactionCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	// create returnTransaction
	returnTransactionModel, err := h.returnTransactionUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(returnTransactionModel.InvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, returnTransactionModel)
}

func (h *ReturnTransactionHandler) Update(c *gin.Context) {
	var req request.ReturnTransactionUpdateRequest
	returnTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	returnTransactionModel, err := h.returnTransactionUsecase.Update(returnTransactionID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(returnTransactionModel.InvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, returnTransactionModel)
}

func (h *ReturnTransactionHandler) Delete(c *gin.Context) {
	returnTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	returnTransactionModel, err := h.returnTransactionUsecase.Show(returnTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	err = h.returnTransactionUsecase.Delete(returnTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(returnTransactionModel.InvoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

//func (h *ReturnTransactionHandler) Index(c *gin.Context) {
//	returnTransactions, returnTransactionPagination, err := h.returnTransactionUsecase.Index(request.NewReturnTransactionPaginationConfig(c.Request.URL.Query()))
//	if err != nil {
//		_ = c.Error(err).SetType(gin.ErrorTypePublic)
//		return
//	}
//
//	c.JSON(http.StatusOK, map[string]interface{}{
//		"data": returnTransactions,
//		"meta": returnTransactionPagination,
//	})
//}

func (h *ReturnTransactionHandler) Show(c *gin.Context) {
	returnTransactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	returnTransactionModel, err := h.returnTransactionUsecase.Show(returnTransactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, returnTransactionModel)
}
