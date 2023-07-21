package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
)

type RestReadAdapter struct {
	URL string
	Client *http.Client
}

func (rga RestReadAdapter) Call(ctx context.Context) (string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", rga.URL, nil)
    if err != nil {
        return "", err
    }

    resp, err := rga.Client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status code %v calling: %s", resp.StatusCode, rga.URL)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var restResponse contract.RestResponseBody
    err = json.Unmarshal(body, &restResponse)
    if err != nil {
        return "", err
    }

    return restResponse.Message, nil
}

