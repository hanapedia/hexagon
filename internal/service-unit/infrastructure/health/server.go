package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/config"
	logger "github.com/hanapedia/hexagon/pkg/operator/log"
)

type ReadyResponse struct {
	Ready bool `json:"ready"`
}

func ServeHealth(ctx context.Context, wg *sync.WaitGroup) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ReadyResponse{Ready: true})
	})
	srv := &http.Server{Addr: fmt.Sprintf(":%v", config.GetEnvs().HEALTH_PORT), Handler: mux}
	wg.Add(1)

	go func() {
		logger.Logger.
			WithField("port", config.GetEnvs().HEALTH_PORT).
			Infof("Starting health server.")
		if err := srv.ListenAndServe(); err != nil {
			logger.Logger.
				WithField("port", config.GetEnvs().HEALTH_PORT).
				WithField("err", err).
				Errorf("Failed to start health server.")
		}
	}()

	go func() {
		<-ctx.Done()
		logger.Logger.Infof("Context cancelled. Health server shutting down.")
		srv.Shutdown(context.Background())
		wg.Done()
	}()
}
