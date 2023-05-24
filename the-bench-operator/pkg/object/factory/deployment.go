package factory

import (
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DeploymentArgs struct {
	Name                   string
	Namespace              string
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

// DeploymentFactory create Deployment
func DeploymentFactory(args *DeploymentArgs) appsv1.Deployment {
	return appsv1.Deployment{
		TypeMeta:   TypeMetaFactory("Deployment", "apps/v1"),
		ObjectMeta: ObjectMetaFactory(args.Name, args.Namespace, map[string]string{}),
		Spec:       DeploymentSpecFactory(args),
	}
}

// DeploymentSpecFactory create deployment specs
func DeploymentSpecFactory(args *DeploymentArgs) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: utils.Int32Ptr(args.Replicas),
		Selector: LabelSelectorFactory(map[string]string{"app": args.Name}),
		Template: PodTemplateFactory(args),
	}

}

// PodTemplateFactory create pod template
func PodTemplateFactory(args *DeploymentArgs) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: ObjectMetaFactory("", "", map[string]string{"app": args.Name}),
		Spec:       PodSpecFactory(args),
	}
}

// PodSpecFactory create pod spec
func PodSpecFactory(args *DeploymentArgs) corev1.PodSpec {
	return corev1.PodSpec{
		Containers: ContainerFactory(args),
		Volumes:    VolumeFactory(args.ConfigVolume, args.EnvVolume),
	}
}

// ContainerFactory create container
func ContainerFactory(args *DeploymentArgs) []corev1.Container {
	return []corev1.Container{
		{
			Name:         args.Name,
			Image:        args.Image,
			Resources:    ContainerResourcesFactory(args),
			Ports:        ContainerPortFactory(args.Ports),
			VolumeMounts: VolumeMountFactory(args.VolumeMounts),
			EnvFrom:      ContainerEnvFactory(args),
		},
	}
}

// ContainerResourcesFactory creates resource requiremnts definition
func ContainerResourcesFactory(args *DeploymentArgs) corev1.ResourceRequirements {
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

// ContainerEnvFactory creates container env from source
func ContainerEnvFactory(args *DeploymentArgs) []corev1.EnvFromSource {
	envFromSource := corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: *LocalObjectReferenceFactory("env"),
		},
	}
	return []corev1.EnvFromSource{envFromSource}
}

// VolumeFactory create volume
func VolumeFactory(configVolumeArgs, envVolumeArgs *ConfigMapVolumeArgs) []corev1.Volume {
	var volumes []corev1.Volume
	if configVolumeArgs != nil {
		volumes = append(volumes, GetConfigMapVolume(configVolumeArgs))
	}
	if configVolumeArgs != nil {
		volumes = append(volumes, GetConfigMapVolume(envVolumeArgs))
	}
	return volumes
}

// GetConfigMapVolume creates the cofigmap volume entry
func GetConfigMapVolume(arg *ConfigMapVolumeArgs) corev1.Volume {
	var items []corev1.KeyToPath
	for key, path := range arg.Items {
		item := corev1.KeyToPath{
			Key:  key,
			Path: path,
		}
		items = append(items, item)
	}
	volume := corev1.Volume{
		Name: arg.Name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: items,
			},
		},
	}
	return volume
}

// ContainerPortFactory create container port slice
func ContainerPortFactory(ports map[string]int32) []corev1.ContainerPort {
	var containerPorts []corev1.ContainerPort
	for name, port := range ports {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          name,
			ContainerPort: port,
		})
	}
	return containerPorts
}

// VolumeMountFactory create container port slice
func VolumeMountFactory(volumes map[string]string) []corev1.VolumeMount {
	var volumeMounts []corev1.VolumeMount
	for name, path := range volumes {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      name,
			MountPath: path,
		})
	}
	return volumeMounts
}
