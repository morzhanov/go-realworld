package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"google.golang.org/grpc/metadata"
)

var GrpcMeta opentracing.BuiltinFormat = 4
var EventsHeaders opentracing.BuiltinFormat = 5

func InjectHttpSpan(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func InjectGrpcSpan(span opentracing.Span, ctx context.Context) (*context.Context, error) {
	meta := make(map[string]string, 0)
	err := span.Tracer().Inject(
		span.Context(),
		GrpcMeta,
		opentracing.TextMapCarrier(meta))
	if err != nil {
		return nil, err
	}

	md := make([]string, 0)
	for k, v := range meta {
		md = append(md, k)
		md = append(md, v)
	}

	ct := metadata.AppendToOutgoingContext(ctx, md...)
	return &ct, nil
}

func InjectEventsSpan(span opentracing.Span, m *kafka.Message) error {
	data := make(map[string]string, 0)
	err := span.Tracer().Inject(
		span.Context(),
		EventsHeaders,
		opentracing.TextMapCarrier(data),
	)
	if err != nil {
		return err
	}

	for k, v := range data {
		h := protocol.Header{Key: k, Value: []byte(v)}
		m.Headers = append(m.Headers, h)
	}
	return nil
}

func ExtractHttpSpan(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

func ExtractGrpcSpan(tracer opentracing.Tracer, ctx context.Context) (opentracing.SpanContext, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("error duting grpc span extracting")
	}

	data := make(map[string]string, 0)
	for k, v := range meta {
		data[k] = v[0]
	}

	return tracer.Extract(
		GrpcMeta,
		opentracing.TextMapCarrier(data),
	)
}

func ExtractEventsSpan(tracer opentracing.Tracer, m *kafka.Message) (opentracing.SpanContext, error) {
	data := make(map[string]string, 0)
	for _, header := range m.Headers {
		data[header.Key] = string(header.Value)
	}

	return tracer.Extract(
		EventsHeaders,
		opentracing.TextMapCarrier(data),
	)
}
