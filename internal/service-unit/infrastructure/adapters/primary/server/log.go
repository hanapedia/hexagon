package server

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/domain"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func Log(handler *domain.PrimaryAdapterHandler, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	logger.Logger.
		Tracef("Invocation completed | %-30s | %10v ms", handler.GetId(), elapsed)
}
