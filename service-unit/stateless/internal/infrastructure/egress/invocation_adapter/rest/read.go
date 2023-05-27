package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/contract"
)

type RestReadAdapter struct {
	URL string
	Client *http.Client
}

func (rga RestReadAdapter) Call() (string, error) {
	resp, err := rga.Client.Get(rga.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status code calling:" + rga.URL)
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
