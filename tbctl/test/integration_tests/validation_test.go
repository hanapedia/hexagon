package tests

import (
	"testing"

	"github.com/hanapedia/the-bench/tbctl/pkg/validation"
)

func TestFileValidation(t *testing.T) {
	validation.ValidateFile("./testdata/service-a.yaml")
}

func TestDirctoryFlatValidation(t *testing.T) {
	validation.ValidateDirectory("./testdata/flat/")
}

func TestDirctoryNestedValidation(t *testing.T) {
	validation.ValidateDirectory("./testdata/nested/")
}

func TestDirctoryInvalidFieldRejection(t *testing.T) {
	sufe, aef := validation.ValidateFile("./testdata/invalid_field/service-invalid.yaml")
	if len(sufe) == 0 && len(aef) == 0 {
		t.Errorf("The function did not return an error")
	}
}

func TestDirctoryInvalidMappingRejection(t *testing.T) {
	err := validation.ValidateDirectory("./testdata/invalid_mapping/")
	if len(err.MappingErrors) == 0 {
		t.Errorf("The function did not return an error")
	}
}
