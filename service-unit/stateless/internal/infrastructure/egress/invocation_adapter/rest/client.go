package rest

import (
	"net/http"
	"time"
)

type RestClient struct {
	Client *http.Client
}

func NewRestClient() RestClient {
	return RestClient{
		Client: &http.Client{
			Timeout: time.Duration(time.Millisecond * 50),
		},
	}
}

func (restClient RestClient) Close() {
	restClient.Client.CloseIdleConnections()
}
