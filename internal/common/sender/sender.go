package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/eventlistener"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

// TODO: move this to commmon package
type Transport int

const (
	RestTransport Transport = iota
	RpcTransport
	EventsTransport
)

type Headers map[string]string

type GrpcClient struct {
	picturesClient  picturesrpc.PicturesClient
	usersClient     usersrpc.UsersClient
	analyticsClient analyticsrpc.AnalyticsClient
	authClient      authrpc.AuthClient
}

type RpcClient int

const (
	UsersRpcClient RpcClient = iota
	PicturesRpcClient
	AnalyticsRpcClient
	AuthRpcClient
)

type RpcRequestInput struct {
	Client RpcClient
	Method string
	Data   []byte
}

type EventsClientItem struct {
	brokers []string
	topic   string
}

type EventsClient struct {
	Auth      *EventsClientItem
	Analytics *EventsClientItem
	Pictures  *EventsClientItem
	Users     *EventsClientItem
}

type BaseEventPayload struct {
	EventId string `json:"eventId"`
}

type EventData struct {
	EventId string
	Data    string
}

type EventsRequestInput struct {
	Service string
	Event   string
	Data    string
}

type Sender struct {
	API          map[string]*ServiceAPI
	restClient   *http.Client
	grpcClient   *GrpcClient
	eventsClient *EventsClient
}

type RestServiceAPIItem struct {
	Method string
	Url    string
}

// TODO: move to common package
type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventsServiceAPIItem struct {
	Event string
}

type ServiceAPI struct {
	Rest   map[string]RestServiceAPIItem
	Events map[string]EventsServiceAPIItem
}

func check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *Sender) PerformRequest(
	transport Transport,
	service string,
	method string,
	input interface{},
	el *eventlistener.EventListener,
) (res interface{}, err error) {
	switch transport {
	case RestTransport:
		params := s.API[service].Rest[method]
		err = s.restRequest(params.Method, params.Url, input, nil, &res)
	case RpcTransport:
		res, err = s.rpcRequest(AuthRpcClient, method, input)
	case EventsTransport:
		uuid := uuid.NewV4().String()
		json, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.eventsRequest(service, method, string(json), uuid, &res, true, el)
	default:
		return nil, fmt.Errorf("Wrong transport type")
	}
	return
}

// TODO: maybe parse functions should be moved to some separate package
func ParseEventsResponse(inputValue string, res interface{}) (payload *EventData, err error) {
	payload = &EventData{}
	err = json.Unmarshal([]byte(inputValue), payload)
	check(err)

	result := reflect.ValueOf(res)
	err = json.Unmarshal([]byte(payload.Data), &result)
	return
}

func ParseRestBody(ctx *gin.Context, input interface{}) (err error) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {

		return err
	}

	in := reflect.ValueOf(input)
	if err = json.Unmarshal(jsonData, &in); err != nil {
		return err
	}
	return nil
}

func (s *Sender) restRequest(
	method string,
	url string,
	data interface{},
	headers *Headers,
	res interface{},
) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	check(err)
	for k, v := range *headers {
		req.Header.Set(k, v)
	}

	response, err := s.restClient.Do(req)
	check(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)

	return json.Unmarshal(body, &res)
}

func (c *Sender) getRpcClient(client RpcClient) (interface{}, error) {
	switch client {
	case UsersRpcClient:
		return &c.grpcClient.usersClient, nil
	case PicturesRpcClient:
		return &c.grpcClient.picturesClient, nil
	case AnalyticsRpcClient:
		return &c.grpcClient.analyticsClient, nil
	case AuthRpcClient:
		return &c.grpcClient.authClient, nil
	default:
		return nil, fmt.Errorf("wrong client")
	}
}

func (s *Sender) rpcRequest(
	client RpcClient,
	method string,
	input interface{},
) (res interface{}, err error) {
	c, err := s.getRpcClient(client)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	fn := reflect.ValueOf(c).Elem().MethodByName(method)

	inputArgs := [2]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(input)}
	returnArgs := fn.Call(inputArgs[:])

	if len(returnArgs) == 1 {
		err = returnArgs[0].Interface().(error)
		return nil, err
	}

	err = returnArgs[1].Interface().(error)
	if err != nil {
		return nil, err
	}
	return returnArgs[0].Interface(), nil
}

func (s *Sender) eventsRequest(
	api string,
	event string,
	data string,
	eventId string,
	res interface{},
	wait bool,
	el *eventlistener.EventListener,
) (err error) {
	params := s.API[api].Events[event]

	eventData := EventData{
		EventId: eventId,
		Data:    data,
	}
	jsonData, err := json.Marshal(&eventData)
	input := EventsRequestInput{
		Service: api,
		Event:   params.Event,
		Data:    string(jsonData),
	}

	var w *kafka.Writer
	switch input.Service {
	case "auth":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  s.eventsClient.Auth.brokers,
			Topic:    s.eventsClient.Auth.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "analytics":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  s.eventsClient.Analytics.brokers,
			Topic:    s.eventsClient.Analytics.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "pictures":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  s.eventsClient.Pictures.brokers,
			Topic:    s.eventsClient.Pictures.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "users":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  s.eventsClient.Users.brokers,
			Topic:    s.eventsClient.Users.topic,
			Balancer: &kafka.LeastBytes{},
		})
	}
	defer w.Close()

	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(input.Event),
			Value: []byte(input.Data),
		},
	)

	if wait {
		response := make(chan []byte)
		l := eventlistener.Listener{Uuid: eventId, Response: response}
		err = el.AddListener(&l)
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

func setupRestClient() *http.Client {
	return &http.Client{}
}

func setupGrpcClient() *GrpcClient {
	// TODO: get addresses from env vars
	picturesAddr := ""
	usersAddr := ""
	analyticsAddr := ""
	authAddr := ""

	conn, err := grpc.Dial(picturesAddr, grpc.WithInsecure(), grpc.WithBlock())
	check(err)
	picturesClient := picturesrpc.NewPicturesClient(conn)

	conn, err = grpc.Dial(usersAddr, grpc.WithInsecure(), grpc.WithBlock())
	check(err)
	usersClient := usersrpc.NewUsersClient(conn)

	conn, err = grpc.Dial(analyticsAddr, grpc.WithInsecure(), grpc.WithBlock())
	check(err)
	analyticsClient := analyticsrpc.NewAnalyticsClient(conn)

	conn, err = grpc.Dial(authAddr, grpc.WithInsecure(), grpc.WithBlock())
	check(err)
	authClient := authrpc.NewAuthClient(conn)

	return &GrpcClient{picturesClient, usersClient, analyticsClient, authClient}
}

func setupEventsClient() *EventsClient {
	// TODO: get values from env vars
	authConnectionUri := ""
	topic := ""
	// ...

	return &EventsClient{
		Auth:      &EventsClientItem{[]string{authConnectionUri}, topic},
		Analytics: &EventsClientItem{[]string{connectionUri}, topic},
		Pictures:  &EventsClientItem{[]string{connectionUri}, topic},
		Users:     &EventsClientItem{[]string{connectionUri}, topic},
	}
}

func NewServiceAPI() map[string]*ServiceAPI {
	// TODO: get services's apis from json files and create ServiceAPI
	// TODO: parse data (json) and create service api struct
	authRestData := make([]byte, 0)
	authEventsData := make([]byte, 0)
	// ...
	return map[string]*ServiceAPI{}
}

func NewSender() *Sender {
	API := NewServiceAPI()
	r := setupRestClient()
	g := setupGrpcClient()
	e := setupEventsClient()

	return &Sender{API, r, g, e}
}
