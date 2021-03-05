// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	v1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Artifacts) DeepCopyInto(out *Artifacts) {
	*out = *in
	if in.Images != nil {
		in, out := &in.Images, &out.Images
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Objects.DeepCopyInto(&out.Objects)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Artifacts.
func (in *Artifacts) DeepCopy() *Artifacts {
	if in == nil {
		return nil
	}
	out := new(Artifacts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Build) DeepCopyInto(out *Build) {
	*out = *in
	in.Artifacts.DeepCopyInto(&out.Artifacts)
	if in.Images != nil {
		in, out := &in.Images, &out.Images
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.LogsBucketRef = in.LogsBucketRef
	in.Options.DeepCopyInto(&out.Options)
	if in.Secret != nil {
		in, out := &in.Secret, &out.Secret
		*out = make([]Secret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Source.DeepCopyInto(&out.Source)
	if in.Step != nil {
		in, out := &in.Step, &out.Step
		*out = make([]Step, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Substitutions != nil {
		in, out := &in.Substitutions, &out.Substitutions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Build.
func (in *Build) DeepCopy() *Build {
	if in == nil {
		return nil
	}
	out := new(Build)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudBuildTrigger) DeepCopyInto(out *CloudBuildTrigger) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudBuildTrigger.
func (in *CloudBuildTrigger) DeepCopy() *CloudBuildTrigger {
	if in == nil {
		return nil
	}
	out := new(CloudBuildTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudBuildTrigger) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudBuildTriggerList) DeepCopyInto(out *CloudBuildTriggerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CloudBuildTrigger, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudBuildTriggerList.
func (in *CloudBuildTriggerList) DeepCopy() *CloudBuildTriggerList {
	if in == nil {
		return nil
	}
	out := new(CloudBuildTriggerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CloudBuildTriggerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudBuildTriggerSpec) DeepCopyInto(out *CloudBuildTriggerSpec) {
	*out = *in
	in.Build.DeepCopyInto(&out.Build)
	out.Github = in.Github
	if in.IgnoredFiles != nil {
		in, out := &in.IgnoredFiles, &out.IgnoredFiles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.IncludedFiles != nil {
		in, out := &in.IncludedFiles, &out.IncludedFiles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Substitutions != nil {
		in, out := &in.Substitutions, &out.Substitutions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.TriggerTemplate = in.TriggerTemplate
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudBuildTriggerSpec.
func (in *CloudBuildTriggerSpec) DeepCopy() *CloudBuildTriggerSpec {
	if in == nil {
		return nil
	}
	out := new(CloudBuildTriggerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudBuildTriggerStatus) DeepCopyInto(out *CloudBuildTriggerStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1alpha1.Condition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudBuildTriggerStatus.
func (in *CloudBuildTriggerStatus) DeepCopy() *CloudBuildTriggerStatus {
	if in == nil {
		return nil
	}
	out := new(CloudBuildTriggerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Github) DeepCopyInto(out *Github) {
	*out = *in
	out.PullRequest = in.PullRequest
	out.Push = in.Push
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Github.
func (in *Github) DeepCopy() *Github {
	if in == nil {
		return nil
	}
	out := new(Github)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Objects) DeepCopyInto(out *Objects) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Timing != nil {
		in, out := &in.Timing, &out.Timing
		*out = make([]Timing, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Objects.
func (in *Objects) DeepCopy() *Objects {
	if in == nil {
		return nil
	}
	out := new(Objects)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Options) DeepCopyInto(out *Options) {
	*out = *in
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SecretEnv != nil {
		in, out := &in.SecretEnv, &out.SecretEnv
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SourceProvenanceHash != nil {
		in, out := &in.SourceProvenanceHash, &out.SourceProvenanceHash
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Volumes, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Options.
func (in *Options) DeepCopy() *Options {
	if in == nil {
		return nil
	}
	out := new(Options)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PullRequest) DeepCopyInto(out *PullRequest) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PullRequest.
func (in *PullRequest) DeepCopy() *PullRequest {
	if in == nil {
		return nil
	}
	out := new(PullRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Push) DeepCopyInto(out *Push) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Push.
func (in *Push) DeepCopy() *Push {
	if in == nil {
		return nil
	}
	out := new(Push)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepoSource) DeepCopyInto(out *RepoSource) {
	*out = *in
	out.RepoRef = in.RepoRef
	if in.Substitutions != nil {
		in, out := &in.Substitutions, &out.Substitutions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepoSource.
func (in *RepoSource) DeepCopy() *RepoSource {
	if in == nil {
		return nil
	}
	out := new(RepoSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Secret) DeepCopyInto(out *Secret) {
	*out = *in
	out.KmsKeyRef = in.KmsKeyRef
	if in.SecretEnv != nil {
		in, out := &in.SecretEnv, &out.SecretEnv
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Secret.
func (in *Secret) DeepCopy() *Secret {
	if in == nil {
		return nil
	}
	out := new(Secret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Source) DeepCopyInto(out *Source) {
	*out = *in
	in.RepoSource.DeepCopyInto(&out.RepoSource)
	out.StorageSource = in.StorageSource
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Source.
func (in *Source) DeepCopy() *Source {
	if in == nil {
		return nil
	}
	out := new(Source)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Step) DeepCopyInto(out *Step) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SecretEnv != nil {
		in, out := &in.SecretEnv, &out.SecretEnv
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Volumes, len(*in))
		copy(*out, *in)
	}
	if in.WaitFor != nil {
		in, out := &in.WaitFor, &out.WaitFor
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageSource) DeepCopyInto(out *StorageSource) {
	*out = *in
	out.BucketRef = in.BucketRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageSource.
func (in *StorageSource) DeepCopy() *StorageSource {
	if in == nil {
		return nil
	}
	out := new(StorageSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Timing) DeepCopyInto(out *Timing) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Timing.
func (in *Timing) DeepCopy() *Timing {
	if in == nil {
		return nil
	}
	out := new(Timing)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TriggerTemplate) DeepCopyInto(out *TriggerTemplate) {
	*out = *in
	out.RepoRef = in.RepoRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TriggerTemplate.
func (in *TriggerTemplate) DeepCopy() *TriggerTemplate {
	if in == nil {
		return nil
	}
	out := new(TriggerTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Volumes) DeepCopyInto(out *Volumes) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Volumes.
func (in *Volumes) DeepCopy() *Volumes {
	if in == nil {
		return nil
	}
	out := new(Volumes)
	in.DeepCopyInto(out)
	return out
}