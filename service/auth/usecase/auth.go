package usecase

import (
	"log"
	"os"
	"time"

	"github.com/PhantomX7/go-pos/models"
	"github.com/PhantomX7/go-pos/service/auth"
	"github.com/PhantomX7/go-pos/service/auth/delivery/http/request"
	"github.com/PhantomX7/go-pos/service/auth/delivery/http/response"
	"github.com/PhantomX7/go-pos/service/role"
	"github.com/PhantomX7/go-pos/service/user"
	"github.com/PhantomX7/go-pos/utils/errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo user.UserRepository
	roleRepo role.RoleRepository
}

type authClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
	Id       uint64 `json:"id"`
}

func New(userRepo user.UserRepository, roleRepo role.RoleRepository) auth.AuthUsecase {
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

	roleM, _ := a.roleRepo.FindByID(userM.RoleId)

	// generate jwt token
	tokenString, err := generateTokenString(userM, roleM)
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

	err := a.userRepo.Insert(&models.User{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		return response.AuthResponse{}, err
	}
	userM, _ := a.userRepo.FindByUsername(request.Username)
	roleM, _ := a.roleRepo.FindByID(userM.RoleId)
	// generate jwt token
	tokenString, err := generateTokenString(userM, roleM)
	if err != nil {
		return response.AuthResponse{}, err
	}

	return response.AuthResponse{
		Token:    tokenString,
		Username: userM.Username,
	}, nil
}

func (a AuthUsecase) GetMe(userID uint64) (response.GetMeResponse, error) {
	userM, err := a.userRepo.FindByID(userID)
	if err != nil {
		return response.GetMeResponse{}, err
	}
	roleM, err := a.roleRepo.FindByID(userM.RoleId)
	if err != nil {
		return response.GetMeResponse{}, err
	}

	return response.GetMeResponse{
		Username: userM.Username,
		Role:     roleM.Name,
	}, nil
}

func generateTokenString(user *models.User, role *models.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		Id:       user.ID,
		Username: user.Username,
		Role:     role.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 180).Unix(),
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
