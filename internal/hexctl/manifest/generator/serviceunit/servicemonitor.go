package serviceunit

import (
	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/manifest/serviceunit"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

type ServiceMonitorManifest struct {
	serviceMonitor *promv1.ServiceMonitor
}

func NewServiceMonitorManifest(suc *model.ServiceUnitConfig, cc *model.ClusterConfig) *ServiceMonitorManifest {
	manifest := ServiceMonitorManifest{
		serviceMonitor: serviceunit.CreateServiceMonitor(suc, cc),
	}

	return &manifest
}

func (smm *ServiceMonitorManifest) Generate(config *model.ServiceUnitConfig, path string) core.ManifestErrors {
	// Open the manifestFile in append mode and with write-only permissions
	file, err := core.CreateFile(path)
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Unable to open output file."),
			},
		}
	}
	defer file.Close()

	serviceMonitorYaml := yaml.GenerateManifest(smm.serviceMonitor)
	_, err = file.WriteString(core.FormatManifest(serviceMonitorYaml))
	if err != nil {
		return core.ManifestErrors{
			Stateless: []core.StatelessManifestError{
				core.NewStatelessManifestError(config, "Failed to write ServiceMonitor manifest"),
			},
		}
	}

	return core.ManifestErrors{}
}
