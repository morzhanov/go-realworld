# generate golang proto files

1. go to proto folder
2. run:

```bash
protoc \
    --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
    --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    analytics.proto
```

## compile all

1. go to `/api/rpc` folder
2. run `bash ./compile_proto.sh`
