package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(kbSize int) (string, error) {
	byteSize := kbSize * 1024
	rawByteSize := byteSize * 3 / 4
	bytes := make([]byte, rawByteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := make([]byte, byteSize)
	base64.StdEncoding.Encode(encoded, bytes)
	return string(encoded[:byteSize]), nil
}
