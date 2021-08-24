package services

import (
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/eventlistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	pmodel "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type APIGatewayService struct {
	sender        *sender.Sender
	eventListener *eventlistener.EventListener
}

func (s *APIGatewayService) Login(transport sender.Transport, input *authrpc.LoginInput) (res *authrpc.AuthResponse, err error) {
	result, err := s.sender.PerformRequest(transport, "auth", "login", input, s.eventListener)
	return result.(*authrpc.AuthResponse), err
}

func (s *APIGatewayService) Signup(transport sender.Transport, input *authrpc.SignupInput) (res *authrpc.AuthResponse, err error) {
	result, err := s.sender.PerformRequest(transport, "auth", "signup", input, s.eventListener)
	return result.(*authrpc.AuthResponse), err
}

func (s *APIGatewayService) GetPictures(transport sender.Transport, userId string) (res []*pmodel.Picture, err error) {
	input := prpc.GetUserPicturesRequest{UserId: userId}
	result, err := s.sender.PerformRequest(transport, "pictures", "getPictures", &input, s.eventListener)
	return result.([]*pmodel.Picture), err
}

func (s *APIGatewayService) GetPicture(transport sender.Transport, userId string, pictureId string) (res *pmodel.Picture, err error) {
	input := prpc.GetUserPictureRequest{
		UserId:    userId,
		PictureId: pictureId,
	}
	result, err := s.sender.PerformRequest(transport, "pictures", "getPicture", &input, s.eventListener)
	return result.(*pmodel.Picture), err
}

func (s *APIGatewayService) CreatePicture(transport sender.Transport, input *prpc.CreateUserPictureRequest) (res *pmodel.Picture, err error) {
	result, err := s.sender.PerformRequest(transport, "pictures", "createPicture", &input, s.eventListener)
	return result.(*pmodel.Picture), err
}

func (s *APIGatewayService) DeletePicture(transport sender.Transport, userId string, pictureId string) error {
	input := prpc.DeleteUserPictureRequest{
		UserId:    userId,
		PictureId: pictureId,
	}
	_, err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener)
	return err
}

func (s *APIGatewayService) GetAnalytics(transport sender.Transport, input *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	result, err := s.sender.PerformRequest(transport, "pictures", "deletePicture", &input, s.eventListener)
	return result.(*anrpc.AnalyticsEntryMessage), err
}

func NewAPIGatewayService(s *sender.Sender, el *eventlistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
