package usecase

import (
	"github.com/PhantomX7/go-pos/auth"
	"github.com/PhantomX7/go-pos/auth/delivery/http/request"
	"github.com/PhantomX7/go-pos/auth/delivery/http/response"
	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/user"
	"github.com/PhantomX7/go-pos/role"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type AuthUsecase struct {
	userRepo user.UserRepository
	roleRepo role.RoleRepository
}

type authClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Id       uint64 `json:"id"`
}

func NewAuthUsecase(userRepo user.UserRepository, roleRepo role.RoleRepository) auth.AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (a AuthUsecase) SignIn(request request.SignInRequest) (response.AuthResponse, error) {
	userM, _ := a.userRepo.FindByUsername(request.Username)
	err := bcrypt.CompareHashAndPassword([]byte(userM.Password), []byte(request.Password))
	if err != nil {
		return response.AuthResponse{}, errors.ErrFailedAuthentication
	}

	// generate jwt token
	tokenString, err := generateTokenString(userM)
	if err != nil {
		return response.AuthResponse{}, err
	}

	return response.AuthResponse{
		Token:    tokenString,
		Username: userM.Username,
	}, nil

}

func (a AuthUsecase) SignUp(request request.SignUpRequest) (response.AuthResponse, error) {
	if request.Passphrase != os.Getenv("PASSPHRASE") {
		return response.AuthResponse{}, errors.ErrForbidden
	}

	userM := models.User{
		Username: request.Username,
		Password: request.Password,
	}
	err := a.userRepo.Insert(&userM)
	if err != nil {
		return response.AuthResponse{}, err
	}
	userM, _ = a.userRepo.FindByUsername(request.Username)
	// generate jwt token
	tokenString, err := generateTokenString(userM)
	if err != nil {
		return response.AuthResponse{}, err
	}

	return response.AuthResponse{
		Token:    tokenString,
		Username: userM.Username,
	}, nil
}

func (a AuthUsecase) GetMe(userID int64) (response.GetMeResponse, error) {
	userM, err := a.userRepo.FindByID(userID)
	if err != nil {
		return response.GetMeResponse{}, err
	}
	return response.GetMeResponse{
		Username: userM.Username,
	}, nil
}

func generateTokenString(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		Id:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 180).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("error-encoding-token", err)
		return "", errors.ErrUnprocessableEntity
	}

	return tokenString, nil
}
