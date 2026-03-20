package admissionpolicygenerator

import (
	"context"

	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/admissionpolicy"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mapWriterBeta struct{}

func (mapWriterBeta) Reconcile(ctx context.Context, c *controller, mpol *policiesv1beta1.MutatingPolicy, genericPolicy engineapi.GenericPolicy) error {
	return reconcileMutatingAdmissionPolicies(ctx, c, mpol, genericPolicy,
		c.client.AdmissionregistrationV1beta1().MutatingAdmissionPolicies(),
		c.client.AdmissionregistrationV1beta1().MutatingAdmissionPolicyBindings(),
		func(name string) *admissionregistrationv1beta1.MutatingAdmissionPolicy {
			return &admissionregistrationv1beta1.MutatingAdmissionPolicy{ObjectMeta: metav1.ObjectMeta{Name: name}}
		},
		func(name string) *admissionregistrationv1beta1.MutatingAdmissionPolicyBinding {
			return &admissionregistrationv1beta1.MutatingAdmissionPolicyBinding{ObjectMeta: metav1.ObjectMeta{Name: name}}
		},
		admissionpolicy.BuildMutatingAdmissionPolicyBeta,
		admissionpolicy.BuildMutatingAdmissionPolicyBindingBeta,
	)
}
