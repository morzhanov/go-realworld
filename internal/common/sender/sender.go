package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

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

func getGrpcClientType(service string) (RpcClient, error) {
	switch service {
	case "analytics":
		return AnalyticsRpcClient, nil
	case "auth":
		return AuthRpcClient, nil
	case "users":
		return UsersRpcClient, nil
	case "pictures":
		return PicturesRpcClient, nil
	default:
		return -1, fmt.Errorf("wrong service type")
	}
}

func (s *Sender) PerformRequest(
	transport Transport,
	service string,
	method string,
	input interface{},
	el *eventslistener.EventListener,
	span *opentracing.Span,
	meta RequestMeta,
	res interface{},
) error {
	apiConfig, err := s.API.GetApiItem(service)
	if err != nil {
		return err
	}
	switch transport {
	case RestTransport:
		params := apiConfig.Rest[method]
		url := fmt.Sprintf("http://%s%s", s.restClient.urls[service], params.Url)
		err = s.restRequest(params.Method, url, input, nil, &res, span, meta)
		if err != nil {
			return err
		}
	case RpcTransport:
		client, err := getGrpcClientType(service)
		if err != nil {
			return err
		}
		params := apiConfig.Grpc[method]
		if err := s.rpcRequest(client, params.Method, input, span, res); err != nil {
			return err
		}
	case EventsTransport:
		uuidVal := uuid.NewV4().String()
		jsonVal, err := json.Marshal(input)
		if err != nil {
			return err
		}
		err = s.eventsRequest(service, method, string(jsonVal), uuidVal, &res, true, el, span)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("wrong transport type")
	}
	return nil
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
		"results",
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
	meta RequestMeta,
) (err error) {
	if !helper.CheckStruct(data) {
		return errors.New("value is not struct")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if meta != nil && meta["queryparams"] != nil {
		url = fmt.Sprintf("%s?%s", url, meta["queryparams"])
	}
	if meta != nil && meta["urlparams"] != nil {
		urlparams := meta["urlparams"].(UrlParams)
		for k, v := range urlparams {
			param := fmt.Sprintf(":%s", k)
			url = strings.Replace(url, param, v, 1)
		}
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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return errors.New(string(body))
	}
	if reflect.ValueOf(res).IsNil() || len(body) == 0 {
		return nil
	}
	if err := response.Body.Close(); err != nil {
		return err
	}
	return json.Unmarshal(body, &res)
}

func (s *Sender) GetTransportFromContext(ctx context.Context) Transport {
	val := ctx.Value("transport")
	return val.(Transport)
}

func (s *Sender) StringToTransport(transport string) (Transport, error) {
	switch transport {
	case "rest":
		return RestTransport, nil
	case "grpc":
		return RpcTransport, nil
	case "events":
		return EventsTransport, nil
	default:
		return -1, fmt.Errorf("wrong transport %s", transport)
	}
}

func (s *Sender) getRpcClient(client RpcClient) (interface{}, error) {
	switch client {
	case UsersRpcClient:
		if s.grpcClient.usersClient == nil {
			return nil, fmt.Errorf("users grpc server is not connected")
		}
		return &s.grpcClient.usersClient, nil
	case PicturesRpcClient:
		if s.grpcClient.usersClient == nil {
			return nil, fmt.Errorf("pictures grpc server is not connected")
		}
		return &s.grpcClient.picturesClient, nil
	case AnalyticsRpcClient:
		if s.grpcClient.usersClient == nil {
			return nil, fmt.Errorf("analytics grpc server is not connected")
		}
		return &s.grpcClient.analyticsClient, nil
	case AuthRpcClient:
		if s.grpcClient.usersClient == nil {
			return nil, fmt.Errorf("auth grpc server is not connected")
		}
		return &s.grpcClient.authClient, nil
	default:
		return nil, fmt.Errorf("wrong client")
	}
}

func (s *Sender) rpcRequest(
	client RpcClient,
	method string,
	input interface{},
	span *opentracing.Span,
	res interface{},
) error {
	if !helper.CheckStruct(input) {
		return errors.New("value is not a structure")
	}
	c, err := s.getRpcClient(client)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ctx, err = tracing.InjectGrpcSpan(*span, ctx)
	if err != nil {
		return err
	}

	fn := reflect.ValueOf(c).Elem().MethodByName(method)
	inputArgs := [2]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(input)}
	returnArgs := fn.Call(inputArgs[:])
	if len(returnArgs) == 1 {
		err = returnArgs[0].Interface().(error)
		return err
	}
	if err, ok := returnArgs[1].Interface().(error); ok && err != nil {
		return err
	}
	reflect.ValueOf(res).Elem().Set(returnArgs[0])
	return nil
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
	eventKey := "results"
	if api != "results" {
		apiConfig, err := s.API.GetApiItem(api)
		if err != nil {
			return err
		}
		eventKey = apiConfig.Events[event].Event
	}

	eventData := events.EventData{
		EventId: eventId,
		Data:    data,
	}
	jsonData, err := json.Marshal(&eventData)
	input := EventsRequestInput{
		Service: api,
		Event:   eventKey,
		Data:    string(jsonData),
	}

	w := kafka.Writer{
		Balancer: &kafka.LeastBytes{},
	}
	switch input.Service {
	case "analytics":
		w.Addr = kafka.TCP(s.eventsClient.Analytics.brokers[0])
		w.Topic = s.eventsClient.Analytics.topic
	case "auth":
		w.Addr = kafka.TCP(s.eventsClient.Auth.brokers[0])
		w.Topic = s.eventsClient.Auth.topic
	case "pictures":
		w.Addr = kafka.TCP(s.eventsClient.Pictures.brokers[0])
		w.Topic = s.eventsClient.Pictures.topic
	case "users":
		w.Addr = kafka.TCP(s.eventsClient.Users.brokers[0])
		w.Topic = s.eventsClient.Users.topic
	case "results":
		w.Addr = kafka.TCP(s.eventsClient.Results.brokers[0])
		w.Topic = s.eventsClient.Results.topic
	}

	m := kafka.Message{
		Key:   []byte(input.Event),
		Value: []byte(input.Data),
	}
	response := make(chan []byte)
	l := eventslistener.Listener{Uuid: eventId, Response: response}
	if wait {
		if err := el.AddListener(&l); err != nil {
			return err
		}
	}
	if err := tracing.InjectEventsSpan(*span, &m); err != nil {
		return err
	}
	if err := w.WriteMessages(context.Background(), m); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	if wait {
		b := <-response
		if err := el.RemoveListener(&l); err != nil {
			return err
		}
		if err := json.Unmarshal(b, &res); err != nil {
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
	s.grpcClient = &GrpcClient{}
	go func() {
		uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.AnalyticsGrpcPort)
		conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			cancel()
			s.logger.Fatal("error during dialing to analytics grpc server, exiting...")
		}
		s.grpcClient.analyticsClient = analyticsrpc.NewAnalyticsClient(conn)
	}()
	go func() {
		uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.AuthGrpcPort)
		conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			cancel()
			s.logger.Fatal("error during dialing to auth grpc server, exiting...")
		}
		s.grpcClient.authClient = authrpc.NewAuthClient(conn)
	}()
	go func() {
		uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.PicturesGrpcPort)
		conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			cancel()
			s.logger.Fatal("error during dialing to pictures grpc server, exiting...")
		}
		s.grpcClient.picturesClient = picturesrpc.NewPicturesClient(conn)
	}()
	go func() {
		uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.UsersGrpcPort)
		conn, err := grpc.Dial(uri, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			cancel()
			s.logger.Fatal("error during dialing to users grpc server, exiting...")
		}
		s.grpcClient.usersClient = usersrpc.NewUsersClient(conn)
	}()
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
	s.setupGrpcClient(c, cancel)
}

func NewSender(ac *config.ApiConfig, l *zap.Logger) *Sender {
	return &Sender{API: ac, logger: l}
}
