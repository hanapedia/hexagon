package v1

import (
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

func TestGetPayloadSize(t *testing.T) {
	tests := []struct {
		name     string
		spec     PayloadSpec
		expected int64
	}{
		{
			name: "Valid Size",
			spec: PayloadSpec{
				Size: "10Gi",
			},
			expected: 10 * 1024 * 1024 * 1024, // 10 GiB in bytes
		},
		{
			name: "Invalid Size",
			spec: PayloadSpec{
				Size: "invalid-size",
			},
			expected: constants.PayloadSizeMap[constants.DefaultPayloadSize], // Default size
		},
		{
			name: "1",
			spec: PayloadSpec{
				Size: "1",
			},
			expected: 1,
		},
		{
			name: "0",
			spec: PayloadSpec{
				Size: "0",
			},
			expected: 0,
		},
		{
			name: "Valid Variant",
			spec: PayloadSpec{
				Variant: constants.SMALL,
			},
			expected: constants.PayloadSizeMap[constants.SMALL],
		},
		{
			name: "Empty Size and Variant",
			spec: PayloadSpec{},
			expected: constants.PayloadSizeMap[constants.DefaultPayloadSize], // Default size
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPayloadSize(tt.spec)
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}
