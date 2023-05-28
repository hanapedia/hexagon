package rest

import (
	"net/http"
	"time"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/infrastructure/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type RestClient struct {
	Client *http.Client
}

func NewRestClient() RestClient {
	client := RestClient{
		Client: &http.Client{
			Timeout: time.Duration(time.Millisecond * 50),
		},
	}

	// enable tracing
	if config.GetEnvs().TRACING {
		client.Client.Transport = otelhttp.NewTransport(http.DefaultTransport)
	}

	return client
}

func (restClient RestClient) Close() {
	restClient.Client.CloseIdleConnections()
}
