package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/domain/contract"
)

type RestReadAdapter struct {
	URL string
	Client *http.Client
	ports.SecondaryPortBase
}

func (rra *RestReadAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
    req, err := http.NewRequestWithContext(ctx, "GET", rra.URL, nil)
    if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

    resp, err := rra.Client.Do(req)
    if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: fmt.Errorf("unexpected status code %v", resp.StatusCode),
		}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

    var restResponse contract.RestResponseBody
    err = json.Unmarshal(body, &restResponse)
    if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

	return ports.SecondaryPortCallResult{
		Payload: restResponse.Payload,
		Error: nil,
	}
}
