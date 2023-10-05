package utils

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	mathRand "math/rand"

	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

func GenerateRandomString(kbSize constants.PayloadSize) (string, error) {
	byteSize := kbSize * 1024
	rawByteSize := byteSize * 3 / 4
	bytes := make([]byte, rawByteSize)
	_, err := cryptoRand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := make([]byte, byteSize)
	base64.StdEncoding.Encode(encoded, bytes)
	return string(encoded[:byteSize]), nil
}

// GeneratePayload generates payload with given size.
// if the size is not given, default payload size is used.
func GeneratePayload(entrySize constants.PayloadSizeVariant) (string, error) {
	var payload string
	var err error
	switch entrySize {
	case constants.SMALL:
		payload, err = GenerateRandomString(constants.SMALLSIZE)
	case constants.MEDIUM:
		payload, err = GenerateRandomString(constants.MEDIUMSIZE)
	case constants.LARGE:
		payload, err = GenerateRandomString(constants.LARGESIZE)
	default:
		payload, err = GenerateRandomString(constants.DefaultPayloadSize)
	}

	return payload, err
}

// generates random integer from 1 to 100
func GetRandomId(min int, max int) int {
	if min >= max {
		max = 100
		min = 1
	}

	randomInt := mathRand.Intn(max) + min
	return randomInt
}
