# Configuration language
see [terminologies](./docs/terminology.md) for unclear terms.
## Basics
- Each service unit can be configured using yaml.
- Each service unit can be defined by configuring set of primary adapters that is exposes to other serivce units.
- Each primary adapter can be defined by assigning an array of secondary adapters.
- Both primary and secondary adapters can be a variant of stateless, stateful, or broker. see [Configuration](#configuration)for possible values.

## Configuration
### Service Unit
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| name          | Name of the service.                          | ""          | true        |
| version       | Version of the image used.                    | ""          | true        |
| deployment    | [Deployment configs](#deployment).            | {}          | false       |
| adapters      | List of [Primary Adapter Configs](#primary-adapter).| []    | true        |

### Deployment
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|
| replicas      | Number replicas for the service.              | 1           | false       |
| gateway       | Configuration for load generator. Should be defined if you want to enable load generator for the service. | {} | false |
| gateway.virtualUsers | Number of virtual users.               | 0           | true        |
| gateway.duration     | Duration of the load test in minutes.  | 0           | true        |
| resources     | Resource limit and request in k8s core v1 format.    | {}      | false    |
| env           | Extra environmental variables in k8s core v1 format. | {}      | false    |

### Primary Adapter
| Parameter     | Description                                   | Default     | Required    |
|---------------|-----------------------------------------------|-------------|-------------|

## Validation
- configuration file written in yaml can be validated using the [cli](../cmd/hexctl/).
- two types of validations:
    - field validation: checks if fields of the configuration file is valid.
    - mapping validation: checks if secondary adapter can be mapped to the destination primary adapter.
- field validation is also ran inside of Service unit once the config file is loaded.
