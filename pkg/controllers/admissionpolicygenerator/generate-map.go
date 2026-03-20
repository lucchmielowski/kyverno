package admissionpolicygenerator

import (
	"context"

	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/admissionpolicy"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
)

func (c *controller) handleMAPGeneration(ctx context.Context, mpol *policiesv1beta1.MutatingPolicy) error {
	if c.mapWriter == nil {
		return nil
	}
	genericPolicy := engineapi.NewMutatingPolicy(mpol)
	if !admissionpolicy.HasMutatingAdmissionPolicyPermission(c.checker) {
		logger.V(2).Info("insufficient permissions to generate MutatingAdmissionPolicies")
		c.updatePolicyStatus(ctx, genericPolicy, false, "insufficient permissions to generate MutatingAdmissionPolicies")
		return nil
	}
	if !admissionpolicy.HasMutatingAdmissionPolicyBindingPermission(c.checker) {
		logger.V(2).Info("insufficient permissions to generate MutatingAdmissionPolicyBindings")
		c.updatePolicyStatus(ctx, genericPolicy, false, "insufficient permissions to generate MutatingAdmissionPolicyBindings")
		return nil
	}
	return c.mapWriter.Reconcile(ctx, c, mpol, genericPolicy)
}
