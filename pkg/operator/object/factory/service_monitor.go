package factory

import (
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

const METRICS_PORT_NAME = "metrics"

type ServiceMonitorArgs struct {
	Name      string
	Namespace string
}

// NewServiceMonitor creates ServiceMonitor Custom Resource defined by
// github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring
func NewServiceMonitor(args *ServiceMonitorArgs) promv1.ServiceMonitor {
	return promv1.ServiceMonitor{
		TypeMeta: NewTypeMeta("ServiceMonitor", "monitoring.coreos.com/v1"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{
			Name:      args.Name,
			Namespace: args.Namespace,
			Labels: map[string]string{
				AppLabel: args.Name,
			},
		},
		),
		Spec: promv1.ServiceMonitorSpec{
			Selector: *NewLabelSelector(map[string]string{AppLabel: args.Name}),
			Endpoints: []promv1.Endpoint{
				{Port: METRICS_PORT_NAME},
			},
		},
	}
}
