package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/domain/contract"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type restWriteAdapter struct {
	url         string
	client      *http.Client
	payloadSize int64
	secondary.SecondaryPortBase
}

func (rwa *restWriteAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	payload := utils.GenerateRandomString(rwa.payloadSize)
	restRequestBody := contract.RestRequestBody{
		Payload: &payload,
	}

	logger.Logger.Debugf("Sending request with %v bytes to %s", rwa.payloadSize, rwa.GetDestId())

	jsonRestRequestBody, err := json.Marshal(restRequestBody)
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rwa.url, bytes.NewReader(jsonRestRequestBody))
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rwa.client.Do(req)
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   fmt.Errorf("unexpected status code %v", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	var restResponse contract.RestResponseBody
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: restResponse.Payload,
		Error:   nil,
	}
}
