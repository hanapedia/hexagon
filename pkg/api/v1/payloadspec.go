package v1

import (
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"k8s.io/apimachinery/pkg/api/resource"
)

// PayloadSpec contains configuration for payload sent by adapters
type PayloadSpec struct {
	Size    string                       `json:"size,omitempty"`
	Variant constants.PayloadSizeVariant `json:"variant,omitempty" validate:"omitempty,oneof=small medium large"`
	Count   int                          `json:"payloadCount,omitempty"`
}

func ParseResource(payloadSpec PayloadSpec) (int64, bool) {
	if payloadSpec.Size != "" {
		quantity, err := resource.ParseQuantity(payloadSpec.Size)
		if err == nil {
			return quantity.AsInt64()
		}
	}
	return 0, false
}

// GetPayloadSize parses resource size using
// "k8s.io/apimachinery/pkg/api/resource", Size is specified.
// It then looks at Variant for size aliases.
// Lastly, when none of them is given, default size is returned.
// The default payload size is defined in "github.com/hanapedia/hexagon/pkg/operator/constants"
func GetPayloadSize(payloadSpec PayloadSpec) int64 {
	if size, ok := ParseResource(payloadSpec); ok {
		logger.Logger.Debugf("parsed payload size, %v", size)
		return size
	}
	size, ok := constants.PayloadSizeMap[payloadSpec.Variant]
	if !ok {
		size = constants.PayloadSizeMap[constants.DefaultPayloadSize]
		logger.Logger.Debugf("resorting to default payload size, %v", size)
	}
	return size
}

func GetPayloadVariant(payloadSpec PayloadSpec) constants.PayloadSizeVariant {
	if size, ok := ParseResource(payloadSpec); ok {
		if size < constants.PayloadSizeMap[constants.SMALL] {
			return constants.SMALL
		}
		if size < constants.PayloadSizeMap[constants.MEDIUM] {
			return constants.MEDIUM
		}
		if size >= constants.PayloadSizeMap[constants.MEDIUM] {
			return constants.LARGE
		}
	}
	if payloadSpec.Variant == "" {
		return constants.DefaultPayloadSize
	}
	return payloadSpec.Variant
}
