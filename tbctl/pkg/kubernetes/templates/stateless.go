package templates

type StatelessManifestTemplateArgs struct {
	Name                   string
	Namespace              string
	Image                  string
	Replicas               int
	ResourceLimitsCPU      string
	ResourceLimitsMemory   string
	ResourceRequestsCPU    string
	ResourceRequestsMemory string
	HTTPPort  int
	GRPCPort  int
}

const StatelessManifestTemplates = `---
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
        resources:
          limits:
            cpu: {{ .ResourceLimitsCPU }}
            memory: {{ .ResourceLimitsMemory }}
          requests:
            cpu: {{ .ResourceRequestsCPU }}
            memory: {{ .ResourceRequestsMemory }}
        ports:
        - containerPort: {{ .HTTPPort }}
          name: http
        - containerPort: {{ .GRPCPort }}
          name: grpc
        volumeMounts:
        - name: config
          mountPath: /app/config/
      volumes:
      - name: config
        configMap:
          name: {{ .Name }}-config
          items:
          - key: rawYAMLContent
            path: service-unit.yaml
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
  - name: http
    protocol: TCP
    port: {{ .HTTPPort }}
    targetPort: {{ .HTTPPort }}
  - name: grpc
    protocol: TCP
    port: {{ .GRPCPort }}
    targetPort: {{ .GRPCPort }}
  type: ClusterIP
`
