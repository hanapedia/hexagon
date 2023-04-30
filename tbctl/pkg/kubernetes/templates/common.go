package templates

type ConfigConfigMap struct {
	Name      string
	Namespace string
	Config    string
}

// indent function must be defined before using this template
const ConfigConfigMapTemplate = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Name }}-config-configMap
  namespace: {{ .Namespace }}
data:
{{ indent .Config 2 }}
`

type EnvConfigMap struct {
	Name      string
	Namespace string
	Envs      string
}

const EnvConfigMapTemplate = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Name }}-env-configMap
  namespace: {{ .Namespace }}
data:
  fromEnvFile: |-
{{ indent .Envs 4 }}
`
