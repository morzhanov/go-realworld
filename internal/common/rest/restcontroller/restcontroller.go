package restcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	errs "github.com/morzhanov/go-realworld/internal/common/errors"
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

type baseRestController struct {
	router *gin.Engine
	tracer opentracing.Tracer
	logger *zap.Logger
	mC     metrics.Collector
}

type BaseRestController interface {
	Listen(ctx context.Context, cancel context.CancelFunc, port string)
	ParseRestBody(ctx *gin.Context, input interface{}) error
	HandleRestError(ctx *gin.Context, err error)
	GetSpan(ctx *gin.Context) *opentracing.Span
	Handler(handler gin.HandlerFunc) gin.HandlerFunc
	Router() *gin.Engine
	Tracer() opentracing.Tracer
	Logger() *zap.Logger
	MC() metrics.Collector
}

func (c *baseRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: c.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancel()
			errs.LogInitializationError(err, "rest controller", c.logger)
			return
		}
	}()

	<-ctx.Done()
	log.Info("Shutdown REST Server ...")

	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		cancel2()
		errs.LogInitializationError(err, "rest controller", c.logger)
	}
}

func (c *baseRestController) ParseRestBody(ctx *gin.Context, input interface{}) error {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	in := reflect.ValueOf(input).Interface()
	return json.Unmarshal(jsonData, &in)
}

func (c *baseRestController) HandleRestError(ctx *gin.Context, err error) {
	c.logger.Error(errors.Unwrap(err).Error())
	if err.Error() == "not authorized" {
		ctx.String(http.StatusUnauthorized, err.Error())
		return
	}
	ctx.String(http.StatusInternalServerError, err.Error())
}

func (c *baseRestController) GetSpan(ctx *gin.Context) *opentracing.Span {
	item, _ := ctx.Get("span")
	span := item.(opentracing.Span)
	return &span
}

func (c *baseRestController) Handler(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := tracing.StartSpanFromHttpRequest(c.tracer, ctx.Request)
		ctx.Set("span", span)
		handler(ctx)
		defer span.Finish()
	}
}

func (c *baseRestController) Router() *gin.Engine        { return c.router }
func (c *baseRestController) Tracer() opentracing.Tracer { return c.tracer }
func (c *baseRestController) Logger() *zap.Logger        { return c.logger }
func (c *baseRestController) MC() metrics.Collector      { return c.mC }

func NewRestController(
	tracer opentracing.Tracer,
	logger *zap.Logger,
	mc metrics.Collector,
) BaseRestController {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders([]string{"authorization"}...)
	router.Use(cors.New(config))
	c := baseRestController{router, tracer, logger, mc}
	c.mC.RegisterMetricsEndpoint(router)
	return &c
}
