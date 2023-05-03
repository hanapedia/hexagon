package logger

import (
	"github.com/hanapedia/the-bench/config/model"
)

func PrintErrors(cve model.ConfigValidationError) {
	for _, err := range cve.ServiceUnitFieldErrors {
		Logger.Errorf(err.Error())
	}
	for _, err := range cve.AdapterFieldErrors {
		Logger.Errorf(err.Error())
	}
	for _, err := range cve.MappingErrors {
		Logger.Errorf(err.Error())
	}
}
