package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/user"
	"github.com/PhantomX7/go-pos/user/delivery/http/request"
)

type UserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserHandler(userUC user.UserUsecase) server.Handler {
	return &UserHandler{
		userUsecase: userUC,
	}
}

func (h *UserHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {

	userRoute := r.Group("/user", m.AuthHandle())
	{
		userRoute.GET("/", h.Index)
		userRoute.GET("/:id", h.Show)
		userRoute.POST("/", h.Create)
		userRoute.PUT("/:id", h.Update)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req request.UserCreateRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	userModel, err := h.userUsecase.Create(req.ToUserModel())
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, userModel)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req request.UserUpdateRequest
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrUnprocessableEntity).SetType(gin.ErrorTypePublic)
	}

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	userModel, err := h.userUsecase.Update(userID, req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, userModel)
}

func (h *UserHandler) Index(c *gin.Context) {
	users, userPagination, err := h.userUsecase.Index(request.NewUserPaginationConfig(c.Request.URL.Query()))
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, map[string]interface{}{
		"data": users,
		"meta": userPagination,
	})
}

func (h *UserHandler) Show(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errors.ErrUnprocessableEntity).SetType(gin.ErrorTypePublic)
	}

	userModel, err := h.userUsecase.Show(userID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, userModel)
}
