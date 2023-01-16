/*
Copyright (c) 2020 - 2023 SUSE LLC

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

// Code generated by main. DO NOT EDIT.

// +k8s:deepcopy-gen=package
// +groupName=fleet.cattle.io
package v1alpha1

import (
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	BundleResourceName                   = "bundles"
	BundleDeploymentResourceName         = "bundledeployments"
	BundleNamespaceMappingResourceName   = "bundlenamespacemappings"
	ClusterResourceName                  = "clusters"
	ClusterGroupResourceName             = "clustergroups"
	ClusterRegistrationResourceName      = "clusterregistrations"
	ClusterRegistrationTokenResourceName = "clusterregistrationtokens"
	ContentResourceName                  = "contents"
	GitRepoResourceName                  = "gitrepos"
	GitRepoRestrictionResourceName       = "gitreporestrictions"
	ImageScanResourceName                = "imagescans"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: fleet.GroupName, Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Bundle{},
		&BundleList{},
		&BundleDeployment{},
		&BundleDeploymentList{},
		&BundleNamespaceMapping{},
		&BundleNamespaceMappingList{},
		&Cluster{},
		&ClusterList{},
		&ClusterGroup{},
		&ClusterGroupList{},
		&ClusterRegistration{},
		&ClusterRegistrationList{},
		&ClusterRegistrationToken{},
		&ClusterRegistrationTokenList{},
		&Content{},
		&ContentList{},
		&GitRepo{},
		&GitRepoList{},
		&GitRepoRestriction{},
		&GitRepoRestrictionList{},
		&ImageScan{},
		&ImageScanList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
