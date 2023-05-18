package generate

import (
	"fmt"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/templates"
)

func GenerateStatelessManifests(input string, dir string, serviceUnitConfig model.ServiceUnitConfig) ManifestErrors {
	var statelessManifestErrors []StatelessManifestError
	templateArgs := templates.StatelessManifestTemplateArgs{
		Name:                   serviceUnitConfig.Name,
		Namespace:              NAMESPACE,
		Image:                  SERVICE_UNIT_IMAGE,
		Replicas:               REPLICAS,
		ResourceLimitsCPU:      LIMIT_CPU,
		ResourceLimitsMemory:   LIMIT_MEM,
		ResourceRequestsCPU:    REQUEST_CPU,
		ResourceRequestsMemory: REQUEST_MEM,
		HTTPPort:               HTTP_PORT,
		GRPCPort:               GRPC_PORT,
	}
	err := RenderAndSave(
		dir,
		fmt.Sprintf("%s-manifest", serviceUnitConfig.Name),
		templates.StatelessManifestTemplates,
		templateArgs,
	)
	if err != nil {
		statelessManifestErrors = append(
			statelessManifestErrors,
			NewStatelessManifestError(serviceUnitConfig, err.Error()),
		)
	}

	var commonManifestErrors []CommonManifestError
	configErrs := GenerateConfigManifest(dir, serviceUnitConfig, input)
	if configErrs != nil {
		commonManifestErrors = append(commonManifestErrors, NewCommonManifestError(serviceUnitConfig, configErrs.Error()))
	}

	return ManifestErrors{stateless: statelessManifestErrors, common: commonManifestErrors}
}
