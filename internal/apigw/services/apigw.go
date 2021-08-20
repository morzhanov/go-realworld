package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	analyticsdto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	analyticsmodel "github.com/morzhanov/go-realworld/internal/analytics/models"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	"github.com/morzhanov/go-realworld/internal/common/eventlistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	picturemodel "github.com/morzhanov/go-realworld/internal/pictures/models"
	uuid "github.com/satori/go.uuid"
)

type Transport int

const (
	rest Transport = iota
	rpc
	events
)

type APIGatewayService struct {
	sender        *sender.Sender
	eventListener *eventlistener.EventListener
}

func (s *APIGatewayService) handleRest(
	api string,
	endpoint string,
	data interface{},
	headers *sender.Headers,
	res interface{},
) (err error) {
	params := s.sender.API[api].Rest[endpoint]
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	response := s.sender.RestRequest(params.Method, params.Url, b, headers)
	return json.Unmarshal(response, &res)
}

func (s *APIGatewayService) handleEvents(api string, event string, data interface{}, res interface{}) (err error) {
	uuid := uuid.NewV4().String()
	params := s.sender.API[api].Events[event]

	// TODO: move this somewhere to common events package
	type RequestEvent struct {
		Uuid string
		Data string
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	eventData := RequestEvent{uuid, string(jsonData)}
	jsonEventData, err := json.Marshal(&eventData)
	if err != nil {
		return err
	}
	input := sender.EventsRequestInput{
		api,
		params.Event,
		string(jsonEventData),
	}

	// TODO: we need a listener only if response needed
	s.sender.EventsRequest(&input)
	response := make(chan []byte)
	l := eventlistener.Listener{uuid, response}
	err = s.eventListener.AddListener(&l)
	if err != nil {
		return err
	}
	b := <-response
	err = json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	return
}

func (s *APIGatewayService) Login(transport Transport, data *authdto.LoginInput) (res *authdto.LoginDto, err error) {
	switch transport {
	case rest:
		err = s.handleRest("auth", "login", data, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.AuthRpcClient)
		if err != nil {
			return nil, err
		}
		authRpcClient := client.(authrpc.AuthClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := authRpcClient.Login(ctx, &authrpc.LoginInput{
			Username: data.Username,
			Password: data.Password,
		})
		if err != nil {
			return nil, err
		}
		return &authdto.LoginDto{res.AccessToken}, nil
	case events:
		err = s.handleEvents("auth", "login", data, &res)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) Signup(transport Transport, data *authdto.SignupInput) (res *authdto.LoginDto, err error) {
	switch transport {
	case rest:
		err = s.handleRest("auth", "signup", data, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.AuthRpcClient)
		if err != nil {
			return nil, err
		}
		authRpcClient := client.(authrpc.AuthClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := authRpcClient.Signup(ctx, &authrpc.SignupInput{
			Username: data.Username,
			Password: data.Password,
		})
		if err != nil {
			return nil, err
		}
		return &authdto.LoginDto{res.AccessToken}, nil
	case events:
		err = s.handleEvents("auth", "signup", data, &res)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) GetPictures(transport Transport, userId string) (res []*picturemodel.Picture, err error) {
	switch transport {
	case rest:
		// TODO: we need somehow to inject userId into URL params
		err = s.handleRest("pictures", "getPictires", nil, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.AuthRpcClient)
		if err != nil {
			return nil, err
		}
		authRpcClient := client.(picturesrpc.PicturesClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := authRpcClient.GetUserPictures(ctx, &picturesrpc.GetUserPicturesRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
		pics := make([]*picturemodel.Picture, len(res.Pictures))
		for _, pic := range res.Pictures {
			pics = append(pics, &picturemodel.Picture{
				ID:     pic.Id,
				Base64: pic.Base64,
				Title:  pic.Title,
				UserId: pic.UserId,
			})
		}
		return pics, nil
	case events:
		// TODO: we need somehow to inject userId into event
		err = s.handleEvents("pictures", "getPictires", data, &res)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
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

func NewAPIGatewayService(s *sender.Sender, el *eventlistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
