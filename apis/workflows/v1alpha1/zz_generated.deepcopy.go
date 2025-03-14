//go:build !ignore_autogenerated

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	k8sv1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parent) DeepCopyInto(out *Parent) {
	*out = *in
	if in.ProjectRef != nil {
		in, out := &in.ProjectRef, &out.ProjectRef
		*out = new(v1beta1.ProjectRef)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parent.
func (in *Parent) DeepCopy() *Parent {
	if in == nil {
		return nil
	}
	out := new(Parent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflow) DeepCopyInto(out *WorkflowsWorkflow) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflow.
func (in *WorkflowsWorkflow) DeepCopy() *WorkflowsWorkflow {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkflowsWorkflow) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowIdentity) DeepCopyInto(out *WorkflowsWorkflowIdentity) {
	*out = *in
	if in.parent != nil {
		in, out := &in.parent, &out.parent
		*out = new(WorkflowsWorkflowParent)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowIdentity.
func (in *WorkflowsWorkflowIdentity) DeepCopy() *WorkflowsWorkflowIdentity {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowList) DeepCopyInto(out *WorkflowsWorkflowList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WorkflowsWorkflow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowList.
func (in *WorkflowsWorkflowList) DeepCopy() *WorkflowsWorkflowList {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkflowsWorkflowList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowObservedState) DeepCopyInto(out *WorkflowsWorkflowObservedState) {
	*out = *in
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
	if in.RevisionId != nil {
		in, out := &in.RevisionId, &out.RevisionId
		*out = new(string)
		**out = **in
	}
	if in.CreateTime != nil {
		in, out := &in.CreateTime, &out.CreateTime
		*out = new(string)
		**out = **in
	}
	if in.UpdateTime != nil {
		in, out := &in.UpdateTime, &out.UpdateTime
		*out = new(string)
		**out = **in
	}
	if in.RevisionCreateTime != nil {
		in, out := &in.RevisionCreateTime, &out.RevisionCreateTime
		*out = new(string)
		**out = **in
	}
	if in.StateError != nil {
		in, out := &in.StateError, &out.StateError
		*out = new(WorkflowsWorkflow_StateError)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowObservedState.
func (in *WorkflowsWorkflowObservedState) DeepCopy() *WorkflowsWorkflowObservedState {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowObservedState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowParent) DeepCopyInto(out *WorkflowsWorkflowParent) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowParent.
func (in *WorkflowsWorkflowParent) DeepCopy() *WorkflowsWorkflowParent {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowParent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowRef) DeepCopyInto(out *WorkflowsWorkflowRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowRef.
func (in *WorkflowsWorkflowRef) DeepCopy() *WorkflowsWorkflowRef {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowSpec) DeepCopyInto(out *WorkflowsWorkflowSpec) {
	*out = *in
	in.Parent.DeepCopyInto(&out.Parent)
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ServiceAccountRef != nil {
		in, out := &in.ServiceAccountRef, &out.ServiceAccountRef
		*out = new(v1beta1.IAMServiceAccountRef)
		**out = **in
	}
	if in.SourceContents != nil {
		in, out := &in.SourceContents, &out.SourceContents
		*out = new(string)
		**out = **in
	}
	if in.KMSCryptoKeyRef != nil {
		in, out := &in.KMSCryptoKeyRef, &out.KMSCryptoKeyRef
		*out = new(v1beta1.KMSCryptoKeyRef)
		**out = **in
	}
	if in.CallLogLevel != nil {
		in, out := &in.CallLogLevel, &out.CallLogLevel
		*out = new(string)
		**out = **in
	}
	if in.UserEnvVars != nil {
		in, out := &in.UserEnvVars, &out.UserEnvVars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowSpec.
func (in *WorkflowsWorkflowSpec) DeepCopy() *WorkflowsWorkflowSpec {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflowStatus) DeepCopyInto(out *WorkflowsWorkflowStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]k8sv1alpha1.Condition, len(*in))
		copy(*out, *in)
	}
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.ExternalRef != nil {
		in, out := &in.ExternalRef, &out.ExternalRef
		*out = new(string)
		**out = **in
	}
	if in.ObservedState != nil {
		in, out := &in.ObservedState, &out.ObservedState
		*out = new(WorkflowsWorkflowObservedState)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflowStatus.
func (in *WorkflowsWorkflowStatus) DeepCopy() *WorkflowsWorkflowStatus {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflowStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkflowsWorkflow_StateError) DeepCopyInto(out *WorkflowsWorkflow_StateError) {
	*out = *in
	if in.Details != nil {
		in, out := &in.Details, &out.Details
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkflowsWorkflow_StateError.
func (in *WorkflowsWorkflow_StateError) DeepCopy() *WorkflowsWorkflow_StateError {
	if in == nil {
		return nil
	}
	out := new(WorkflowsWorkflow_StateError)
	in.DeepCopyInto(out)
	return out
}
