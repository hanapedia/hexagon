# Configuration language
see [terminologies](./docs/terminology.md) for unclear terms.
## Basics
- The configuration is done using yaml for each service unit.
- Each service unit can be defined by configuring set of ingress adapters that is exposes to other serivce units.
- Each ingress adapter can be defined by assigning an array of egress adapters.
- Bothe ingress and egress adapters can be a variant of stateless, stateful, or broker

## Docs (coming soon)

## Validation
- configuration file written in yaml can be validated using the [cli](../tbctl/).
- two types of validations:
    - field validation: checks if fields of the configuration file is valid.
    - mapping validation: checks if egress adapter can be mapped to the destination ingress adapter.
- field validation is also ran inside of Service unit once the config file is loaded.
