package rest_interface

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type RestWriteInterface struct {
	URL string
}

type RestWriteResponse struct {
	Message string `json:"message"`
}

func (rpi RestWriteInterface) Run() (string, error) {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	resp, err := client.Post(rpi.URL, "application/json", nil)
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

	var restWriteResponse RestWriteResponse
	err = json.Unmarshal(body, &restWriteResponse)
	if err != nil {
		return "", err
	}

	return restWriteResponse.Message, nil
}

