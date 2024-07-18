package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain/contract"
)

type restReadAdapter struct {
	url string
	client *http.Client
	secondary.SecondaryPortBase
}

func (rra *restReadAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
    req, err := http.NewRequestWithContext(ctx, "GET", rra.url, nil)
    if err != nil {
        return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

    resp, err := rra.client.Do(req)
    if err != nil {
        return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: fmt.Errorf("unexpected status code %v", resp.StatusCode),
		}
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

    var restResponse contract.RestResponseBody
    err = json.Unmarshal(body, &restResponse)
    if err != nil {
        return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
    }

	return secondary.SecondaryPortCallResult{
		Payload: restResponse.Payload,
		Error: nil,
	}
}
