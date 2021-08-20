package sender

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
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

func (c *Sender) RestRequest(method string, url string, data []byte, headers *Headers) []byte {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	check(err)
	for k, v := range *headers {
		req.Header.Set(k, v)
	}

	res, err := c.restClient.Do(req)
	check(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	check(err)
	return body
}

func (c *Sender) GetRpcClient(client RpcClient) (interface{}, error) {
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

func (c *Sender) EventsRequest(input *EventsRequestInput) {
	var w *kafka.Writer
	switch input.Service {
	case "auth":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  c.eventsClient.Auth.brokers,
			Topic:    c.eventsClient.Auth.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "analytics":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  c.eventsClient.Analytics.brokers,
			Topic:    c.eventsClient.Analytics.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "pictures":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  c.eventsClient.Pictures.brokers,
			Topic:    c.eventsClient.Pictures.topic,
			Balancer: &kafka.LeastBytes{},
		})
	case "users":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  c.eventsClient.Users.brokers,
			Topic:    c.eventsClient.Users.topic,
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
