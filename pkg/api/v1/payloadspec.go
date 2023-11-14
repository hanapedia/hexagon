package v1

import (
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"k8s.io/apimachinery/pkg/api/resource"
)

// PayloadSpec contains configuration for payload sent by adapters
type PayloadSpec struct {
	Size    string                       `json:"size,omitempty"`
	Variant constants.PayloadSizeVariant `json:"class,omitempty" validate:"omitempty,oneof=small medium large"`
	Count   int                          `json:"payloadCount,omitempty"`
}

func GetPayloadSize(payloadSpec PayloadSpec) int64 {
	if payloadSpec.Size != "" {
		quantity, err := resource.ParseQuantity(payloadSpec.Size)
		if err == nil {
			if size, ok := quantity.AsInt64(); ok {
				return size
			}
		}
	}
	size, ok := constants.PayloadSizeMap[payloadSpec.Variant]
	if !ok {
		size = constants.PayloadSizeMap[constants.DefaultPayloadSize]
	}
	return size
}
