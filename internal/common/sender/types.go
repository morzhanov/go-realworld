package sender

import (
	"context"
	analyticsrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/grpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/grpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/grpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
)

type Transport int

const (
	RestTransport Transport = iota
	RpcTransport
	EventsTransport
)

type RequestMeta map[string]interface{}

type UrlParams map[string]string

type RestApiBaseUrls map[string]string

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
	Results   *EventsClientItem
}

type RestClient struct {
	http *http.Client
	urls RestApiBaseUrls
}

type EventsRequestInput struct {
	Service string
	Event   string
	Data    string
}

type sender struct {
	API          *config.ApiConfig
	restClient   *RestClient
	grpcClient   *GrpcClient
	eventsClient *EventsClient
	logger       *zap.Logger
}

type Sender interface {
	Connect(c *config.Config, cancel context.CancelFunc)
	PerformRequest(transport Transport, service string, method string, input interface{}, el eventslistener.EventListener, span *opentracing.Span, meta RequestMeta, res interface{}) error
	SendEventsResponse(eventUuid string, value interface{}, span *opentracing.Span) error
	GetTransportFromContext(ctx context.Context) Transport
	StringToTransport(transport string) (Transport, error)
	TransportToString(transport Transport) (string, error)
	GetAPI() *config.ApiConfig
	CheckStruct(val interface{}) bool
}
