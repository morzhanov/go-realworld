package services

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	. "github.com/morzhanov/go-realworld/internal/auth/config"
	. "github.com/morzhanov/go-realworld/internal/auth/dto"
	. "github.com/morzhanov/go-realworld/internal/users/models"
)

// TODO: kuber should proxy all rest/grpc/events requests to this service
// TODO: service should handle access tokens from rest/grpc/events calls and check them
// TODO: also service should handle login and signup calls
type AuthService struct {
	db *sqlx.DB
}

func (s *AuthService) Login(data *LoginInput) (res *LoginDto, err error) {
	// TODO: validate user password (users.ValidateUserPassword)
	// TODO: get user data by username from users service (users.GetUserDataByUsername)
	user := User{}

	token, err := createJwt(user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginDto{AccessToken: token}, nil
}

func (s *AuthService) Signup(data *SignupInput) (res *LoginDto, err error) {
	// TODO: create new user via users.CreateUser
	user := User{}

	token, err := createJwt(user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginDto{AccessToken: token}, nil
}

func createJwt(userId string) (res string, err error) {
	const secret = "jdnfksdmfksd"

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func verifyJwt(tokenString string) error {
	const secret = "jdnfksdmfksd"

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: this GET endpoint should be called from ingress
func (s *AuthService) ValidateRestRequest(data *ValidateRestRequestInput) error {
	for _, route := range PUBLIC_ROUTES {
		if route == data.Path {
			return nil
		}
	}
	return verifyJwt(data.AccessToken)
	// TODO: inject userId into headers and return result to ingress
}

// TODO: this GET endpoint should be called from client service
func (s *AuthService) ValidateRpcRequest(data *ValidateRpcRequestInput) error {
	for _, route := range PUBLIC_RPC_METHODS {
		if route == data.Method {
			return nil
		}
	}
	return verifyJwt(data.AccessToken)
	// TODO: return userId
}

// TODO: this GET endpoint should be called from client service
func (s *AuthService) ValidateEventsRequest(data *ValidateEventsRequestInput) error {
	for _, route := range PUBLIC_EVENTS {
		if route == data.Event {
			return nil
		}
	}
	return verifyJwt(data.AccessToken)
	// TODO: return userId
}

func NewAuthService(db *sqlx.DB) *AuthService {
	return &AuthService{db}
}
