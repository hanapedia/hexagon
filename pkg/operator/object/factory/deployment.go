package factory

import (
	"github.com/hanapedia/the-bench/pkg/operator/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DeploymentArgs struct {
	Name                   string
	Namespace              string
	Annotations            map[string]string
	Image                  string
	Replicas               int32
	ResourceLimitsCPU      string
	ResourceLimitsMemory   string
	ResourceRequestsCPU    string
	ResourceRequestsMemory string
	Ports                  map[string]int32
	VolumeMounts           map[string]string
	ConfigVolume           *ConfigMapVolumeArgs
	EnvVolume              *ConfigMapVolumeArgs
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
		Selector: NewLabelSelector(map[string]string{"app": args.Name}),
		Template: NewPodTemplate(args),
	}

}

// NewPodTemplate create pod template
func NewPodTemplate(args *DeploymentArgs) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: NewObjectMeta(ObjectMetaOptions{Labels: map[string]string{"app": args.Name}}),
		Spec:       NewPodSpec(args),
	}
}

// NewPodSpec create pod spec
func NewPodSpec(args *DeploymentArgs) corev1.PodSpec {
	return corev1.PodSpec{
		Containers: NewContainer(args),
		Volumes:    NewVolume(args.ConfigVolume, args.EnvVolume),
	}
}

// NewContainer create container
func NewContainer(args *DeploymentArgs) []corev1.Container {
	return []corev1.Container{
		{
			Name:         args.Name,
			Image:        args.Image,
			Resources:    NewContainerResources(args),
			Ports:        NewContainerPort(args.Ports),
			VolumeMounts: NewVolumeMount(args.VolumeMounts),
			// EnvFrom:      ContainerEnvFactory(args),
		},
	}
}

// NewContainerResources creates resource requiremnts definition
func NewContainerResources(args *DeploymentArgs) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(args.ResourceLimitsCPU),
			corev1.ResourceMemory: resource.MustParse(args.ResourceLimitsMemory),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(args.ResourceRequestsCPU),
			corev1.ResourceMemory: resource.MustParse(args.ResourceRequestsMemory),
		},
	}
}

// NewContainerEnv creates container env from source
func NewContainerEnv(args *DeploymentArgs) []corev1.EnvFromSource {
	envFromSource := corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: *NewLocalObjectReference("env"),
		},
	}
	return []corev1.EnvFromSource{envFromSource}
}

// NewVolume create volume
func NewVolume(configVolumeArgs, envVolumeArgs *ConfigMapVolumeArgs) []corev1.Volume {
	var volumes []corev1.Volume
	if configVolumeArgs != nil {
		volumes = append(volumes, GetConfigMapVolume("config", configVolumeArgs))
	}
	if envVolumeArgs != nil {
		volumes = append(volumes, GetConfigMapVolume("env", envVolumeArgs))
	}
	return volumes
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
