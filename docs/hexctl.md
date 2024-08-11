## hexctl
`hexctl` is the CLI tool used to generate Kubernetes manifest from Hexagon configuration files.

### Commands
`hexctl generate`: generate Kubernetes manifests
- `-f`: **directory path** for input Hexagon configuration files
    - this is the directory that CLI tries to look for `hexagon-config.yaml` file for global/cluster configurations (not required).
    - CLI looks for all `yaml` and `yml` files under this directory recursively, and *treats each file as a configuration for a single service*. 
- `-o`: **directory path** for output Kubernetes manifest files
    - the generated files will have the same name as the input file but with a suffix such as `-service-unit.yaml`, or `-broker.yaml`
        - `-service-unit.yaml` will contain all the Kubernetes object definitions required to run the service. 
        - `-broker.yaml` will contain `Topic` custom resource, which is defined by [Strimzi Kafka Operator](https://strimzi.io/). 

`hexctl validate`: validates Hexagon configuration files.
- `-f`: **directory path or file path** for input Hexagon configuration file(s).
    - if **file path** is give, it only validates the fields of the configuration.
    - if **directory path** is given, it validates the fields of all the configurations, and also the mapping between primary adapters and secondary adapters.
        - it would check if corresponding primary adapter definition is present for each secondary adapter definition.

`hexctl graph`: (experimental) Generates [GraphML](http://graphml.graphdrawing.org/) for the application service graph.
- `-f`: **directory path** for input Hexagon configuration file(s).
- `-o`: **directory path** for output GraphML file.
