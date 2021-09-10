package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

func (s *Sender) PerformRequest(
	transport Transport,
	service string,
	method string,
	input interface{},
	el *eventslistener.EventListener,
	span *opentracing.Span,
) (res interface{}, err error) {
	switch transport {
	case RestTransport:
		apiConfig, err := s.API.GetApiItem(service)
		if err != nil {
			return nil, err
		}
		params := apiConfig.Rest[method]
		url := fmt.Sprintf("http://%s%s", s.restClient.urls[service], params.Url)
		err = s.restRequest(params.Method, url, input, nil, &res, span)
		if err != nil {
			return nil, err
		}
	case RpcTransport:
		res, err = s.rpcRequest(AuthRpcClient, method, input, span)
	case EventsTransport:
		uuidVal := uuid.NewV4().String()
		jsonVal, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		err = s.eventsRequest(service, method, string(jsonVal), uuidVal, &res, true, el, span)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("wrong transport type")
	}
	return
}

func (s *Sender) SendEventsResponse(eventUuid string, value interface{}, span *opentracing.Span) error {
	if !helper.CheckStruct(value) {
		return errors.New("value is not struct")
	}

	payload, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	return s.eventsRequest(
		"response",
		"response",
		string(payload),
		eventUuid,
		nil,
		false,
		nil,
		span,
	)
}

func (s *Sender) restRequest(
	method string,
	url string,
	data interface{},
	headers *http.Header,
	res interface{},
	span *opentracing.Span,
) (err error) {
	if !helper.CheckStruct(data) {
		return errors.New("value is not struct")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v[0])
		}
	}

	err = tracing.InjectHttpSpan(*span, req)
	if err != nil {
		return err
	}

	response, err := s.restClient.http.Do(req)
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
	span *opentracing.Span,
) (res interface{}, err error) {
	if !helper.CheckStruct(input) {
		return nil, err
	}

	c, err := s.getRpcClient(client)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()
	ctx, err = tracing.InjectGrpcSpan(*span, ctx)
	if err != nil {
		return nil, err
	}

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
	span *opentracing.Span,
) (err error) {
	apiConfig, err := s.API.GetApiItem(api)
	if err != nil {
		return err
	}
	params := apiConfig.Events[event]

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

	m := kafka.Message{
		Key:   []byte(input.Event),
		Value: []byte(input.Data),
	}
	tracing.InjectEventsSpan(*span, &m)
	w.WriteMessages(context.Background(), m)

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

func (s *Sender) setupRestClient(c *config.Config) {
	restBaseUrls := RestApiBaseUrls{
		"analytics": fmt.Sprintf("%v:%v", c.RestAddr, c.AnalyticsRestPort),
		"auth":      fmt.Sprintf("%v:%v", c.RestAddr, c.AuthRestPort),
		"pictures":  fmt.Sprintf("%v:%v", c.RestAddr, c.PicturesRestPort),
		"users":     fmt.Sprintf("%v:%v", c.RestAddr, c.UsersRestPort),
		"apigw":     fmt.Sprintf("%v:%v", c.RestAddr, c.ApiGWRestPort),
	}
	s.restClient = &RestClient{
		http: &http.Client{},
		urls: restBaseUrls,
	}
}

func (s *Sender) setupGrpcClient(c *config.Config, cancel context.CancelFunc) {
	conn, err := grpc.Dial(c.PicturesGrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		s.logger.Fatal("error during dialing to pictures grpc server, exiting...")
	}
	picturesClient := picturesrpc.NewPicturesClient(conn)

	conn, err = grpc.Dial(c.UsersGrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		s.logger.Fatal("error during dialing to users grpc server, exiting...")
	}
	usersClient := usersrpc.NewUsersClient(conn)

	conn, err = grpc.Dial(c.AnalyticsGrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		s.logger.Fatal("error during dialing to analytics grpc server, exiting...")
	}
	analyticsClient := analyticsrpc.NewAnalyticsClient(conn)

	conn, err = grpc.Dial(c.AuthGrpcAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		s.logger.Fatal("error during dialing to auth grpc server, exiting...")
	}
	authClient := authrpc.NewAuthClient(conn)

	s.grpcClient = &GrpcClient{
		picturesClient,
		usersClient,
		analyticsClient,
		authClient,
	}
}

func (s *Sender) setupEventsClient(c *config.Config) {
	kafkaUri := c.KafkaUri

	s.eventsClient = &EventsClient{
		Auth:      &EventsClientItem{[]string{kafkaUri}, c.AuthKafkaTopic},
		Analytics: &EventsClientItem{[]string{kafkaUri}, c.AnalyticsKafkaTopic},
		Pictures:  &EventsClientItem{[]string{kafkaUri}, c.PicturesKafkaTopic},
		Users:     &EventsClientItem{[]string{kafkaUri}, c.UsersKafkaTopic},
		Results:   &EventsClientItem{[]string{kafkaUri}, c.ResultsKafkaTopic},
	}
}

func (s *Sender) Connect(c *config.Config, cancel context.CancelFunc) {
	s.setupEventsClient(c)
	s.setupRestClient(c)
	go s.setupGrpcClient(c, cancel)
}

func NewSender(ac *config.ApiConfig, l *zap.Logger) *Sender {
	return &Sender{API: ac, logger: l}
}
