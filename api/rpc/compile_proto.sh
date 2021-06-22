# !bash

protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      ./analytics/analytics.proto

protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      ./auth/auth.proto

protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      ./pictures/pictures.proto

protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      ./users/users.proto
