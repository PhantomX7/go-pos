package http

import (
	"github.com/PhantomX7/go-pos/utils/errors"
	"net/http"
	"strconv"

	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/invoice"
	"github.com/PhantomX7/go-pos/invoice/delivery/http/request"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	invoiceUsecase invoice.InvoiceUsecase
}

func NewInvoiceHandler(invoiceUC invoice.InvoiceUsecase) server.Handler {
	return &InvoiceHandler{
		invoiceUsecase: invoiceUC,
	}
}

func (h *InvoiceHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	invoiceRoute := r.Group("/invoice", m.AuthHandle())
	{
		invoiceRoute.GET("/", h.Index)
		invoiceRoute.GET("/:id", h.Show)
		invoiceRoute.POST("/", h.Create)
		invoiceRoute.PUT("/:id", h.Update)
		invoiceRoute.DELETE("/:id", h.Delete)
	}
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	var req request.InvoiceCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	invoiceModel, err := h.invoiceUsecase.Create(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, invoiceModel)
}

func (h *InvoiceHandler) Update(c *gin.Context) {
	var req request.InvoiceUpdateRequest
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	invoiceModel, err := h.invoiceUsecase.Update(invoiceID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, invoiceModel)
}

func (h *InvoiceHandler) Delete(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	err = h.invoiceUsecase.Delete(invoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.Status(http.StatusOK)
}

func (h *InvoiceHandler) Index(c *gin.Context) {
	invoices, invoicePagination, err := h.invoiceUsecase.Index(request.NewInvoicePaginationConfig(c.Request.URL.Query()))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data": invoices,
		"meta": invoicePagination,
	})
}

func (h *InvoiceHandler) Show(c *gin.Context) {
	invoiceID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrNotFound).SetType(gin.ErrorTypePublic)
	}

	invoiceModel, err := h.invoiceUsecase.Show(invoiceID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, invoiceModel)
}
