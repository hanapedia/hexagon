package utils

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	mathRand "math/rand"
	"time"

	"github.com/hanapedia/the-bench/config/constants"
)

func GenerateRandomString(kbSize int) (string, error) {
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

func GeneratePayloadWithRepositorySize(entrySize constants.RepositoryEntryVariant) (string, error) {
	var payload string
	var err error
	switch entrySize {
	case constants.SMALL:
		payload, err = GenerateRandomString(int(constants.SMALLSIZE))
	case constants.MEDIUM:
		payload, err = GenerateRandomString(int(constants.MEDIUMSIZE))
	case constants.LARGE:
		payload, err = GenerateRandomString(int(constants.LARGESIZE))
	}
	return payload, err
}

// generates random integer from 1 to 100
func GetRandomId(min int, max int) int {
	if min >= max {
		max = 100
		min = 1
	}
	// Seed the random number generator with the current time
	mathRand.Seed(time.Now().UnixNano())
	randomInt := mathRand.Intn(max) + min
	return randomInt
}
