//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProducerConfig) DeepCopyInto(out *ProducerConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerEgressAdapterConfig.
func (in *ProducerConfig) DeepCopy() *ProducerConfig {
	if in == nil {
		return nil
	}
	out := new(ProducerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsumerConfig) DeepCopyInto(out *ConsumerConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BrokerIngressAdapterConfig.
func (in *ConsumerConfig) DeepCopy() *ConsumerConfig {
	if in == nil {
		return nil
	}
	out := new(ConsumerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecondaryAdapterConfig) DeepCopyInto(out *SecondaryAdapterConfig) {
	*out = *in
	if in.InvocationConfig != nil {
		in, out := &in.InvocationConfig, &out.InvocationConfig
		*out = new(InvocationConfig)
		**out = **in
	}
	if in.RepositoryConfig != nil {
		in, out := &in.RepositoryConfig, &out.RepositoryConfig
		*out = new(RepositoryClientConfig)
		**out = **in
	}
	if in.ProducerConfig != nil {
		in, out := &in.ProducerConfig, &out.ProducerConfig
		*out = new(ProducerConfig)
		**out = **in
	}
	if in.Id != nil {
		in, out := &in.Id, &out.Id
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressAdapterConfig.
func (in *SecondaryAdapterConfig) DeepCopy() *SecondaryAdapterConfig {
	if in == nil {
		return nil
	}
	out := new(SecondaryAdapterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressAdapter) DeepCopyInto(out *IngressAdapter) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressAdapter.
func (in *IngressAdapter) DeepCopy() *IngressAdapter {
	if in == nil {
		return nil
	}
	out := new(IngressAdapter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IngressAdapter) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressAdapterList) DeepCopyInto(out *IngressAdapterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IngressAdapter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressAdapterList.
func (in *IngressAdapterList) DeepCopy() *IngressAdapterList {
	if in == nil {
		return nil
	}
	out := new(IngressAdapterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IngressAdapterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrimaryAdapterSpec) DeepCopyInto(out *PrimaryAdapterSpec) {
	*out = *in
	if in.ServerConfig != nil {
		in, out := &in.ServerConfig, &out.ServerConfig
		*out = new(ServerConfig)
		**out = **in
	}
	if in.ConsumerConfig != nil {
		in, out := &in.ConsumerConfig, &out.ConsumerConfig
		*out = new(ConsumerConfig)
		**out = **in
	}
	if in.RepositoryConfig != nil {
		in, out := &in.RepositoryConfig, &out.RepositoryConfig
		*out = new(RepositoryConfig)
		**out = **in
	}
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]Step, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressAdapterSpec.
func (in *PrimaryAdapterSpec) DeepCopy() *PrimaryAdapterSpec {
	if in == nil {
		return nil
	}
	out := new(PrimaryAdapterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressAdapterStatus) DeepCopyInto(out *IngressAdapterStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressAdapterStatus.
func (in *IngressAdapterStatus) DeepCopy() *IngressAdapterStatus {
	if in == nil {
		return nil
	}
	out := new(IngressAdapterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalAdapterConfig) DeepCopyInto(out *InternalAdapterConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalAdapterConfig.
func (in *InternalAdapterConfig) DeepCopy() *InternalAdapterConfig {
	if in == nil {
		return nil
	}
	out := new(InternalAdapterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceUnit) DeepCopyInto(out *ServiceUnit) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceUnit.
func (in *ServiceUnit) DeepCopy() *ServiceUnit {
	if in == nil {
		return nil
	}
	out := new(ServiceUnit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceUnit) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceUnitConfig) DeepCopyInto(out *ServiceUnitConfig) {
	*out = *in
	if in.AdapterConfigs != nil {
		in, out := &in.AdapterConfigs, &out.AdapterConfigs
		*out = make([]PrimaryAdapterSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceUnitConfig.
func (in *ServiceUnitConfig) DeepCopy() *ServiceUnitConfig {
	if in == nil {
		return nil
	}
	out := new(ServiceUnitConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceUnitList) DeepCopyInto(out *ServiceUnitList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceUnit, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceUnitList.
func (in *ServiceUnitList) DeepCopy() *ServiceUnitList {
	if in == nil {
		return nil
	}
	out := new(ServiceUnitList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceUnitList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceUnitSpec) DeepCopyInto(out *ServiceUnitSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceUnitSpec.
func (in *ServiceUnitSpec) DeepCopy() *ServiceUnitSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceUnitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceUnitStatus) DeepCopyInto(out *ServiceUnitStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceUnitStatus.
func (in *ServiceUnitStatus) DeepCopy() *ServiceUnitStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceUnitStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryClientConfig) DeepCopyInto(out *RepositoryClientConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulEgressAdapterConfig.
func (in *RepositoryClientConfig) DeepCopy() *RepositoryClientConfig {
	if in == nil {
		return nil
	}
	out := new(RepositoryClientConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryConfig) DeepCopyInto(out *RepositoryConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulIngressAdapterConfig.
func (in *RepositoryConfig) DeepCopy() *RepositoryConfig {
	if in == nil {
		return nil
	}
	out := new(RepositoryConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InvocationConfig) DeepCopyInto(out *InvocationConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatelessEgressAdapterConfig.
func (in *InvocationConfig) DeepCopy() *InvocationConfig {
	if in == nil {
		return nil
	}
	out := new(InvocationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerConfig) DeepCopyInto(out *ServerConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatelessIngressAdapterConfig.
func (in *ServerConfig) DeepCopy() *ServerConfig {
	if in == nil {
		return nil
	}
	out := new(ServerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Step) DeepCopyInto(out *Step) {
	*out = *in
	if in.AdapterConfig != nil {
		in, out := &in.AdapterConfig, &out.AdapterConfig
		*out = new(SecondaryAdapterConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Step.
func (in *Step) DeepCopy() *Step {
	if in == nil {
		return nil
	}
	out := new(Step)
	in.DeepCopyInto(out)
	return out
}
