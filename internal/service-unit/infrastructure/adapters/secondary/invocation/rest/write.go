package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/domain/contract"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
	"github.com/hanapedia/the-bench/pkg/service-unit/utils"
)

type restWriteAdapter struct {
	url     string
	client  *http.Client
	payload constants.PayloadSizeVariant
	ports.SecondaryPortBase
}

func (rwa *restWriteAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	payload, err := utils.GeneratePayload(rwa.payload)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	restRequestBody := contract.RestRequestBody{
		Message: fmt.Sprintf("Posting %s payload of random text to %s", rwa.payload, rwa.GetDestId()),
		Payload: &payload,
	}
	jsonRestRequestBody, err := json.Marshal(restRequestBody)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rwa.url, bytes.NewReader(jsonRestRequestBody))
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rwa.client.Do(req)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   fmt.Errorf("unexpected status code %v", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	var restResponse contract.RestResponseBody
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return ports.SecondaryPortCallResult{
		Payload: restResponse.Payload,
		Error:   nil,
	}
}
