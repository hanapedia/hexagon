package generate

import (
	"fmt"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
)

type ManifestErrors struct {
	stateless []StatelessManifestError
	broker    []BrokerManifestError
	stateful  []StatefulManifestError
	common    []CommonManifestError
}

func (me ManifestErrors) Print() {
	for _, err := range me.stateless {
		fmt.Println(err.message)
	}
	for _, err := range me.broker {
		fmt.Println(err.message)
	}
	for _, err := range me.stateful {
		fmt.Println(err.message)
	}
	for _, err := range me.common {
		fmt.Println(err.message)
	}
}

func (me *ManifestErrors) Extend(other ManifestErrors) {
	me.stateless = append(me.stateless, other.stateless...)
	me.broker = append(me.broker, other.broker...)
	me.stateful = append(me.stateful, other.stateful...)
	me.common = append(me.common, other.common...)
}

func (me ManifestErrors) Exist() bool {
	return len(me.stateless) > 0 || len(me.broker) > 0 || len(me.stateful) > 0 || len(me.common) > 0
}

// stateless error
type CommonManifestError struct {
	message string
}

func (e *CommonManifestError) Error() string {
	return e.message
}

func NewCommonManifestError(serviceUnitConfig model.ServiceUnitConfig, message string) CommonManifestError {
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

func NewStatelessManifestError(serviceUnitConfig model.ServiceUnitConfig, message string) StatelessManifestError {
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

func NewBrokerManifestError(brokerAdapterConfig model.BrokerIngressAdapterConfig, message string) BrokerManifestError {
	return BrokerManifestError{
		message: fmt.Sprintf(
			"Error generating broker manifest for %s for topic %s: %s",
			brokerAdapterConfig.Variant,
			brokerAdapterConfig.Topic,
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

func NewStatefulManifestError(serviceUnitConfig model.ServiceUnitConfig, message string) StatefulManifestError {
	return StatefulManifestError{
		message: fmt.Sprintf(
			"Error generating stateful manifest for %s: %s",
			serviceUnitConfig.Name,
			message,
		),
	}
}
