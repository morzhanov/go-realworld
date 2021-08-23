package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	andto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	anevents "github.com/morzhanov/go-realworld/internal/analytics/events"
	anmodel "github.com/morzhanov/go-realworld/internal/analytics/models"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	authevents "github.com/morzhanov/go-realworld/internal/auth/events"
	"github.com/morzhanov/go-realworld/internal/common/eventlistener"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	pdto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	pevents "github.com/morzhanov/go-realworld/internal/pictures/events"
	pmodel "github.com/morzhanov/go-realworld/internal/pictures/models"
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
	method string,
	url string,
	data interface{},
	headers *sender.Headers,
	res interface{},
) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	response := s.sender.RestRequest(method, url, b, headers)
	return json.Unmarshal(response, &res)
}

func (s *APIGatewayService) handleEvents(
	api string,
	event string,
	data string,
	eventId string,
	res interface{},
	wait bool,
) (err error) {
	params := s.sender.API[api].Events[event]
	input := sender.EventsRequestInput{
		Service: api,
		Event:   params.Event,
		Data:    data,
	}

	if wait {
		s.sender.EventsRequest(&input)
		response := make(chan []byte)
		l := eventlistener.Listener{Uuid: eventId, Response: response}
		err = s.eventListener.AddListener(&l)
		if err != nil {
			return err
		}
		b := <-response
		err = json.Unmarshal(b, &res)
		if err != nil {
			return err
		}
	}
	return
}

func (s *APIGatewayService) Login(transport Transport, data *authdto.LoginInput) (res *authdto.LoginDto, err error) {
	switch transport {
	case rest:
		params := s.sender.API["auth"].Rest["login"]
		err = s.handleRest(params.Method, params.Url, data, nil, &res)
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
		return &authdto.LoginDto{AccessToken: res.AccessToken}, nil
	case events:
		uuid := uuid.NewV4().String()
		input := authevents.LoginInput{
			BaseEventPayload: authevents.BaseEventPayload{EventId: uuid},
			LoginInput: authdto.LoginInput{
				Username: data.Username,
				Password: data.Password,
			},
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("auth", "login", string(json), uuid, &res, true)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) Signup(transport Transport, data *authdto.SignupInput) (res *authdto.LoginDto, err error) {
	switch transport {
	case rest:
		params := s.sender.API["auth"].Rest["signup"]
		err = s.handleRest(params.Method, params.Url, data, nil, &res)
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
		return &authdto.LoginDto{AccessToken: res.AccessToken}, nil
	case events:
		uuid := uuid.NewV4().String()
		input := authevents.Signup{
			BaseEventPayload: authevents.BaseEventPayload{EventId: uuid},
			SignupInput: authdto.SignupInput{
				Username: data.Username,
				Password: data.Password,
			},
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("pictures", "signup", string(json), uuid, &res, true)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) GetPictures(transport Transport, userId string) (res []*pmodel.Picture, err error) {
	switch transport {
	case rest:
		params := s.sender.API["pictures"].Rest["getPictures"]
		url := strings.Replace(params.Url, ":userId", userId, 0)
		err = s.handleRest(params.Method, url, nil, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.PicturesRpcClient)
		if err != nil {
			return nil, err
		}
		picRpcClient := client.(prpc.PicturesClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := picRpcClient.GetUserPictures(ctx, &prpc.GetUserPicturesRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
		pics := make([]*pmodel.Picture, len(res.Pictures))
		for _, pic := range res.Pictures {
			pics = append(pics, &pmodel.Picture{
				ID:     pic.Id,
				Base64: pic.Base64,
				Title:  pic.Title,
				UserId: pic.UserId,
			})
		}
		return pics, nil
	case events:
		uuid := uuid.NewV4().String()
		input := pevents.GetUserPicturesInput{
			BaseEventPayload: pevents.BaseEventPayload{EventId: uuid},
			UserId:           userId,
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("pictures", "pictures:get", string(json), uuid, &res, true)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) GetPicture(transport Transport, userId string, pictureId string) (res *pmodel.Picture, err error) {
	switch transport {
	case rest:
		params := s.sender.API["pictures"].Rest["getPicture"]
		url := strings.Replace(params.Url, ":userId", userId, 0)
		url = strings.Replace(url, ":pictureId", pictureId, 0)
		err = s.handleRest(params.Method, url, nil, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.PicturesRpcClient)
		if err != nil {
			return nil, err
		}
		picRpcClient := client.(prpc.PicturesClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := picRpcClient.GetUserPicture(ctx, &prpc.GetUserPictureRequest{
			UserId:    userId,
			PictureId: pictureId,
		})
		if err != nil {
			return nil, err
		}
		return &pmodel.Picture{
			ID:     res.Id,
			Base64: res.Base64,
			Title:  res.Title,
			UserId: res.UserId,
		}, nil
	case events:
		uuid := uuid.NewV4().String()
		input := pevents.GetUserPictureInput{
			BaseEventPayload: pevents.BaseEventPayload{EventId: uuid},
			UserId:           userId,
			PictireId:        pictureId,
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("pictures", "pictures:get_one", string(json), uuid, &res, true)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) CreatePicture(transport Transport, userId string, data *pdto.CreatePicturesDto) (res *pmodel.Picture, err error) {
	switch transport {
	case rest:
		params := s.sender.API["pictures"].Rest["createPicture"]
		url := strings.Replace(params.Url, ":userId", userId, 0)
		err = s.handleRest(params.Method, url, nil, nil, &res)
	case rpc:
		client, err := s.sender.GetRpcClient(sender.PicturesRpcClient)
		if err != nil {
			return nil, err
		}
		picRpcClient := client.(prpc.PicturesClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		res, err := picRpcClient.CreateUserPicture(ctx, &prpc.CreateUserPictureRequest{
			UserId: userId,
			Title:  data.Title,
			Base64: data.Base64,
		})
		if err != nil {
			return nil, err
		}
		return &pmodel.Picture{
			ID:     res.Id,
			Base64: res.Base64,
			Title:  res.Title,
			UserId: res.UserId,
		}, nil
	case events:
		uuid := uuid.NewV4().String()
		input := pevents.CreateUserPictureInput{
			BaseEventPayload: pevents.BaseEventPayload{EventId: uuid},
			UserId:           userId,
			CreatePicturesDto: pdto.CreatePicturesDto{
				Title:  data.Title,
				Base64: data.Base64,
			},
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("pictures", "pictures:create", string(json), uuid, &res, true)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

func (s *APIGatewayService) DeletePicture(transport Transport, userId string, pictureId string) error {
	switch transport {
	case rest:
		params := s.sender.API["pictures"].Rest["deletePicture"]
		url := strings.Replace(params.Url, ":userId", userId, 0)
		url = strings.Replace(url, ":pictireId", pictureId, 0)
		err := s.handleRest(params.Method, url, nil, nil, nil)
		return err
	case rpc:
		client, err := s.sender.GetRpcClient(sender.PicturesRpcClient)
		if err != nil {
			return err
		}
		picRpcClient := client.(prpc.PicturesClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		_, err = picRpcClient.DeleteUserPicture(ctx, &prpc.DeleteUserPictureRequest{
			UserId:    userId,
			PictureId: pictureId,
		})
		return err
	case events:
		uuid := uuid.NewV4().String()
		input := pevents.DeleteUserPictureInput{
			BaseEventPayload: pevents.BaseEventPayload{EventId: uuid},
			UserId:           userId,
			PictireId:        pictureId,
		}
		json, err := json.Marshal(input)
		if err != nil {
			return err
		}
		err = s.handleEvents("pictures", "pictures:create", string(json), uuid, nil, true)
		return err
	default:
		return fmt.Errorf("Wrong transport type")
	}
}

func (s *APIGatewayService) GetAnalytics(transport Transport, input *andto.GetLogsInput) (res *anmodel.AnalyticsEntry, err error) {
	switch transport {
	case rest:
		params := s.sender.API["analytics"].Rest["get"]
		var url string
		if input.Offset != 0 {
			url = strings.Replace(params.Url, ":offset", strconv.Itoa(input.Offset), 0)
		} else {
			url = params.Url
		}
		err := s.handleRest(params.Method, url, nil, nil, &res)
		return res, err
	case rpc:
		client, err := s.sender.GetRpcClient(sender.AnalyticsRpcClient)
		if err != nil {
			return nil, err
		}
		picRpcClient := client.(anrpc.AnalyticsClient)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
		result, err := picRpcClient.GetLog(ctx, &anrpc.GetLogRequest{
			Offset: int32(input.Offset),
		})
		return &anmodel.AnalyticsEntry{
			ID:        result.Id,
			UserID:    result.UserId,
			Operation: result.Operation,
			Data:      result.Data,
		}, err
	case events:
		uuid := uuid.NewV4().String()
		input := anevents.GetLogsInput{
			BaseEventPayload: anevents.BaseEventPayload{EventId: uuid},
			GetLogsInput: andto.GetLogsInput{
				Offset: input.Offset,
			},
		}
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.handleEvents("pictures", "pictures:create", string(json), uuid, &res, true)
		return res, err
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
}

func NewAPIGatewayService(s *sender.Sender, el *eventlistener.EventListener) *APIGatewayService {
	return &APIGatewayService{s, el}
}
