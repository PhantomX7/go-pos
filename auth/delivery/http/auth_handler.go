package http

import (
	"fmt"
	"github.com/PhantomX7/go-pos/app/api/middleware"
	"github.com/PhantomX7/go-pos/app/api/server"
	"github.com/PhantomX7/go-pos/auth"
	"github.com/PhantomX7/go-pos/auth/delivery/http/request"
	"github.com/PhantomX7/go-pos/utils/errors"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authUsecase auth.AuthUsecase
}

func NewAuthHandler(authUC auth.AuthUsecase) server.Handler {
	return &AuthHandler{
		authUsecase: authUC,
	}
}

func (h *AuthHandler) Register(r *gin.RouterGroup, m *middleware.Middleware) {
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/signin", h.SignIn)
		authRoute.POST("/signup", h.SignUp)
		authRoute.GET("/me", m.AuthHandle(), h.GetMe)
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req request.SignInRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	res, err := h.authUsecase.SignIn(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest

	// validate request
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	res, err := h.authUsecase.SignUp(req)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, err := strconv.ParseInt(fmt.Sprint(claims["id"]), 10, 64)
	if err != nil {
		log.Println("error-parsing-id-from-token:", err)
		_ = c.Error(errors.ErrUnprocessableEntity).SetType(gin.ErrorTypePublic)
		return
	}
	res, err := h.authUsecase.GetMe(userID)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(http.StatusOK, res)
}
