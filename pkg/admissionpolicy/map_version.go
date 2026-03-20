package admissionpolicy

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

const admissionregistrationGroup = "admissionregistration.k8s.io"

// SupportedMutatingAdmissionPolicyGroupVersions returns admissionregistration.k8s.io API versions
// that this Kyverno build can use for MutatingAdmissionPolicy, in preference order (newest supported first).
//
// TODO: Only versions with a working client/informer/writer are listed. When MAP graduates to v1 and Kyverno
// implements it, add {Group: admissionregistrationGroup, Version: "v1"} at the front of the slice.
func SupportedMutatingAdmissionPolicyGroupVersions() []schema.GroupVersion {
	return []schema.GroupVersion{
		{Group: admissionregistrationGroup, Version: "v1beta1"},
		{Group: admissionregistrationGroup, Version: "v1alpha1"},
	}
}

// DiscoverMutatingAdmissionPolicyVersion returns the first GroupVersion from
// SupportedMutatingAdmissionPolicyGroupVersions that the apiserver serves and that lists the
// mutatingadmissionpolicies resource. If none match, found is false; err is set only when every
// probe failed with a discovery error (e.g. complete discovery failure).
func DiscoverMutatingAdmissionPolicyVersion(sri discovery.ServerResourcesInterface) (gv schema.GroupVersion, found bool, err error) {
	var lastErr error
	for _, cand := range SupportedMutatingAdmissionPolicyGroupVersions() {
		list, e := sri.ServerResourcesForGroupVersion(cand.String())
		if e != nil {
			lastErr = e
			continue
		}
		for i := range list.APIResources {
			if list.APIResources[i].Name == "mutatingadmissionpolicies" {
				return cand, true, nil
			}
		}
	}
	if lastErr != nil {
		return schema.GroupVersion{}, false, lastErr
	}
	return schema.GroupVersion{}, false, nil
}
