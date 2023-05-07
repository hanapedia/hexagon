package model

import "github.com/go-playground/validator/v10"

func (suc ServiceUnitConfig) Validate() []InvalidServiceUnitFieldValueError {
	validate := validator.New()
	var errs []InvalidServiceUnitFieldValueError
	err := validate.Struct(suc)
	if err != nil {
		errs = append(errs, mapInvalidServiceUnitFieldValueErrors(err, suc)...)
	}

	return errs
}

func (step Step) Validate(serviceName string, ingressAdapterConfig IngressAdapterConfig) []InvalidStepFieldValueError {
	var stepFieldErrors []InvalidStepFieldValueError
	if step.EgressAdapterConfig == nil {
		stepFieldErrors = append(stepFieldErrors, NewInvalidStepFieldValueError(ingressAdapterConfig.GetId(serviceName)))
	}

	return stepFieldErrors
}

func (iac IngressAdapterConfig) GetId(serviceName string) string {
	var id string
	if iac.StatelessIngressAdapterConfig != nil {
		id = iac.StatelessIngressAdapterConfig.GetId(serviceName)
	}
	if iac.BrokerIngressAdapterConfig != nil {
		id = iac.BrokerIngressAdapterConfig.GetId(serviceName)
	}
	if iac.StatefulIngressAdapterConfig != nil {
		id = iac.StatefulIngressAdapterConfig.GetId(serviceName)
	}
	return id
}

func (eac EgressAdapterConfig) GetId() string {
	var id string
	if eac.StatelessEgressAdapterConfig != nil {
		id = eac.StatelessEgressAdapterConfig.GetId()
	}
	if eac.BrokerEgressAdapterConfig != nil {
		id = eac.BrokerEgressAdapterConfig.GetId()
	}
	if eac.StatefulEgressAdapterConfig != nil {
		id = eac.StatefulEgressAdapterConfig.GetId()
	}
	return id
}
