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
	result, err := s.sender.PerformRequest(transport, "auth", method, input, s.eventListener, span)
	return result.(*authrpc.ValidationResponse), err
}

func (s *APIGatewayService) Login(
	transport sender.Transport,
	input *authrpc.LoginInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	result, err := s.sender.PerformRequest(transport, "auth", "login", input, s.eventListener, span)
	return result.(*authrpc.AuthResponse), err
}

func (s *APIGatewayService) Signup(
	transport sender.Transport,
	input *authrpc.SignupInput,
	span *opentracing.Span,
) (res *authrpc.AuthResponse, err error) {
	result, err := s.sender.PerformRequest(transport, "auth", "signup", input, s.eventListener, span)
	return result.(*authrpc.AuthResponse), err
}

func (s *APIGatewayService) GetPictures(
	transport sender.Transport,
	userId string,
	span *opentracing.Span,
) (res []*pmodel.Picture, err error) {
	input := prpc.GetUserPicturesRequest{UserId: userId}
	result, err := s.sender.PerformRequest(transport, "pictures", "getPictures", &input, s.eventListener, span)
	return result.([]*pmodel.Picture), err
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
	result, err := s.sender.PerformRequest(transport, "pictures", "getPicture", &input, s.eventListener, span)
	return result.(*pmodel.Picture), err
}

func (s *APIGatewayService) CreatePicture(
	transport sender.Transport,
	input *prpc.CreateUserPictureRequest,
	span *opentracing.Span,
) (res *pmodel.Picture, err error) {
	result, err := s.sender.PerformRequest(transport, "pictures", "createPicture", &input, s.eventListener, span)
	return result.(*pmodel.Picture), err
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
	_, err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span)
	return err
}

func (s *APIGatewayService) GetAnalytics(
	transport sender.Transport,
	input *anrpc.GetLogRequest,
	span *opentracing.Span,
) (res *anrpc.AnalyticsEntryMessage, err error) {
	result, err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener, span)
	return result.(*anrpc.AnalyticsEntryMessage), err
}

func NewAPIGatewayService(s *sender.Sender, el *eventslistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
