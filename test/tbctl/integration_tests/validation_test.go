package integration_test

import (
	"testing"

	"github.com/hanapedia/hexagon/internal/tbctl/validation"
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

func TestDirctoryInvalidServiceFieldRejection(t *testing.T) {
	errs := validation.ValidateFile("./testdata/invalid_service_unit/service-unit.yaml")
	if errs.Exist() {
		t.Errorf("The function did not return an error")
	}
}

func TestDirctoryInvalidFieldRejection(t *testing.T) {
	errs := validation.ValidateFile("./testdata/invalid_field/service-invalid.yaml")
	if errs.Exist() {
		t.Errorf("The function did not return an error")
	}
}

func TestDirctoryInvalidMappingRejection(t *testing.T) {
	err := validation.ValidateDirectory("./testdata/invalid_mapping/")
	if len(err.MappingErrors) == 0 {
		t.Errorf("The function did not return an error")
	}
}
