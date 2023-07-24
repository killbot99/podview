package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type ReadyReplicasChangedPredicate struct {
	predicate.Funcs
}

// predicate for only caring if the number of ready replicas have changed.
func (ReadyReplicasChangedPredicate) Update(e event.UpdateEvent) bool {
	old, ok := e.ObjectOld.(*appsv1.Deployment)
	if !ok {
		return false // in case someone uses this predicate for an object that isn't a deployment
	}
	new := e.ObjectNew.(*appsv1.Deployment)

	return !(old.Status.ReadyReplicas == new.Status.ReadyReplicas)
}
