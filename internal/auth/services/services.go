package services

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	authrpc "github.com/morzhanov/go-realworld/api/grpc/auth"
	usersrpc "github.com/morzhanov/go-realworld/api/grpc/users"
	aconfig "github.com/morzhanov/go-realworld/internal/auth/config"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
)

type authService struct {
	sender            sender.Sender
	el                eventslistener.EventListener
	accessTokenSecret string
}

type AuthService interface {
	Login(ctx context.Context, data *authrpc.LoginInput, span *opentracing.Span) (res *authrpc.AuthResponse, err error)
	Signup(ctx context.Context, data *authrpc.SignupInput, span *opentracing.Span) (res *authrpc.AuthResponse, err error)
	ValidateRestRequest(data *authrpc.ValidateRestRequestInput) (res *authrpc.ValidationResponse, err error)
	ValidateRpcRequest(data *authrpc.ValidateRpcRequestInput) (res *authrpc.ValidationResponse, err error)
	ValidateEventsRequest(data *authrpc.ValidateEventsRequestInput) (res *authrpc.ValidationResponse, err error)
}

func (s *authService) createJwt(userId string) (res string, err error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(s.accessTokenSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) verifyJwt(tokenString string) (res *authrpc.ValidationResponse, err error) {
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return &authrpc.ValidationResponse{UserId: claims["user_id"].(string)}, nil
}

func (s *authService) Login(
	ctx context.Context,
	data *authrpc.LoginInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	transport := s.sender.GetTransportFromContext(ctx)

	d := usersrpc.ValidateUserPasswordRequest{Username: data.Username, Password: data.Password}
	err = s.sender.PerformRequest(transport, "users", "validatePassword", &d, s.el, span, nil, nil)
	if err != nil {
		return nil, err
	}

	d2 := usersrpc.GetUserDataByUsernameRequest{Username: data.Username}
	user := &usersrpc.UserMessage{}
	queryparams := fmt.Sprintf("username=%s", data.Username)
	meta := sender.RequestMeta{"queryparams": queryparams}
	if err = s.sender.PerformRequest(transport, "users", "getUserByUsername", &d2, s.el, span, meta, user); err != nil {
		return nil, err
	}

	token, err := s.createJwt(user.Id)
	if err != nil {
		return nil, err
	}
	return &authrpc.AuthResponse{AccessToken: token}, nil
}

func (s *authService) Signup(
	ctx context.Context,
	data *authrpc.SignupInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	transport := s.sender.GetTransportFromContext(ctx)

	d := usersrpc.CreateUserRequest{Username: data.Username, Password: data.Password}
	user := &usersrpc.UserMessage{}
	if err = s.sender.PerformRequest(transport, "users", "createUser", &d, s.el, span, nil, &user); err != nil {
		return nil, err
	}

	token, err := s.createJwt(user.Id)
	if err != nil {
		return nil, err
	}
	return &authrpc.AuthResponse{AccessToken: token}, nil
}

func (s *authService) ValidateRestRequest(data *authrpc.ValidateRestRequestInput) (res *authrpc.ValidationResponse, err error) {
	for _, route := range aconfig.PUBLIC_ROUTES {
		if route == data.Path {
			return nil, nil
		}
	}
	return s.verifyJwt(data.AccessToken)
}

func (s *authService) ValidateRpcRequest(data *authrpc.ValidateRpcRequestInput) (res *authrpc.ValidationResponse, err error) {
	for _, route := range aconfig.PUBLIC_RPC_METHODS {
		if route == data.Method {
			return nil, nil
		}
	}
	return s.verifyJwt(data.AccessToken)
}

func (s *authService) ValidateEventsRequest(data *authrpc.ValidateEventsRequestInput) (res *authrpc.ValidationResponse, err error) {
	return s.verifyJwt(data.AccessToken)
}

func NewAuthService(
	s sender.Sender,
	el eventslistener.EventListener,
	c *config.Config,
) AuthService {
	return &authService{s, el, c.AccessTokenSecret}
}
