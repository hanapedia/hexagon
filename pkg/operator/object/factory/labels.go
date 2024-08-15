package factory

import "fmt"

const HostnameLabel = "kubernetes.io/hostname"

const LabelPrefix = "hexagon.hanapedia.link"

var AppLabel = fmt.Sprintf("%v/app", LabelPrefix)

var ServiceUnitVariantLabel = fmt.Sprintf("%v/variant", LabelPrefix)

type ServiceUnitVariantValues float64

const (
	STATELESS ServiceUnitVariantValues = iota
	STATEFULL
	LOAD_GENERATOR
)

func (suvv ServiceUnitVariantValues) AsString() string {
	switch suvv {
	case STATELESS:
		return "stateless"
	case STATEFULL:
		return "statefull"
	case LOAD_GENERATOR:
		return "LOAD_GENERATOR"
	}
	return "stateless"
}

