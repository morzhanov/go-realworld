package services

import (
	analyticsdto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	analyticsmodel "github.com/morzhanov/go-realworld/internal/analytics/models"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	picturemodel "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type APIGatewayService struct {
}

func (s *APIGatewayService) Login(data *authdto.LoginInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
}

func (s *APIGatewayService) Signup(data *authdto.SignupInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
}

func (s *APIGatewayService) GetPictures(userId string) (res []*picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) GetPicture(userId string, pictureId string) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) CreatePicture(userId string, data *picturedto.CreatePicturesDto) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) DeletePicture(userId string, pictureId string) error {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) GetAnalytics(input *analyticsdto.GetLogsInput) (res *analyticsmodel.AnalyticsEntry, err error) {
	// TODO: send request to pictures service
}

func NewAPIGatewayService() *APIGatewayService {
	return &APIGatewayService{}
}
