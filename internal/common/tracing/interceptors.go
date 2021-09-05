package tracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func InjectHttpSpan(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func InjectGrpcSpan(span opentracing.Span /*, ...*/) error {
	// TODO: ...
}

func InjectEventsSpan(span opentracing.Span /*, ...*/) error {
	// TODO: ...
}

func ExtractHttpSpan(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

func ExtractGrpcSpan(tracer opentracing.Tracer /*, ...*/) (opentracing.SpanContext, error) {
	// TODO: ...
}

func ExtractEventsSpan(tracer opentracing.Tracer /*, ...*/) (opentracing.SpanContext, error) {
	// TODO: ...
}
