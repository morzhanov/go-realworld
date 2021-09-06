package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/segmentio/kafka-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func StartSpanFromHttpRequest(tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, _ := ExtractHttpSpan(tracer, r)
	return tracer.StartSpan("http-receive", ext.RPCServerOption(spanCtx))
}

func StartSpanFromGrpcRequest(tracer opentracing.Tracer, ctx context.Context) opentracing.Span {
	spanCtx, _ := ExtractGrpcSpan(tracer, ctx)
	return tracer.StartSpan("grpc-receive", ext.RPCServerOption(spanCtx))
}

func StartSpanFromEventsRequest(tracer opentracing.Tracer, m *kafka.Message) opentracing.Span {
	spanCtx, _ := ExtractEventsSpan(tracer, m)
	return tracer.StartSpan("event-receive", ext.RPCServerOption(spanCtx))
}

func NewTracer(ctx context.Context, c *config.Config, logger *zap.Logger) (*opentracing.Tracer, error) {
	cfg := jconfig.Configuration{
		ServiceName: c.ServiceName,
	}

	tracer, closer, err := cfg.NewTracer(jconfig.Logger(NewJeagerLogger(logger)))
	if err != nil {
		return nil, fmt.Errorf("cannot init Jaeger tracer: %v", err)
	}

	go func() {
		<-ctx.Done()
		closer.Close()
	}()

	return &tracer, nil
}
