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

	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func (s *Sender) PerformRequest(
	transport Transport,
	service string,
	method string,
	input interface{},
	el *eventslistener.EventListener,
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

func (s *Sender) SendEventsResponse(eventUuid string, value interface{}) error {
	if !helper.CheckStruct(value) {
		log.Fatal("Value is not struct")
	}

	payload, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	s.eventsRequest("response", "response", string(payload), eventUuid, nil, false, nil)
	return nil
}

func (s *Sender) restRequest(
	method string,
	url string,
	data interface{},
	headers *Headers,
	res interface{},
) (err error) {
	if !helper.CheckStruct(data) {
		log.Fatal("Value is not struct")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	for k, v := range *headers {
		req.Header.Set(k, v)
	}

	response, err := s.restClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

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
	if !helper.CheckStruct(input) {
		log.Fatal("Value is not struct")
	}

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
	el *eventslistener.EventListener,
) (err error) {
	params := s.API[api].Events[event]

	eventData := events.EventData{
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
	case "results":
		w = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  s.eventsClient.Users.brokers,
			Topic:    s.eventsClient.Results.topic,
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
		l := eventslistener.Listener{Uuid: eventId, Response: response}
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

func setupGrpcClient() (*GrpcClient, error) {
	picturesAddr := viper.GetString("PICTURES_GRPC_ADDR")
	usersAddr := viper.GetString("USERS_GRPC_ADDR")
	analyticsAddr := viper.GetString("ANALYTICS_GRPC_ADDR")
	authAddr := viper.GetString("AUTH_GRPC_ADDR")

	conn, err := grpc.Dial(picturesAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	picturesClient := picturesrpc.NewPicturesClient(conn)

	conn, err = grpc.Dial(usersAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	usersClient := usersrpc.NewUsersClient(conn)

	conn, err = grpc.Dial(analyticsAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	analyticsClient := analyticsrpc.NewAnalyticsClient(conn)

	conn, err = grpc.Dial(authAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	authClient := authrpc.NewAuthClient(conn)

	return &GrpcClient{picturesClient, usersClient, analyticsClient, authClient}, nil
}

func setupEventsClient() *EventsClient {
	authConnectionUri := viper.GetString("AUTH_KAFKA_URI")
	authTopic := viper.GetString("AUTH_KAFKA_TOPIC")
	analyticsConnectionUri := viper.GetString("ANALYTICS_KAFKA_URI")
	analyticsTopic := viper.GetString("ANALYTICS_KAFKA_TOPIC")
	picturesConnectionUri := viper.GetString("PICTURES_KAFKA_URI")
	picturesTopic := viper.GetString("PICTURES_KAFKA_TOPIC")
	usersConnectionUri := viper.GetString("USERS_KAFKA_URI")
	usersTopic := viper.GetString("USERS_KAFKA_TOPIC")
	resultsConnectionUri := viper.GetString("RESULTS_KAFKA_URI")
	resultsTopic := viper.GetString("RESULTS_KAFKA_TOPIC")

	return &EventsClient{
		Auth:      &EventsClientItem{[]string{authConnectionUri}, authTopic},
		Analytics: &EventsClientItem{[]string{analyticsConnectionUri}, analyticsTopic},
		Pictures:  &EventsClientItem{[]string{picturesConnectionUri}, picturesTopic},
		Users:     &EventsClientItem{[]string{usersConnectionUri}, usersTopic},
		Results:   &EventsClientItem{[]string{resultsConnectionUri}, resultsTopic},
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

func NewSender() (*Sender, error) {
	API := NewServiceAPI()
	r := setupRestClient()
	e := setupEventsClient()
	g, err := setupGrpcClient()
	if err != nil {
		return nil, err
	}

	return &Sender{API, r, g, e}, err
}
