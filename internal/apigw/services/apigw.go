package services

import (
	"log"

	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	analyticsdto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	analyticsmodel "github.com/morzhanov/go-realworld/internal/analytics/models"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	picturemodel "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type Transport int

const (
	rest Transport = iota
	rpc
	events
)

// TODO: inject here urls for all transports and all services
type APIGatewayService struct {
	sender *sender.Sender
}

func (s *APIGatewayService) Login(transport Transport, data *authdto.LoginInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
	switch transport {
	case rest:
		authLoginParams := s.sender.API["auth"].Rest["login"]
		s.sender.RestRequest(authLoginParams.Method, authLoginParams.Url, data, headers)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.AuthRpcClient)
		check(err)
		authRpcClient := client.(authrpc.AuthClient)
		authRpcClient.ValidateRpcRequest( /*...*/ )
	case events:
		authLoginParams := s.sender.API["auth"].Events["login"]
		s.sender.EventsRequest(sender.EventsRequestInput{"auth", authLoginParams.Event, "..."})
		// TODO: wait for value, add listener
	}
}

func check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *APIGatewayService) Signup(transport Transport, data *authdto.SignupInput) (res *authdto.LoginDto, err error) {
	// TODO: send request to auth service
}

func (s *APIGatewayService) GetPictures(transport Transport, userId string) (res []*picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) GetPicture(transport Transport, userId string, pictureId string) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) CreatePicture(transport Transport, userId string, data *picturedto.CreatePicturesDto) (res *picturemodel.Picture, err error) {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) DeletePicture(transport Transport, userId string, pictureId string) error {
	// TODO: send request to pictures service
}

func (s *APIGatewayService) GetAnalytics(transport Transport, input *analyticsdto.GetLogsInput) (res *analyticsmodel.AnalyticsEntry, err error) {
	// TODO: send request to pictures service
}

func NewAPIGatewayService(s *sender.Sender) *APIGatewayService {
	return &APIGatewayService{s}
}
