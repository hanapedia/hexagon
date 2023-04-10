package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type RestReadAdapter struct {
	URL string
}

type RestReadResponse struct {
	Message string `json:"message"`
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
		return "", errors.New("unexpected status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var restReadResponse RestReadResponse
	err = json.Unmarshal(body, &restReadResponse)
	if err != nil {
		return "", err
	}

	return restReadResponse.Message, nil
}

