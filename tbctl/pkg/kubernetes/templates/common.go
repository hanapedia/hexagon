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
  name: {{ .Name }}-config
  namespace: {{ .Namespace }}
data:
  rawYAMLContent: |
{{ indent .Config 4 }}
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
  name: {{ .Name }}-env
  namespace: {{ .Namespace }}
data:
  fromEnvFile: |
{{ indent .Envs 4 }}
`
