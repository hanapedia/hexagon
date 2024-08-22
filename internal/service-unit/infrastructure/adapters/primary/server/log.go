package server

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	logger "github.com/hanapedia/hexagon/pkg/operator/log"
)

func Log(handler *domain.PrimaryAdapterHandler, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	logger.Logger.
		Infof("Invocation completed | %-30s | %10v ms", handler.GetId(), elapsed)
}
