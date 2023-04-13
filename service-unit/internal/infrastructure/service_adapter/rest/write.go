package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type RestWriteAdapter struct {
	URL string
}

type RestWriteResponse struct {
	Message string `json:"message"`
}

func (rwa RestWriteAdapter) Call() (string, error) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	resp, err := client.Post(rwa.URL, "application/json", nil)
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

	var restWriteResponse RestWriteResponse
	err = json.Unmarshal(body, &restWriteResponse)
	if err != nil {
		return "", err
	}

	return restWriteResponse.Message, nil
}

