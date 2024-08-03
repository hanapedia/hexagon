package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	logger "github.com/hanapedia/hexagon/pkg/operator/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	PromRegistry *prometheus.Registry
)

func init() {
	PromRegistry = prometheus.NewRegistry()
}

func ServeMetrics(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(PromRegistry, promhttp.HandlerOpts{Registry: PromRegistry}))
	srv := &http.Server{Addr: fmt.Sprintf(":%v", config.GetEnvs().METRICS_PORT), Handler: mux}
	wg.Add(1)

	go func() {
		logger.Logger.
			WithField("port", config.GetEnvs().METRICS_PORT).
			Infof("Starting metrics server.")
		if err := srv.ListenAndServe(); err != nil {
			logger.Logger.
				WithField("port", config.GetEnvs().METRICS_PORT).
				WithField("err", err).
				Errorf("Failed to start metrics server.")
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Logger.Infof("Context cancelled. Metrics server shutting down.")
		srv.Shutdown(context.Background())
		wg.Done()
	}()
}
