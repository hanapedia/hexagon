# Breaking changes
## version 0.0.5 -> 0.1.0
### Field name changes
- `.ingressAdapters` -> `adapters`
- `.ingressAdapters.0.broker` -> `.adapters.0.consumer`
- `.ingressAdapters.0.stateless` -> `.adapters.0.server`
- `.ingressAdapters.0.stateful` -> `.adapters.0.repository`
- `.ingressAdapters.0.steps.0.egressAdapter.broker` -> `.ingressAdapters.0.steps.0.adapter.producer`
- `.ingressAdapters.0.steps.0.egressAdapter.stateless` -> `.ingressAdapters.0.steps.0.adapter.invocation`
- `.ingressAdapters.0.steps.0.egressAdapter.stateful` -> `.ingressAdapters.0.steps.0.adapter.repository`

### New Fields
- `deployment`: config for deployment
- `deployment.gateway`: config for Gateway moved
- `deployment.resources`: config for resource allocation
- `deployment.env`: config for additional env vars
