package admissionpolicygenerator

import (
	controllerutils "github.com/kyverno/kyverno/pkg/utils/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// MAP / MAP binding handlers use metav1.Object so they work for any served admissionregistration.k8s.io MAP API version.
func (c *controller) addMAP(obj interface{}) {
	acc, ok := obj.(metav1.Object)
	if !ok {
		return
	}
	c.enqueueMAPFromMeta(acc)
}

func (c *controller) updateMAP(oldObj, newObj interface{}) {
	oldAcc, ok := oldObj.(metav1.Object)
	if !ok {
		return
	}
	newAcc, ok := newObj.(metav1.Object)
	if !ok {
		return
	}
	if oldAcc.GetResourceVersion() == newAcc.GetResourceVersion() {
		return
	}
	c.enqueueMAPFromMeta(newAcc)
}

func (c *controller) deleteMAP(obj interface{}) {
	acc, ok := obj.(metav1.Object)
	if !ok {
		return
	}
	c.enqueueMAPFromMeta(acc)
}

func (c *controller) enqueueMAPFromMeta(m metav1.Object) {
	if m == nil {
		return
	}
	refs := m.GetOwnerReferences()
	if len(refs) != 1 || refs[0].Kind != "MutatingPolicy" {
		return
	}
	mpol, err := c.mpolLister.Get(refs[0].Name)
	if err != nil {
		return
	}
	c.enqueueMP(mpol)
}

func (c *controller) addMAPbinding(obj interface{}) {
	acc, ok := obj.(metav1.Object)
	if !ok {
		return
	}
	c.enqueueMAPbindingFromMeta(acc)
}

func (c *controller) updateMAPbinding(oldObj, newObj interface{}) {
	oldAcc, ok := oldObj.(metav1.Object)
	if !ok {
		return
	}
	newAcc, ok := newObj.(metav1.Object)
	if !ok {
		return
	}
	if oldAcc.GetResourceVersion() == newAcc.GetResourceVersion() {
		return
	}
	c.enqueueMAPbindingFromMeta(newAcc)
}

func (c *controller) deleteMAPbinding(obj interface{}) {
	acc, ok := obj.(metav1.Object)
	if !ok {
		return
	}
	c.enqueueMAPbindingFromMeta(acc)
}

func (c *controller) enqueueMAPbindingFromMeta(mb metav1.Object) {
	if mb == nil {
		return
	}
	refs := mb.GetOwnerReferences()
	if len(refs) != 1 || refs[0].Kind != "MutatingPolicy" {
		return
	}
	mpol, err := c.mpolLister.Get(refs[0].Name)
	if err != nil {
		return
	}
	c.enqueueMP(mpol)
}

func registerMAPInformerHandlers(c *controller, mapPolicyInformer, mapBindingInformer cache.SharedInformer) {
	if mapPolicyInformer != nil {
		if _, err := controllerutils.AddEventHandlers(mapPolicyInformer, c.addMAP, c.updateMAP, c.deleteMAP); err != nil {
			logger.Error(err, "failed to register MutatingAdmissionPolicy event handlers")
		}
	}
	if mapBindingInformer != nil {
		if _, err := controllerutils.AddEventHandlers(mapBindingInformer, c.addMAPbinding, c.updateMAPbinding, c.deleteMAPbinding); err != nil {
			logger.Error(err, "failed to register MutatingAdmissionPolicyBinding event handlers")
		}
	}
}
