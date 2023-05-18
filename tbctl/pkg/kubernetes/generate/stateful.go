package generate

import (
	"errors"
	"fmt"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/templates"
)

type StatefulVaraintDetails struct {
	image string
	env   string
}

func GenerateStatefulManifests(dir string, serviceUnitConfig model.ServiceUnitConfig) ManifestErrors {
	var statefulManifestErrors []StatefulManifestError
	variantDetails, err := getVariantDetails(serviceUnitConfig.IngressAdapterConfigs)
	if err != nil {
		statefulManifestErrors = append(
			statefulManifestErrors,
			NewStatefulManifestError(serviceUnitConfig, err.Error()),
		)
	}
	statefulManifestTemplateArgs := templates.StatefulManifestTemplateArgs{
		Name:                   serviceUnitConfig.Name,
		Namespace:              NAMESPACE,
		Image:                  variantDetails.image,
		Replicas:               REPLICAS,
		ResourceLimitsCPU:      LIMIT_CPU,
		ResourceLimitsMemory:   LIMIT_MEM,
		ResourceRequestsCPU:    REQUEST_CPU,
		ResourceRequestsMemory: REQUEST_MEM,
		MONGOPort:              MONGO_PORT,
		POSTGREPort:            POSTGRE_PORT,
	}

	// generate deployment and service manifest
	err = RenderAndSave(
		dir,
		fmt.Sprintf("%s-manifest", serviceUnitConfig.Name),
		templates.StatefulManifestTemplates,
		statefulManifestTemplateArgs,
	)
	if err != nil {
		statefulManifestErrors = append(
			statefulManifestErrors,
			NewStatefulManifestError(serviceUnitConfig, err.Error()),
		)
	}

	var commonManifestErrors []CommonManifestError
	envErrs := GenerateEnvManifest(dir, serviceUnitConfig, variantDetails.env)
	if envErrs != nil {
		commonManifestErrors = append(commonManifestErrors, NewCommonManifestError(serviceUnitConfig, envErrs.Error()))
	}

	return ManifestErrors{stateful: statefulManifestErrors, common: commonManifestErrors}
}

func getVariantDetails(ingressAdapterConfigs []model.IngressAdapterSpec) (StatefulVaraintDetails, error) {
	for _, ingresingressAdapterConfig := range ingressAdapterConfigs {
		if ingresingressAdapterConfig.StatefulIngressAdapterConfig != nil {
			switch ingresingressAdapterConfig.StatefulIngressAdapterConfig.Variant {
			case constants.MONGO:
				return StatefulVaraintDetails{image: MONGO_IMAGE, env: templates.MongoEnvs}, nil
			case constants.POSTGRE:
				return StatefulVaraintDetails{image: POSTGRE_IMAGE, env: templates.PostgreEnvs}, nil
			}
		}
	}
	return StatefulVaraintDetails{}, errors.New("No stateful ingress apdapter found.")
}
