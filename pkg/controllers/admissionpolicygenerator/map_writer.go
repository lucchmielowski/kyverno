package admissionpolicygenerator

import (
	"context"
	"fmt"

	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// MutatingAdmissionPolicyWriter reconciles a Kyverno MutatingPolicy to Kubernetes MutatingAdmissionPolicy resources.
type MutatingAdmissionPolicyWriter interface {
	Reconcile(ctx context.Context, c *controller, mpol *policiesv1beta1.MutatingPolicy, genericPolicy engineapi.GenericPolicy) error
}

// NewMutatingAdmissionPolicyWriter returns a writer for the discovered admissionregistration.k8s.io MAP API version.
func NewMutatingAdmissionPolicyWriter(gv schema.GroupVersion) (MutatingAdmissionPolicyWriter, error) {
	switch gv.Version {
	case "v1alpha1":
		return mapWriterAlpha{}, nil
	case "v1beta1":
		return mapWriterBeta{}, nil
	case "v1":
		return nil, fmt.Errorf("MutatingAdmissionPolicy admissionregistration.k8s.io/v1 is not supported in this Kyverno version yet")
	default:
		return nil, fmt.Errorf("unsupported MutatingAdmissionPolicy API version %q", gv.Version)
	}
}
