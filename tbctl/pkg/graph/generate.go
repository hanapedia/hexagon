package graph

import (
	"io/ioutil"

	"github.com/hanapedia/the-bench/tbctl/lib/graphml"
	"github.com/hanapedia/the-bench/tbctl/pkg/loader"
	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/defaults"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/logger"
)

type GraphGenerator struct {
	// serviceUnitConfigs is the array of parsed service unit configs.
	serviceUnitConfigs []*model.ServiceUnitConfig
}

func newGraphGenerator(paths []string) GraphGenerator {
	configs := make([]*model.ServiceUnitConfig, 0, len(paths))
	for _, path := range paths {
		config := loader.GetConfig(path)
		defaults.SetDefauls(&config)
		configs = append(configs, &config)
	}

	return GraphGenerator{
		serviceUnitConfigs: configs,
	}
}

// generate generates GraphML compatible Graph
func (gg GraphGenerator) generate() []byte {
	graph := graphml.NewGraph()

	// Loop through ServiceUnitConfigs
	for _, config := range gg.serviceUnitConfigs {
		serviceName := config.Name
		for _, ia := range config.IngressAdapterConfigs {

			if ia.StatefulIngressAdapterConfig != nil {
				// Add node data from Egress
				continue
			}

			if ia.BrokerIngressAdapterConfig != nil {
				err := addEdge(graph, ia.BrokerIngressAdapterConfig.Topic, serviceName, "subscribe")
				if err != nil {
					logger.Logger.Fatalf(
						"Error adding edge (%s, %s). %s",
						ia.BrokerIngressAdapterConfig.Topic,
						serviceName,
						err,
					)
				}
			}

			parseSteps(graph, serviceName, ia.Steps)
			graph.SetNodeData(serviceName, "type" , "stateless")
		}
	}

	data, err := graph.ToGraphML()
	if err != nil {
		logger.Logger.Fatalf("Error parsing to GraphML. %s", err)
	}

	return data
}

// addEdge adds edge to the graph with metadata.
func addEdge(graph *graphml.Graph, source, destination, edgeLabel string) error {
	// Extra edge needs to be added for topic subscription.
	err := graph.AddEdge(source, destination)
	if err != nil {
		if err == graphml.ErrEdgeAlreadyExists {
			return nil
		}
		return err
	}
	err = graph.SetEdgeData(source, destination, "type", edgeLabel)
	if err != nil {
		if err == graphml.ErrEdgeDataAlreadySet {
			return nil
		}
		return err
	}

	switch edgeLabel {
	case "tcp":
		err = graph.SetNodeData(destination, "type", "stateful")
	case "http":
		err = graph.SetNodeData(destination, "type", "stateless")
	case "publish":
		err = graph.SetNodeData(destination, "type", "broker")
	case "subscribe":
		err = graph.SetNodeData(source, "type", "broker")
	}
	if err != nil {
		if err == graphml.ErrNodeDataAlreadySet {
			return nil
		}
		return err
	}

	return nil
}

// parseSteps parses steps and update graph.
func parseSteps(graph *graphml.Graph, serviceName string, steps []model.Step) {
	for _, step := range steps {
		var destination, edgeLabel string
		if step.EgressAdapterConfig.BrokerEgressAdapterConfig != nil {
			destination = step.EgressAdapterConfig.BrokerEgressAdapterConfig.Topic
			edgeLabel = "publish"
		}

		if step.EgressAdapterConfig.StatefulEgressAdapterConfig != nil {
			destination = step.EgressAdapterConfig.StatefulEgressAdapterConfig.Name
			edgeLabel = "tcp"
		}

		if step.EgressAdapterConfig.StatelessEgressAdapterConfig != nil {
			destination = step.EgressAdapterConfig.StatelessEgressAdapterConfig.Service
			edgeLabel = "http"
		}

		err := addEdge(
			graph,
			serviceName,
			destination,
			edgeLabel,
		)
		if err != nil {
		logger.Logger.Fatalf(
			"Error adding edge (%s, %s). %s",
			serviceName,
			destination,
			err,
		)
		}
	}
}

func GenerateFromDirectory(input, output string) {
	paths, err := loader.GetYAMLFiles(input)
	if err != nil {
		logger.Logger.Fatalf("Error reading from directory %s. %s", input, err)
	}

	gg := newGraphGenerator(paths)
	data := gg.generate()

	err = ioutil.WriteFile(output, data, 0644)
	if err != nil {
		logger.Logger.Fatalf("Error writing GraphML to a file. %s", err)
	}
}
