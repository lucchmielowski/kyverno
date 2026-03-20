package admissionpolicygenerator

import (
	"context"
	"fmt"

	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	controllerutils "github.com/kyverno/kyverno/pkg/utils/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mapAPIObject is a Kubernetes API object usable with controllerutils.Update (metadata + DeepCopy).
type mapAPIObject[T any] interface {
	metav1.Object
	DeepCopy() T
}

// reconcileMutatingAdmissionPolicies implements MAP + MAPBinding sync for a specific admissionregistration API version.
func reconcileMutatingAdmissionPolicies[TPolicy mapAPIObject[TPolicy], TBinding mapAPIObject[TBinding]](
	ctx context.Context,
	c *controller,
	mpol *policiesv1beta1.MutatingPolicy,
	genericPolicy engineapi.GenericPolicy,
	policyClient controllerutils.ObjectClient[TPolicy],
	bindingClient controllerutils.ObjectClient[TBinding],
	newPolicy func(name string) TPolicy,
	newBinding func(name string) TBinding,
	buildPolicy func(TPolicy, *policiesv1beta1.MutatingPolicy, []policiesv1beta1.PolicyException),
	buildBinding func(TBinding, *policiesv1beta1.MutatingPolicy),
) error {
	mapName := "mpol-" + mpol.GetName()
	mapBindingName := constructBindingName(mapName)

	observedMAP, mapErr := policyClient.Get(ctx, mapName, metav1.GetOptions{})
	var mapNotFound bool
	if mapErr != nil {
		if !apierrors.IsNotFound(mapErr) {
			return fmt.Errorf("failed to get mutatingadmissionpolicy %s: %w", mapName, mapErr)
		}
		mapNotFound = true
	}

	observedMAPbinding, mapBindingErr := bindingClient.Get(ctx, mapBindingName, metav1.GetOptions{})
	var mapBindingNotFound bool
	if mapBindingErr != nil {
		if !apierrors.IsNotFound(mapBindingErr) {
			return fmt.Errorf("failed to get mutatingadmissionpolicybinding %s: %w", mapBindingName, mapBindingErr)
		}
		mapBindingNotFound = true
	}

	wantMap := mpol.GetSpec().GenerateMutatingAdmissionPolicyEnabled()
	shouldDelete := !wantMap
	var reason string
	if wantMap {
		if len(mpol.GetStatus().Autogen.Configs) > 0 {
			shouldDelete = true
			reason = "skip generating MutatingAdmissionPolicy: pod controllers autogen is enabled."
		}
	} else {
		reason = "skip generating MutatingAdmissionPolicy: not enabled."
	}

	if shouldDelete {
		if !mapNotFound {
			if err := policyClient.Delete(ctx, mapName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
		if !mapBindingNotFound {
			if err := bindingClient.Delete(ctx, mapBindingName, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
		c.updatePolicyStatus(ctx, genericPolicy, false, reason)
		return nil
	}

	celexceptions, err := c.getCELExceptions(mpol.GetName())
	if err != nil {
		return fmt.Errorf("failed to get celexceptions by name %s: %w", mpol.GetName(), err)
	}

	if mapNotFound {
		observedMAP = newPolicy(mapName)
	}
	if mapBindingNotFound {
		observedMAPbinding = newBinding(mapBindingName)
	}

	if observedMAP.GetResourceVersion() == "" {
		buildPolicy(observedMAP, mpol, celexceptions)
		_, err := policyClient.Create(ctx, observedMAP, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create mutatingadmissionpolicy %s: %w", observedMAP.GetName(), err)
		}
	} else {
		_, err := controllerutils.Update(ctx, observedMAP, policyClient,
			func(observed TPolicy) error {
				buildPolicy(observed, mpol, celexceptions)
				return nil
			})
		if err != nil {
			return fmt.Errorf("failed to update mutatingadmissionpolicy %s: %w", observedMAP.GetName(), err)
		}
	}

	if observedMAPbinding.GetResourceVersion() == "" {
		buildBinding(observedMAPbinding, mpol)
		_, err := bindingClient.Create(ctx, observedMAPbinding, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create mutatingadmissionpolicybinding %s: %w", observedMAPbinding.GetName(), err)
		}
	} else {
		_, err := controllerutils.Update(ctx, observedMAPbinding, bindingClient,
			func(observed TBinding) error {
				buildBinding(observed, mpol)
				return nil
			})
		if err != nil {
			return fmt.Errorf("failed to update mutatingadmissionpolicybinding %s: %w", observedMAPbinding.GetName(), err)
		}
	}

	c.updatePolicyStatus(ctx, genericPolicy, true, "")
	return nil
}
