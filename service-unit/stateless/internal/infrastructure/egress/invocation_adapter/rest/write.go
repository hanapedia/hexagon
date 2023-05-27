package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
	"github.com/hanapedia/the-bench/service-unit/stateless/pkg/utils"
)

type RestWriteAdapter struct {
	URL string
	Client *http.Client
}

func (rwa RestWriteAdapter) Call() (string, error) {
	payload, err := utils.GenerateRandomString(constants.PayloadSize)
	if err != nil {
		return "", err
	}

	restRequestBody := contract.RestRequestBody{
		Message: fmt.Sprintf("Posting %vkB of random text to %s", constants.PayloadSize, rwa.URL),
		Payload: &payload,
	}
	jsonRestRequestBody, err := json.Marshal(restRequestBody)
	if err != nil {
		return "", err
	}

	resp, err := rwa.Client.Post(rwa.URL, "application/json", bytes.NewReader(jsonRestRequestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status code" + rwa.URL)
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
