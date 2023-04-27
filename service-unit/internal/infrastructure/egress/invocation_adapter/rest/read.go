package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/contract"
)

type RestReadAdapter struct {
	URL string
}

func (rga RestReadAdapter) Call() (string, error) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	resp, err := client.Get(rga.URL)
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
