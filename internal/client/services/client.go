package services

import (
	analyticsdto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	analyticsmodel "github.com/morzhanov/go-realworld/internal/analytics/models"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	picturemodel "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type ClientService struct {
}

func (s *ClientService) Login(data *authdto.LoginInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
}

func (s *ClientService) Signup(data *authdto.SignupInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
}

func (s *ClientService) GetPictures(userId string) (res []*picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *ClientService) GetPicture(userId string, pictureId string) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *ClientService) CreatePicture(userId string, data *picturedto.CreatePicturesDto) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *ClientService) DeletePicture(userId string, pictureId string) error {
	// TODO: send request to pictures service
}

func (s *ClientService) GetAnalytics(input *analyticsdto.GetLogsInput) (res *analyticsmodel.AnalyticsEntry, err error) {
	// TODO: send request to pictures service
}

func NewClientService() *ClientService {
	return &ClientService{}
}
