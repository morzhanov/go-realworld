package restcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
)

type BaseRestController struct {
	Router *gin.Engine
	Tracer *opentracing.Tracer
	Logger *zap.Logger
	MC     *metrics.MetricsCollector
}

func (c *BaseRestController) Listen(
	ctx context.Context,
	port string,
) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: c.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("REST Server Failed to start", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutdown REST Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("REST Server shutdown failed", zap.Error(err))
	}
	return nil
}

func (c *BaseRestController) ParseRestBody(ctx *gin.Context, input interface{}) error {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	in := reflect.ValueOf(input)
	return json.Unmarshal(jsonData, &in)
}

func (c *BaseRestController) HandleRestError(ctx *gin.Context, err error) {
	ctx.String(http.StatusInternalServerError, err.Error())
}

func (c *BaseRestController) GetSpan(ctx *gin.Context) *opentracing.Span {
	item, _ := ctx.Get("span")
	span := item.(opentracing.Span)
	return &span
}

func (c *BaseRestController) Handler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := tracing.StartSpanFromHttpRequest(*c.Tracer, ctx.Request)
		ctx.Set("span", span)
		handler(ctx)
		defer span.Finish()
	}
}

func NewRestController(
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) *BaseRestController {
	router := gin.Default()
	c := BaseRestController{router, tracer, logger, mc}
	c.MC.RegisterMetricsEndpoint(router)
	return &c
}
