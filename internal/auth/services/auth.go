package services

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
	aconfig "github.com/morzhanov/go-realworld/internal/auth/config"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
)

type AuthService struct {
	db                *sqlx.DB
	sender            *sender.Sender
	el                *eventslistener.EventListener
	accessTokenSecret string
}

func getTransport(ctx context.Context) sender.Transport {
	val := ctx.Value("transport")
	return val.(sender.Transport)
}

func (s *AuthService) Login(
	ctx context.Context,
	data *authrpc.LoginInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	transport := getTransport(ctx)

	d := usersrpc.ValidateUserPasswordRequest{Username: data.Username, Password: data.Password}
	_, err = s.sender.PerformRequest(transport, "users", "validatePassword", &d, s.el, span)
	if err != nil {
		return nil, err
	}

	d2 := usersrpc.GetUserDataByUsernameRequest{Username: data.Username}
	r, err := s.sender.PerformRequest(transport, "users", "getUserByUsername", &d2, s.el, span)
	if err != nil {
		return nil, err
	}
	user := r.(*usersrpc.UserMessage)

	token, err := createJwt(user.Id)
	if err != nil {
		return nil, err
	}
	return &authrpc.AuthResponse{AccessToken: token}, nil
}

func (s *AuthService) Signup(
	ctx context.Context,
	data *authrpc.SignupInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	transport := getTransport(ctx)

	d := usersrpc.CreateUserRequest{Username: data.Username, Password: data.Password}
	r, err := s.sender.PerformRequest(transport, "users", "createUser", &d, s.el, span)
	if err != nil {
		return nil, err
	}
	user := r.(*usersrpc.UserMessage)

	token, err := createJwt(user.Id)
	if err != nil {
		return nil, err
	}
	return &authrpc.AuthResponse{AccessToken: token}, nil
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

func (s *AuthService) verifyJwt(tokenString string) (res *authrpc.ValidationResponse, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return &authrpc.ValidationResponse{UserId: token.Raw}, nil
}

func (s *AuthService) ValidateRestRequest(data *authrpc.ValidateRestRequestInput) (res *authrpc.ValidationResponse, err error) {
	for _, route := range aconfig.PUBLIC_ROUTES {
		if route == data.Path {
			return nil, nil
		}
	}
	return s.verifyJwt(data.AccessToken)
}

func (s *AuthService) ValidateRpcRequest(data *authrpc.ValidateRpcRequestInput) (res *authrpc.ValidationResponse, err error) {
	for _, route := range aconfig.PUBLIC_RPC_METHODS {
		if route == data.Method {
			return nil, nil
		}
	}
	return s.verifyJwt(data.AccessToken)
}

func (s *AuthService) ValidateEventsRequest(data *authrpc.ValidateEventsRequestInput) (res *authrpc.ValidationResponse, err error) {
	for _, route := range aconfig.PUBLIC_EVENTS {
		if route == data.Event {
			return nil, nil
		}
	}
	return s.verifyJwt(data.AccessToken)
}

func NewAuthService(
	db *sqlx.DB,
	s *sender.Sender,
	el *eventslistener.EventListener,
	c *config.Config,
) *AuthService {
	return &AuthService{db, s, el, c.AccessTokenSecret}
}
