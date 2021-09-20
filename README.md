# Go Real World example

## Idea

The main purpose of the project is to gain Golang development skills and knowledge and apply them in practice.

The project idea is quite simple: a basinc image uploading service. Fom simplicity images should be uploaded in the image/base64 format.

Project contains custom API gateway service in order to transform REST calls to the desired transport calls. Available transports are:

- REST
- gRPC
- Events (via Kafka broker)

Project contains Users and Auth services in order to register and authenticate users.

Also project contains Analytics service which writes API call logs to the Kafka topic.

## Project architecture

<img src="https://i.ibb.co/McV4YQ6/1632127632115.png" alt="project architecture"/>

### API Gateway

Api gateway is the main and most complex part of the system. API gateway receives REST requests from the Web App and proxies them to the internal services with requested transport.

### Web App

Web App is a react application that provides UI tools to:

- Register new users
- Authenticate existing users and provide access to the internal APIs
- Perform picture-related requests
  - Upload a picture
  - Get all pictures
  - Get a single picture
  - Delete a picture
- Retrieve analytics data

### Internal services

- Analytics service - responsible for loogging and providing request call logs information
- Users service - responsible for user-related CRUD operations
- Auth service - responsible for user's registration and authentication flow
- Pictures service - responsible for picture-related CRUD operations

## Project Structure

- /api - API-related packages and files
  - /api/grpc - .proto definition files and compiled .go files for grpc. Also contains `compile_proto.sh` script to compile .proto files into .go files.
  - /api/postman - contains postman collection for local API calls through REST and GRPC transports
- /build - for further development. Should contain CI/CD configuration files
- /cmd - contains `main.go` files for each service. This folder is the starting point for each project service.
- /configs - contains `*.env` files which contain environment variables declaration for each service. Also contains `api.yml` files with API schema declaration for each transport.
- /deploy - contains docker-compose and Dockerfile files.
  - docker-compose.deps.yml - contains docker-compose config with dependant images (ELK, Prometheus, Jeager, etc.)
  - docker-compose.yml - main docker-compose config files
- /internal - contains internal application packages
  - /internal/common - contains common application packages
  - /internal/common/config - configuration service
  - /internal/common/db - database initialization and management service 
  - /internal/common/errors - error creation/parsing helper functions
  - /internal/common/events - contains events-related helper functions and basic EventsController interface
  - /internal/common/grpc - contains grpc-related helper functions and basic GrpcService interface
  - /internal/common/logger - contains main logger configuration
  - /internal/common/metrics - Prometheus metrics collector service
  - /internal/common/mq - Kafka message queue configuration service
  - /internal/common/rest - contains rest-related helper functions and basic RestController interface
  - /internal/common/sender - contains main interface and structure to manage cross-service communication
  - /internal/common/tracing - Jaeger tracer configuration files and main interface

## Used Libraries and Technologies

### Libs

- <a href="github.com/dgrijalva/jwt-go">jwt-go</a> - jwt parsing and verification
- <a href="github.com/gin-gonic/gin">gin</a> - main REST api transport library 
- <a href="github.com/golang-migrate/migrate/v4">migrate</a> - database migrations library
- <a href="github.com/jmoiron/sqlx">sqlx</a> - enhanced golang SQL client
- <a href="github.com/lib/pq">pq</a> - golang postgresql driver
- <a href="github.com/opentracing/opentracing-go">opentracing-go</a> - golang opentracing library
- <a href="github.com/prometheus/client_golang">prometheus</a> - golang Prometheus library
- <a href="github.com/segmentio/kafka-go">kafka-go</a> - golang Kafka library
- <a href="github.com/spf13/viper">viper</a> - config files parsing library
- <a href="github.com/uber/jaeger-client-go">jaeger-client-go</a> - golang Jaeger library
- <a href="go.uber.org/zap">zap</a> - logging library
- <a href="google.golang.org/grpc">grpc</a> - gRPC library
- <a href="google.golang.org/protobuf">protobuf</a> - protobuf library
- <a href="github.com/spf13/viper">viper</a> - config files parsing library
- <a href="github.com/spf13/viper">viper</a> - config files parsing library
- <a href="github.com/spf13/viper">viper</a> - config files parsing library

### Technologies

- Go - main project programming language
- Typescript - web app main programming language
- PostgreSQL - main Database management system
- Jaeger - system tracing component
- Elasticsearch - ELK stack database
- Kibana - ELK stack UI application
- Filebeat - ELK stack logs consumer component
- Prometheus - metrics collector component
- Grafana - metrics UI component
- Zookeeper - Kafka brokers base communication configuration component
- Kafka - main component for event-driven communication between services
- Docker & Docker compose - main containerization components

## Internal Common package

`common` package provides common packages that used in all application services.

- config package provides application configuration component. The package uses `viper` for configuration files parsing.
- db package uses `sqlx` and `pq` libraries and provides PostgreSQL client component.
- errors package provides error handling features.
- events package provides base Events controller component and additional events-related functions.
- grpc package provides base Grpc service component.
- logger service uses `zap` as a base logger and provides configured application logger.
- metrics service uses `prometheus` and provides metrics collector component.
- mq package configures `kafka` connections and provides functions to create Kafka Reader and Writer components.
- rest package provides base REST controller and additional rest-related functions (like req/res body parsing).
- sender package is the core communication component. The package provides `perforRequest` method that allows to perform API calls to the services via desired transport.
- tracing package uses `jaeger` and provides main tracing component.

## Internal Service packages

Each service handles requests from 3 types of transports (gRPC, REST, Events). Each service has it's own database connection for data persistence.

Packages:

- events - contains event controller
- grpc - contains gRPC controller
- rest - contains REST controller
- services - contains core domain business logic
- models - contains database models
- migrations - SQL migrations for `migrate` package

## Local Running

You can run all services with docker-compose:

```bash
cd deploy
docker-compose up -d
```

Also, you can run each service separately:
```bash
go run ./cmd/<service-name>/main.go
```

In order to run services separatelly you should deploy application dependencies:
```
cd deploy
docker-compose up -f ./docker-compose.deps.yml -d
```

### API calls

`/api/postman` directory contains Postman collection. This collection should be imported into local postman application.

<img src="https://i.ibb.co/tpv20fy/1632129593899.png" alt="postman"/>

#### REST
REST calls could be performed via postman collection.

#### gRPC
gRPC calls could be performed via postman collection. In order to perform these request you should intall <a href="https://github.com/jnewmano/grpc-json-proxy">grpc-json-proxy</a> locally.

#### Events
Events calls could be performed via <a href="https://github.com/edenhill/kcat">kafkacat</a>. 

Example:
```
kcat -C -b localhost:29092 -t auth -P -K:
login:{"username":"username", "password":"pwd"}
```

## Further development ideas

- add caching layer with Redis and consistent hashing
- add OpenAPI for REST services
- add Vault to store secrets
- add tests
- add benchmarks
- add go perf metrics
- add Kubernetes (with Helm or Kustomize)
- add CI/CD configuration