package admissionpolicygenerator

import (
	"context"

	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/admissionpolicy"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	admissionregistrationv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mapWriterAlpha struct{}

func (mapWriterAlpha) Reconcile(ctx context.Context, c *controller, mpol *policiesv1beta1.MutatingPolicy, genericPolicy engineapi.GenericPolicy) error {
	return reconcileMutatingAdmissionPolicies(ctx, c, mpol, genericPolicy,
		c.client.AdmissionregistrationV1alpha1().MutatingAdmissionPolicies(),
		c.client.AdmissionregistrationV1alpha1().MutatingAdmissionPolicyBindings(),
		func(name string) *admissionregistrationv1alpha1.MutatingAdmissionPolicy {
			return &admissionregistrationv1alpha1.MutatingAdmissionPolicy{ObjectMeta: metav1.ObjectMeta{Name: name}}
		},
		func(name string) *admissionregistrationv1alpha1.MutatingAdmissionPolicyBinding {
			return &admissionregistrationv1alpha1.MutatingAdmissionPolicyBinding{ObjectMeta: metav1.ObjectMeta{Name: name}}
		},
		admissionpolicy.BuildMutatingAdmissionPolicy,
		admissionpolicy.BuildMutatingAdmissionPolicyBinding,
	)
}
