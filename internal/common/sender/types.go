package sender

import (
	analyticsrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/grpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/grpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/grpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
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

type Sender struct {
	API          *config.ApiConfig
	restClient   *RestClient
	grpcClient   *GrpcClient
	eventsClient *EventsClient
	logger       *zap.Logger
}
