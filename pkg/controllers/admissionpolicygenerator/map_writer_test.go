package admissionpolicygenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestNewMutatingAdmissionPolicyWriter(t *testing.T) {
	t.Run("v1alpha1", func(t *testing.T) {
		w, err := NewMutatingAdmissionPolicyWriter(schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v1alpha1"})
		require.NoError(t, err)
		assert.IsType(t, mapWriterAlpha{}, w)
	})
	t.Run("v1beta1", func(t *testing.T) {
		w, err := NewMutatingAdmissionPolicyWriter(schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v1beta1"})
		require.NoError(t, err)
		assert.IsType(t, mapWriterBeta{}, w)
	})
	t.Run("v1 unsupported", func(t *testing.T) {
		_, err := NewMutatingAdmissionPolicyWriter(schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v1"})
		require.Error(t, err)
	})
	t.Run("unknown version", func(t *testing.T) {
		_, err := NewMutatingAdmissionPolicyWriter(schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v2"})
		require.Error(t, err)
	})
}
