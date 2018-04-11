/*
Copyright 2018 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	fedv1a1 "github.com/marun/federation-v2/pkg/apis/federation/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VersionCompareType int

const (
	ResourceVersion VersionCompareType = iota
	Generation
)

type ComparisonHelper interface {
	GetVersion(objectMeta *metav1.ObjectMeta) string
	Equivalent(objMeta1, objectMeta2 *metav1.ObjectMeta) bool
}

// NewComparisonHelper instantiates and returns a Resource or Generation Helper
// struct that implements the ComparisonHelper interface based on the version
// comparison type passed in.
func NewComparisonHelper(compareType VersionCompareType) (ComparisonHelper, error) {
	switch compareType {
	case ResourceVersion:
		return &ResourceHelper{}, nil
	case Generation:
		return &GenerationHelper{}, nil
	default:
		return nil, fmt.Errorf("Unrecognized version comparison type %v", compareType)
	}
}

type GenerationHelper struct{}

// GetVersion returns a string containing the version in the resource's
// ObjectMeta using the resource comparison type to perform for that
// resource.
func (GenerationHelper) GetVersion(objectMeta *metav1.ObjectMeta) string {
	return strconv.FormatInt(objectMeta.Generation, 10)
}

// Equivalent returns true if both object metas passed in are equivalent, false
// otherwise.
func (GenerationHelper) Equivalent(obj1Meta, obj2Meta *metav1.ObjectMeta) bool {
	return ObjectMetaEquivalent(*obj1Meta, *obj2Meta)
}

type ResourceHelper struct{}

// GetVersion returns a string containing the version in the resource's
// ObjectMeta using the resource comparison type to perform for that
// resource.
func (ResourceHelper) GetVersion(objectMeta *metav1.ObjectMeta) string {
	return objectMeta.ResourceVersion
}

// Equivalent returns true for ResourceVersion comparison as it doesn't require
// comparing ObjectMeta.
func (ResourceHelper) Equivalent(obj1Meta, obj2Meta *metav1.ObjectMeta) bool {
	return true
}

// SortClusterVersions ASCII sorts the given cluster versions slice
// based on cluster name.
func SortClusterVersions(versions []fedv1a1.ClusterObjectVersion) {
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].ClusterName < versions[j].ClusterName
	})
}

// PropagatedVersionStatusEquivalent returns true if both statuses are equal by
// comparing Template and Override version, and their ClusterVersion slices;
// false otherwise.
func PropagatedVersionStatusEquivalent(pvs1, pvs2 *fedv1a1.PropagatedVersionStatus) bool {
	return pvs1.TemplateVersion == pvs2.TemplateVersion &&
		pvs1.OverrideVersion == pvs2.OverrideVersion &&
		reflect.DeepEqual(pvs1.ClusterVersions, pvs2.ClusterVersions)
}
