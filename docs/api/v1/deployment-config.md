## Deployment
replicas, resources and envs are of same type as the parameters for Kubernetes Deployment API Object, and they are passed straight through when generating Kubernetes manifests.
The definitions for Go struct is in [`pkg/api/v1/deploymentspec.go`](../../../pkg/api/v1/deploymentspec.go)

| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| replicas      | Number replicas for the service.              | 1           | false       |
| gateway       | Configuration for load generator. Should be defined if you want to enable load generator for the service. [see](#gateway) | {} | false |
| resources     | Resource limit and request in k8s core v1 format.    | {}      | false    |
| env           | Extra environmental variables in k8s core v1 format. | {}      | false    |
| topologySpreadConstraint | Whether to turn on topologySpreadConstraint. | false      | false    |

\* TopologySpreadConstraint is applied with following values if turned on.
This feature is just for generating manifests with sane default. Adjust if needed using tools like Kustomize
```yaml
    topologySpreadConstraints:
        - maxSkew: 1
          topologyKey: kubernetes.io/hostname
          whenUnsatisfiable: ScheduleAnyway
          labelSelector:
            matchLabels:
              hexagon.hanapedia.link/app: <<Service Unit Name>>
```

### Gateway
| Parameter     | Description                           | Default     | Required    |
|---------------|---------------------------------------|-------------|-------------|
| virtualUsers | Number of virtual users.               | 0           | true        |
| duration     | Duration of the load test in minutes.  | 0           | true        |
