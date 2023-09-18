# Configuration language
see [terminologies](./docs/terminology.md) for unclear terms.
## Basics
- The configuration is done using yaml for each service unit.
- Each service unit can be defined by configuring set of primary adapters that is exposes to other serivce units.
- Each primary adapter can be defined by assigning an array of secondary adapters.
- Bothe primary and secondary adapters can be a variant of stateless, stateful, or broker

## Docs (coming soon)

## Validation
- configuration file written in yaml can be validated using the [cli](../cmd/tbctl/).
- two types of validations:
    - field validation: checks if fields of the configuration file is valid.
    - mapping validation: checks if secondary adapter can be mapped to the destination primary adapter.
- field validation is also ran inside of Service unit once the config file is loaded.
