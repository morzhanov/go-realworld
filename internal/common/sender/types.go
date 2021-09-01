package sender

import (
	"net/http"

	analyticsrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	authrpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	picturesrpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	usersrpc "github.com/morzhanov/go-realworld/api/rpc/users"
)

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
	Results   *EventsClientItem
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
