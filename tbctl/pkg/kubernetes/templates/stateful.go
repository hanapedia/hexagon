package templates

type StatefulManifestTemplateArgs struct {
	Name                   string
	Namespace              string
	Image                  string
	Replicas               int
	ResourceLimitsCPU      string
	ResourceLimitsMemory   string
	ResourceRequestsCPU    string
	ResourceRequestsMemory string
	MONGOPort              int
	POSTGREPort            int
}

const StatefulManifestTemplates = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: {{ .Name }}
  template:
    metadata:
      labels:
        app: {{ .Name }}
    spec:
      containers:
      - name: {{ .Name }}
        image: {{ .Image }}
        ports:
        - containerPort: {{ .MONGOPort }}
          name: mongo
        envFrom:
        - configMapRef:
            name: {{ .Name }}-env-configmap
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  selector:
    app: {{ .Name }}
  ports:
  - name: mongo
    port: {{ .MONGOPort }}
    targetPort: {{ .MONGOPort }}
  - name: postgre
    port: {{ .POSTGREPort }}
    targetPort: {{ .POSTGREPort }}
  type: ClusterIP
`
