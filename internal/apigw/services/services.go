package services

import (
	"fmt"

	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	pmodel "github.com/morzhanov/go-realworld/internal/pictures/models"
	"github.com/opentracing/opentracing-go"
)

type APIGatewayService struct {
	sender        *sender.Sender
	eventListener *eventslistener.EventListener
}

func (s *APIGatewayService) getAccessToken(ctx *gin.Context) string {
	authorization := ctx.GetHeader("Authorization")
	return authorization[6:]
}

func (s *APIGatewayService) CheckAuth(
	ctx *gin.Context,
	transport sender.Transport,
	apiName string,
	key string,
	span *opentracing.Span,
) (res *authrpc.ValidationResponse, err error) {
	accessToken := s.getAccessToken(ctx)

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
		method = "verifyRestRequest"
	case sender.RpcTransport:
		input = &authrpc.ValidateRpcRequestInput{
			Method:      key,
			AccessToken: accessToken,
		}
		method = "verifyRpcRequest"
	case sender.EventsTransport:
		api, err := s.sender.API.GetApiItem(apiName)
		if err != nil {
			return nil, err
		}
		input = &authrpc.ValidateEventsRequestInput{
			Event:       api.Events[key].Event,
			AccessToken: accessToken,
		}
		method = "verifyEventsRequest"
	default:
		return nil, fmt.Errorf("not valid transport %v", transport)
	}
	res = &authrpc.ValidationResponse{}
	if err = s.sender.PerformRequest(transport, "auth", method, input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return res, nil
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
) (res []*pmodel.Picture, err error) {
	input := prpc.GetUserPicturesRequest{UserId: userId}
	res = []*pmodel.Picture{}
	if err := s.sender.PerformRequest(transport, "pictures", "getPictures", &input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) GetPicture(
	transport sender.Transport,
	userId string,
	pictureId string,
	span *opentracing.Span,
) (res *pmodel.Picture, err error) {
	input := prpc.GetUserPictureRequest{
		UserId:    userId,
		PictureId: pictureId,
	}
	res = &pmodel.Picture{}
	if err := s.sender.PerformRequest(transport, "pictures", "getPicture", &input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return
}

func (s *APIGatewayService) CreatePicture(
	transport sender.Transport,
	input *prpc.CreateUserPictureRequest,
	span *opentracing.Span,
) (res *pmodel.Picture, err error) {
	res = &pmodel.Picture{}
	if err := s.sender.PerformRequest(transport, "pictures", "createPicture", &input, s.eventListener, span, nil, res); err != nil {
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
	return s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span, nil, nil)
}

func (s *APIGatewayService) GetAnalytics(
	transport sender.Transport,
	input *anrpc.GetLogRequest,
	span *opentracing.Span,
) (res *anrpc.AnalyticsEntryMessage, err error) {
	res = &anrpc.AnalyticsEntryMessage{}
	if err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span, nil, res); err != nil {
		return nil, err
	}
	return
}

func NewAPIGatewayService(s *sender.Sender, el *eventslistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
