package defaults

const (
	LOAD_GENERATOR_IMAGE_NAME       = "hiroki11hanada/tb-load-generator"
	LIMIT_MEM_LG                    = "1Gi"
	LG_K6_PROMETHEUS_RW_SERVER_URL  = "http://prometheus-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090/api/v1/write"
	LG_K6_PROMETHEUS_RW_TREND_STATS = "p(95),p(99),avg"
)
