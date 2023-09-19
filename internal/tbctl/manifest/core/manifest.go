package core

import (
	v1 "github.com/hanapedia/the-bench/pkg/api/v1"
)

// Manifest interface should be implemented by structs that represent
// kubernetes resources required by entities such as service unit
type Manifest interface {
	// Generate generates manifest in
	Generate(*v1.ServiceUnitConfig, string) ManifestErrors
}
