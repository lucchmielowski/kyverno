package admissionpolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDiscoverMutatingAdmissionPolicyVersion(t *testing.T) {
	t.Run("selects v1beta1 when v1 has no MAP and both beta and alpha list MAP", func(t *testing.T) {
		client := fake.NewClientset()
		client.Fake.Resources = []*metav1.APIResourceList{
			{
				GroupVersion: "admissionregistration.k8s.io/v1",
				APIResources: []metav1.APIResource{{Name: "validatingadmissionpolicies"}},
			},
			{
				GroupVersion: "admissionregistration.k8s.io/v1beta1",
				APIResources: []metav1.APIResource{{Name: "mutatingadmissionpolicies"}},
			},
			{
				GroupVersion: "admissionregistration.k8s.io/v1alpha1",
				APIResources: []metav1.APIResource{{Name: "mutatingadmissionpolicies"}},
			},
		}
		gv, found, err := DiscoverMutatingAdmissionPolicyVersion(client.Discovery())
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, schema.GroupVersion{Group: admissionregistrationGroup, Version: "v1beta1"}, gv)
	})

	t.Run("selects v1beta1 when v1 lists MAP but Kyverno only supports beta and alpha", func(t *testing.T) {
		client := fake.NewClientset()
		client.Fake.Resources = []*metav1.APIResourceList{
			{
				GroupVersion: "admissionregistration.k8s.io/v1",
				APIResources: []metav1.APIResource{
					{Name: "validatingadmissionpolicies"},
					{Name: "mutatingadmissionpolicies"},
				},
			},
			{
				GroupVersion: "admissionregistration.k8s.io/v1beta1",
				APIResources: []metav1.APIResource{{Name: "mutatingadmissionpolicies"}},
			},
			{
				GroupVersion: "admissionregistration.k8s.io/v1alpha1",
				APIResources: []metav1.APIResource{{Name: "mutatingadmissionpolicies"}},
			},
		}
		gv, found, err := DiscoverMutatingAdmissionPolicyVersion(client.Discovery())
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, schema.GroupVersion{Group: admissionregistrationGroup, Version: "v1beta1"}, gv)
	})

	t.Run("selects v1alpha1 when only alpha serves MAP", func(t *testing.T) {
		client := fake.NewClientset()
		client.Fake.Resources = []*metav1.APIResourceList{
			{
				GroupVersion: "admissionregistration.k8s.io/v1alpha1",
				APIResources: []metav1.APIResource{{Name: "mutatingadmissionpolicies"}},
			},
		}
		gv, found, err := DiscoverMutatingAdmissionPolicyVersion(client.Discovery())
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, schema.GroupVersion{Group: admissionregistrationGroup, Version: "v1alpha1"}, gv)
	})

	t.Run("not found when MAP absent", func(t *testing.T) {
		client := fake.NewClientset()
		client.Fake.Resources = []*metav1.APIResourceList{}
		_, found, err := DiscoverMutatingAdmissionPolicyVersion(client.Discovery())
		if err != nil {
			assert.False(t, found)
			return
		}
		assert.False(t, found)
	})
}

func TestSupportedMutatingAdmissionPolicyGroupVersions_order(t *testing.T) {
	vers := SupportedMutatingAdmissionPolicyGroupVersions()
	require.GreaterOrEqual(t, len(vers), 2)
	assert.Equal(t, "v1beta1", vers[0].Version)
	assert.Equal(t, "v1alpha1", vers[1].Version)
	for _, gv := range vers {
		assert.Equal(t, admissionregistrationGroup, gv.Group)
	}
}
