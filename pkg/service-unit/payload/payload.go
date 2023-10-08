package payload

import (
	"github.com/hanapedia/the-bench/pkg/common/utils"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

// GeneratePayload generates payload with given size.
// if the size is not given, default payload size is used.
func GeneratePayload(entrySize constants.PayloadSizeVariant) (string, error) {
	var payload string
	var err error
	switch entrySize {
	case constants.SMALL:
		payload, err = utils.GenerateRandomString(constants.SMALLSIZE)
	case constants.MEDIUM:
		payload, err = utils.GenerateRandomString(constants.MEDIUMSIZE)
	case constants.LARGE:
		payload, err = utils.GenerateRandomString(constants.LARGESIZE)
	default:
		payload, err = utils.GenerateRandomString(constants.DefaultPayloadSize)
	}

	return payload, err
}
