package auth

import (
	"github.com/PhantomX7/go-pos/service/auth/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/auth/delivery/http/response"
)

type AuthUsecase interface {
	SignIn(request request.SignInRequest) (response.AuthResponse, error)
	SignUp(request request.SignUpRequest) (response.AuthResponse, error)
	GetMe(userID uint64) (response.GetMeResponse, error)
}
