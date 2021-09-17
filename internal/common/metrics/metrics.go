package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type collector struct {
	opsProcessed prometheus.Counter
}

type Collector interface {
	RegisterMetricsEndpoint(router *gin.Engine)
	RecordBaseMetrics(ctx context.Context)
}

func (mc *collector) RegisterMetricsEndpoint(router *gin.Engine) {
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
}

func (mc *collector) RecordBaseMetrics(ctx context.Context) {
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				mc.opsProcessed.Inc()
				time.Sleep(2 * time.Second)
			}
		}
	}()
}

func NewMetricsCollector(c *config.Config) Collector {
	mc := &collector{}
	mc.opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%v_ops_total", c.ServiceName),
		Help: "The total number of processed events",
	})
	return mc
}
