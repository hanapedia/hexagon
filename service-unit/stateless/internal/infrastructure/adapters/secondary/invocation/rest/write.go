package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/application/ports"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
)

type RestWriteAdapter struct {
	URL string
	Client *http.Client
	ports.SecondaryPortBase
}

func (rwa *RestWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GenerateRandomString(constants.PayloadSize)
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}

	restRequestBody := contract.RestRequestBody{
		Message: fmt.Sprintf("Posting %vkB of random text to %s", constants.PayloadSize, rwa.URL),
		Payload: &payload,
	}
	jsonRestRequestBody, err := json.Marshal(restRequestBody)
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rwa.URL, bytes.NewReader(jsonRestRequestBody))
	if err != nil {
        return ports.SecondaryPortCallResult{
			Payload: nil,
			Error: err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rwa.Client.Do(req)
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

