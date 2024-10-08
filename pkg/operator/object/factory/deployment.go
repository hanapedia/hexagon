package factory

import (
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/health"
	"github.com/hanapedia/hexagon/pkg/operator/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentArgs struct {
	Name                           string
	Namespace                      string
	Annotations                    map[string]string
	Image                          string
	Replicas                       int32
	Resource                       *corev1.ResourceRequirements
	Ports                          map[string]int32
	VolumeMounts                   map[string]string
	Envs                           []corev1.EnvVar
	ConfigVolume                   *ConfigMapVolumeArgs
	EnableTopologySpreadConstraint bool
	DisableReadinessProbe          bool
	ReadinessProbePort             intstr.IntOrString
}

type ConfigMapVolumeArgs struct {
	Name  string
	Items map[string]string
}

// NewDeployment create Deployment
func NewDeployment(args *DeploymentArgs) appsv1.Deployment {
	return appsv1.Deployment{
		TypeMeta:   NewTypeMeta("Deployment", "apps/v1"),
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace, Annotations: args.Annotations}),
		Spec:       NewDeploymentSpec(args),
	}
}

// NewDeploymentSpec create deployment specs
func NewDeploymentSpec(args *DeploymentArgs) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: utils.Int32Ptr(args.Replicas),
		Selector: NewLabelSelector(map[string]string{AppLabel: args.Name}),
		Template: NewPodTemplate(args),
	}

}

// NewPodTemplate create pod template
func NewPodTemplate(args *DeploymentArgs) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{Labels: map[string]string{AppLabel: args.Name}}),
		Spec:       NewPodSpec(args),
	}
}

// NewPodSpec create pod spec
func NewPodSpec(args *DeploymentArgs) corev1.PodSpec {
	return corev1.PodSpec{
		Containers:                NewContainer(args),
		Volumes:                   NewVolume(args.ConfigVolume),
		TopologySpreadConstraints: NewTopologySpreadConstraint(args),
	}
}

// NewContainer create container
func NewContainer(args *DeploymentArgs) []corev1.Container {
	resources := corev1.ResourceRequirements{}
	if args.Resource != nil {
		resources = *args.Resource
	}
	return []corev1.Container{
		{
			Name:            args.Name,
			Image:           args.Image,
			ImagePullPolicy: corev1.PullAlways,
			Resources:       resources,
			Ports:           NewContainerPort(args.Ports),
			VolumeMounts:    NewVolumeMount(args.VolumeMounts),
			Env:             args.Envs,
			ReadinessProbe:  NewReadinessProbe(args),
		},
	}
}

// NewVolume create volume
func NewVolume(configVolumeArgs *ConfigMapVolumeArgs) []corev1.Volume {
	var volumes []corev1.Volume
	if configVolumeArgs != nil {
		volumes = append(volumes, GetConfigMapVolume("config", configVolumeArgs))
	}
	return volumes
}

// TODO: parameterize magic number
func NewTopologySpreadConstraint(args *DeploymentArgs) []corev1.TopologySpreadConstraint {
	var tsc []corev1.TopologySpreadConstraint
	if args.EnableTopologySpreadConstraint {
		tsc = append(tsc, corev1.TopologySpreadConstraint{
			MaxSkew:           1,
			TopologyKey:       HostnameLabel,
			WhenUnsatisfiable: corev1.ScheduleAnyway,
			LabelSelector:     NewLabelSelector(map[string]string{AppLabel: args.Name}),
		})
	}
	return tsc
}

// TODO: parameterize magic number
func NewReadinessProbe(args *DeploymentArgs) *corev1.Probe {
	if args.DisableReadinessProbe {
		return nil
	}
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
			Path: health.HEALTH_PATH,
			Port: args.ReadinessProbePort,
		}},
		InitialDelaySeconds: 5,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
	}
}

// GetConfigMapVolume creates the cofigmap volume entry
func GetConfigMapVolume(key string, arg *ConfigMapVolumeArgs) corev1.Volume {
	var items []corev1.KeyToPath
	for key, path := range arg.Items {
		item := corev1.KeyToPath{
			Key:  key,
			Path: path,
		}
		items = append(items, item)
	}
	volume := corev1.Volume{
		Name: key,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: *NewLocalObjectReference(arg.Name),
				Items:                items,
			},
		},
	}
	return volume
}

// NewContainerPort create container port slice
func NewContainerPort(ports map[string]int32) []corev1.ContainerPort {
	var containerPorts []corev1.ContainerPort
	for name, port := range ports {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          name,
			ContainerPort: port,
		})
	}
	return containerPorts
}

// NewVolumeMount create container port slice
func NewVolumeMount(volumes map[string]string) []corev1.VolumeMount {
	var volumeMounts []corev1.VolumeMount
	for name, path := range volumes {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      name,
			MountPath: path,
		})
	}
	return volumeMounts
}
