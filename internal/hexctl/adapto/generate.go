package adapto

import (
	"path/filepath"
	"time"

	"github.com/hanapedia/hexagon/internal/hexctl/loader"
	"github.com/hanapedia/hexagon/internal/hexctl/manifest/core"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/hexctl/graph"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/operator/yaml"
)

type AdaptoGenerator struct {
	// ServiceUnitConfigs is a map containing ServiceUnitConfigs keyed by the input path
	ServiceUnitConfigs map[string]*model.ServiceUnitConfig
	ClusterConfig      *model.ClusterConfig
	Graph              *graph.Graph
}

func GenerateFromDirectory(input, output string) {
	// Attempt to parse the cluster config file first
	clusterConfig := loader.GetClusterConfig(filepath.Join(input, constants.CLUSTER_CONFIG_FILE_NAME))

	paths, err := loader.GetYAMLFiles(input)
	if err != nil {
		logger.Logger.Fatalf("Error reading from directory %s. %s", input, err)
	}

	serviceUnitConfigs := map[string]*model.ServiceUnitConfig{}
	for _, path := range paths {
		outputPath := loader.ReplaceInputDirectory(path, input, output)
		serviceUnitConfigs[outputPath] = loader.GetServiceUnitConfig(path)
	}
	ag := AdaptoGenerator{
		ClusterConfig:      clusterConfig,
		ServiceUnitConfigs: serviceUnitConfigs,
	}
	ag.GenerateGraph()
	ag.PatchServiceUnitConfigs()
	ag.RewriteYAML()
}

// generate generates GraphML compatible Graph
func (ag *AdaptoGenerator) GenerateGraph() {
	graph := graph.Graph{}

	// Loop through ServiceUnitConfigs
	for _, suc := range ag.ServiceUnitConfigs {
		serviceName := suc.Name
		for _, primary := range suc.AdapterConfigs {
			graph.AddNode(primary.GetId(serviceName), map[string]interface{}{"config": primary})
			for _, taskSpec := range primary.TaskSpecs {
				if taskSpec.AdapterConfig.Type() == model.Invocation {
					graph.AddEdge(primary.GetId(serviceName), taskSpec.AdapterConfig.GetId())
				}
			}
		}
	}
	ag.Graph = &graph
}

// generate generates GraphML compatible Graph
func (ag *AdaptoGenerator) PatchServiceUnitConfigs() {
	// apply depth first search
	// callsMap contains the recursive call counters for each primary adapter
	// should set the adaptive timeout based on (counter + 1) * base timeout
	callsMap := ag.Graph.CalculateRecursiveCalls()
	baseTimeout, err := ag.ClusterConfig.GetBaseTimeout()
	if err != nil {
		baseTimeout = model.DEFAULT_ADAPTIVE_TIMEOUT_MAX
	}
	for _, suc := range ag.ServiceUnitConfigs {
		for _, primary := range suc.AdapterConfigs {
			for _, taskSpec := range primary.TaskSpecs {
				if taskSpec.AdapterConfig.Type() != model.Invocation {
					continue
				}
				calls, ok := callsMap[taskSpec.AdapterConfig.GetId()]
				if !ok {
					logger.Logger.Errorf("No matching adapter found.")
					continue
				}
				taskSpec.Resiliency = ag.ClusterConfig.Resiliency
				if ag.ClusterConfig.Resiliency.AdaptiveCallTimeout.Enabled {
					taskSpec.Resiliency.AdaptiveCallTimeout.LatencySLO = (baseTimeout * time.Duration(calls+1)).String()
				} else {
					taskSpec.Resiliency.CallTimeout = (baseTimeout * time.Duration(calls+1)).String()
				}
			}
		}
	}
}

// RewriteYAML rewrites service unit configs with patched fields
func (ag *AdaptoGenerator) RewriteYAML() {
	for path, suc := range ag.ServiceUnitConfigs {
		// Open the manifestFile in append mode and with write-only permissions
		file, err := core.CreateFile(path)
		if err != nil {
			logger.Logger.WithField("path", path).Errorf("Error creating output file.")
			continue
		}
		defer file.Close()

		serviceUnitConfigYaml := yaml.GenerateManifest(suc)
		_, err = file.WriteString(core.FormatManifest(serviceUnitConfigYaml))
		if err != nil {
			logger.Logger.WithField("path", path).Errorf("Error writing to output file.")
			continue
		}
	}
}
