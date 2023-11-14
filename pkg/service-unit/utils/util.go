package util

import (
	crand "crypto/rand"
	"encoding/base64"
	mrand "math/rand"
)

func GenerateRandomString(byteSize int64) (string, error) {
	rawByteSize := byteSize * 3 / 4
	bytes := make([]byte, rawByteSize)
	_, err := crand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := make([]byte, byteSize)
	base64.StdEncoding.Encode(encoded, bytes)
	return string(encoded[:byteSize]), nil
}

// generates random integer from 1 to 100
func GetRandomId(min int, max int) int {
	if min >= max {
		max = 100
		min = 1
	}

	randomInt := mrand.Intn(max) + min
	return randomInt
}
