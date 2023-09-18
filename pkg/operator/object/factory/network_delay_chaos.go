package factory

import "github.com/hanapedia/the-bench/pkg/operator/object/crd"

type NetworkChaosArgs struct {
	Name            string
	Namespace       string
	TargetNamespace string
	Selector        map[string]string
	Duration        string
	Latency         string
	Jitter          string
}

// NetworkChaosFactory create type network chaos kubernetes custom resource objects for chaos-mesh.
func NetworkChaosFactory(args *NetworkChaosArgs) crd.NetworkChaos {
	return crd.NetworkChaos{
		TypeMeta:   TypeMetaFactory("NetworkChaos", "chaos-mesh.org/v1alpha1"),
		ObjectMeta: ObjectMetaFactory(ObjectMetaOptions{Name: args.Name, Namespace: args.Namespace}),
		Spec:       NetworkChaosSpecFactory(args),
	}
}

// NetworkChaosSpecFactory create type namespace kubernetes objects.
func NetworkChaosSpecFactory(args *NetworkChaosArgs) crd.NetworkChaosSpec {
	return crd.NetworkChaosSpec{
		PodSelector: crd.PodSelector{
			Selector: crd.PodSelectorSpec{
				GenericSelectorSpec: crd.GenericSelectorSpec{
					Namespaces:     []string{args.TargetNamespace},
					LabelSelectors: args.Selector,
				},
			},
			Mode: crd.AllMode,
		},
		Duration:  &args.Duration,
		Action:    crd.DelayAction,
		Direction: crd.Both,
		TcParameter: crd.TcParameter{
			Delay: &crd.DelaySpec{
				Latency: args.Latency,
				Jitter:  args.Jitter,
			},
		},
		Target: &crd.PodSelector{
			Selector: crd.PodSelectorSpec{
				GenericSelectorSpec: crd.GenericSelectorSpec{
					Namespaces: []string{args.TargetNamespace},
				},
			},
			Mode: crd.AllMode,
		},
	}

}
