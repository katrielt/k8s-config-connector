// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Config Connector and manual
//     changes will be clobbered when the file is regenerated.
//
// ----------------------------------------------------------------------------

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

package v1alpha1

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ManagedsslcertificateManaged struct {
	/* Immutable. Domains for which a managed SSL certificate will be valid.  Currently,
	there can be up to 100 domains in this list. */
	Domains []string `json:"domains"`
}

type ComputeManagedSSLCertificateSpec struct {
	/* The unique identifier for the resource. */
	// +optional
	CertificateId *int `json:"certificateId,omitempty"`

	/* Immutable. An optional description of this resource. */
	// +optional
	Description *string `json:"description,omitempty"`

	/* Immutable. Properties relevant to a managed certificate.  These will be used if the
	certificate is managed (as indicated by a value of 'MANAGED' in 'type'). */
	// +optional
	Managed *ManagedsslcertificateManaged `json:"managed,omitempty"`

	/* The project that this resource belongs to. */
	ProjectRef v1alpha1.ResourceRef `json:"projectRef"`

	/* Immutable. Optional. The name of the resource. Used for creation and acquisition. When unset, the value of `metadata.name` is used as the default. */
	// +optional
	ResourceID *string `json:"resourceID,omitempty"`

	/* Immutable. Enum field whose value is always 'MANAGED' - used to signal to the API
	which type this is. Default value: "MANAGED" Possible values: ["MANAGED"]. */
	// +optional
	Type *string `json:"type,omitempty"`
}

type ComputeManagedSSLCertificateStatus struct {
	/* Conditions represent the latest available observations of the
	   ComputeManagedSSLCertificate's current state. */
	Conditions []v1alpha1.Condition `json:"conditions,omitempty"`
	/* Creation timestamp in RFC3339 text format. */
	// +optional
	CreationTimestamp *string `json:"creationTimestamp,omitempty"`

	/* Expire time of the certificate in RFC3339 text format. */
	// +optional
	ExpireTime *string `json:"expireTime,omitempty"`

	/* ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller. If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource. */
	// +optional
	ObservedGeneration *int `json:"observedGeneration,omitempty"`

	// +optional
	SelfLink *string `json:"selfLink,omitempty"`

	/* Domains associated with the certificate via Subject Alternative Name. */
	// +optional
	SubjectAlternativeNames []string `json:"subjectAlternativeNames,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ComputeManagedSSLCertificate is the Schema for the compute API
// +k8s:openapi-gen=true
type ComputeManagedSSLCertificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComputeManagedSSLCertificateSpec   `json:"spec,omitempty"`
	Status ComputeManagedSSLCertificateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ComputeManagedSSLCertificateList contains a list of ComputeManagedSSLCertificate
type ComputeManagedSSLCertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComputeManagedSSLCertificate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComputeManagedSSLCertificate{}, &ComputeManagedSSLCertificateList{})
}