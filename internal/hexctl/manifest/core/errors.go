package core

import (
	"fmt"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

type ManifestErrors struct {
	Stateless     []StatelessManifestError
	Broker        []BrokerManifestError
	Stateful      []StatefulManifestError
	Common        []CommonManifestError
	LoadGenerator []LoadGeneratorManifestError
}

func (me ManifestErrors) Print() {
	for _, err := range me.Stateless {
		fmt.Println(err.message)
	}
	for _, err := range me.Broker {
		fmt.Println(err.message)
	}
	for _, err := range me.Stateful {
		fmt.Println(err.message)
	}
	for _, err := range me.Common {
		fmt.Println(err.message)
	}
	for _, err := range me.LoadGenerator {
		fmt.Println(err.message)
	}
}

func (me *ManifestErrors) Extend(other ManifestErrors) {
	me.Stateless = append(me.Stateless, other.Stateless...)
	me.Broker = append(me.Broker, other.Broker...)
	me.Stateful = append(me.Stateful, other.Stateful...)
	me.Common = append(me.Common, other.Common...)
	me.LoadGenerator = append(me.LoadGenerator, other.LoadGenerator...)
}

func (me ManifestErrors) Exist() bool {
	return (len(me.Stateless) > 0 ||
		len(me.Broker) > 0 ||
		len(me.Stateful) > 0 ||
		len(me.Common) > 0 ||
		len(me.LoadGenerator) > 0)
}

// stateless error
type LoadGeneratorManifestError struct {
	message string
}

func (e *LoadGeneratorManifestError) Error() string {
	return e.message
}

func NewLoadGeneratorManifestError(serviceUnitConfig *model.ServiceUnitConfig, message string) LoadGeneratorManifestError {
	return LoadGeneratorManifestError{
		message: fmt.Sprintf(
			"Error generating config map manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}

// stateless error
type CommonManifestError struct {
	message string
}

func (e *CommonManifestError) Error() string {
	return e.message
}

func NewCommonManifestError(serviceUnitConfig *model.ServiceUnitConfig, message string) CommonManifestError {
	return CommonManifestError{
		message: fmt.Sprintf(
			"Error generating config map manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}

// stateless error
type StatelessManifestError struct {
	message string
}

func (e *StatelessManifestError) Error() string {
	return e.message
}

func NewStatelessManifestError(serviceUnitConfig *model.ServiceUnitConfig, message string) StatelessManifestError {
	return StatelessManifestError{
		message: fmt.Sprintf(
			"Error generating stateless manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}

// broker error
type BrokerManifestError struct {
	message string
}

func (e *BrokerManifestError) Error() string {
	return e.message
}

func NewBrokerManifestError(brokerAdapterConfig *model.ConsumerConfig, message string) BrokerManifestError {
	return BrokerManifestError{
		message: fmt.Sprintf(
			"Error generating broker manifest for %s for topic %s: %s",
			brokerAdapterConfig.Variant,
			brokerAdapterConfig.Topic,
			message,
		),
	}
}

func NewBrokerManifestFileError(serviceUnitConfig *model.ServiceUnitConfig, message string) BrokerManifestError {
	return BrokerManifestError{
		message: fmt.Sprintf(
			"Error generating broker manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}


// stateful error
type StatefulManifestError struct {
	message string
}

func (e *StatefulManifestError) Error() string {
	return e.message
}

func NewStatefulManifestError(serviceUnitConfig *model.ServiceUnitConfig, message string) StatefulManifestError {
	return StatefulManifestError{
		message: fmt.Sprintf(
			"Error generating stateful manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}
