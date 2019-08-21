package http

import (
	"github.com/PhantomX7/go-pos/service/invoice"
	"github.com/PhantomX7/go-pos/service/transaction"
	"github.com/PhantomX7/go-pos/service/transaction/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase transaction.TransactionUsecase
	invoiceUsecase     invoice.InvoiceUsecase
}

func NewTransactionHandler(transactionUC transaction.TransactionUsecase, invoiceUC invoice.InvoiceUsecase) server.Handler {
	return &TransactionHandler{
		transactionUsecase: transactionUC,
		invoiceUsecase:     invoiceUC,
	}
}

func (h *TransactionHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	transactionRoute := r.Group("/transaction", m.AuthHandle())
	{
		//transactionRoute.GET("/", h.Index)
		transactionRoute.GET("/:id", h.Show)
		transactionRoute.POST("/", h.Create)
		transactionRoute.PUT("/:id", h.Update)
		transactionRoute.DELETE("/:id", h.Delete)
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req request.TransactionCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	// create transaction
	transactionModel, err := h.transactionUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(transactionModel.InvoiceId)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, transactionModel)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	var req request.TransactionUpdateRequest
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	transactionModel, err := h.transactionUsecase.Update(transactionID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(transactionModel.InvoiceId)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, transactionModel)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	transactionModel, err := h.transactionUsecase.Show(transactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	err = h.transactionUsecase.Delete(transactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	// sync invoice
	err = h.invoiceUsecase.SyncInvoice(transactionModel.InvoiceId)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

//func (h *TransactionHandler) Index(c *gin.Context) {
//	transactions, transactionPagination, err := h.transactionUsecase.Index(request.NewTransactionPaginationConfig(c.Request.URL.Query()))
//	if err != nil {
//		_ = c.Error(err).SetType(gin.ErrorTypePublic)
//		return
//	}
//
//	c.JSON(http.StatusOK, map[string]interface{}{
//		"data": transactions,
//		"meta": transactionPagination,
//	})
//}

func (h *TransactionHandler) Show(c *gin.Context) {
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	transactionModel, err := h.transactionUsecase.Show(transactionID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, transactionModel)
}
