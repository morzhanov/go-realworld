package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/grpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/grpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
)

type APIGatewayService struct {
	sender        *sender.Sender
	eventListener *eventslistener.EventListener
}

func (s *APIGatewayService) getAccessToken(ctx *gin.Context) (string, error) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("not authorized")
	}
	return authorization[7:], nil
}

func (s *APIGatewayService) CheckAuth(
	ctx *gin.Context,
	transport sender.Transport,
	apiName string,
	key string,
	span *opentracing.Span,
) (res *authrpc.ValidationResponse, err error) {
	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var input interface{}
	var method string
	switch transport {
	case sender.RestTransport:
		api, err := s.sender.API.GetApiItem(apiName)
		if err != nil {
			return nil, err
		}
		input = &authrpc.ValidateRestRequestInput{
			Path:        api.Rest[key].Url,
			AccessToken: accessToken,
		}
		method = "validateRestRequest"
	case sender.RpcTransport:
		input = &authrpc.ValidateRpcRequestInput{
			Method:      key,
			AccessToken: accessToken,
		}
		method = "validateRpcRequest"
	case sender.EventsTransport:
		api, err := s.sender.API.GetApiItem(apiName)
		if err != nil {
			return nil, err
		}
		input = &authrpc.ValidateEventsRequestInput{
			Event:       api.Events[key].Event,
			AccessToken: accessToken,
		}
		method = "validateEventsRequest"
	default:
		return nil, fmt.Errorf("not valid transport %v", transport)
	}
	if err = s.sender.PerformRequest(transport, "auth", method, input, s.eventListener, span, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func createMetaWithUserId(userId string) sender.RequestMeta {
	return sender.RequestMeta{"urlparams": sender.UrlParams{
		"userId": userId,
	}}
}

func (s *APIGatewayService) Login(
	transport sender.Transport,
	input *authrpc.LoginInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	res = &authrpc.AuthResponse{}
	if err := s.sender.PerformRequest(transport, "auth", "login", input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) Signup(
	transport sender.Transport,
	input *authrpc.SignupInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	res = &authrpc.AuthResponse{}
	if err = s.sender.PerformRequest(transport, "auth", "signup", input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) GetPictures(
	transport sender.Transport,
	userId string,
	span *opentracing.Span,
) (res *prpc.PicturesMessage, err error) {
	input := prpc.GetUserPicturesRequest{UserId: userId}
	res = &prpc.PicturesMessage{}
	meta := createMetaWithUserId(userId)
	if err := s.sender.PerformRequest(transport, "pictures", "getPictures", &input, s.eventListener, span, meta, &res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) GetPicture(
	transport sender.Transport,
	userId string,
	pictureId string,
	span *opentracing.Span,
) (res *prpc.PictureMessage, err error) {
	input := prpc.GetUserPictureRequest{
		UserId:    userId,
		PictureId: pictureId,
	}
	meta := createMetaWithUserId(userId)
	urlparams := meta["urlparams"].(sender.UrlParams)
	urlparams["id"] = pictureId
	res = &prpc.PictureMessage{}
	if err := s.sender.PerformRequest(transport, "pictures", "getPicture", &input, s.eventListener, span, meta, &res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) CreatePicture(
	transport sender.Transport,
	input *prpc.CreateUserPictureRequest,
	span *opentracing.Span,
) (res *prpc.PictureMessage, err error) {
	res = &prpc.PictureMessage{}
	meta := createMetaWithUserId(input.UserId)
	if err := s.sender.PerformRequest(transport, "pictures", "createPicture", input, s.eventListener, span, meta, &res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) DeletePicture(
	transport sender.Transport,
	userId string,
	pictureId string,
	span *opentracing.Span,
) error {
	input := prpc.DeleteUserPictureRequest{
		UserId:    userId,
		PictureId: pictureId,
	}
	meta := createMetaWithUserId(userId)
	return s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span, meta, nil)
}

func (s *APIGatewayService) GetAnalytics(
	transport sender.Transport,
	input *anrpc.GetLogRequest,
	span *opentracing.Span,
) (res *anrpc.AnalyticsEntryMessage, err error) {
	res = &anrpc.AnalyticsEntryMessage{}
	if err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span, nil, &res); err != nil {
		return nil, err
	}
	return
}

func NewAPIGatewayService(s *sender.Sender, el *eventslistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
