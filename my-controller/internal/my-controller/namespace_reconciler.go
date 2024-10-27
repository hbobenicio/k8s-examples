package mycontroller

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type NamespaceReconciler struct {
	// Client is the APIServer client
	Client client.Client
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	//NOTE I thing logger with all values from context is just too much info. you can use it, but I prefer not
	//     as default
	// logger := log.FromContext(ctx)
	reconcileID := controller.ReconcileIDFromContext(ctx)
	logger := log.Log.WithValues("reconcileID", reconcileID)

	logger.Info("got new reconciling request event", "namespace", req.Name)

	// Getting current system state (be mindful that this client uses caching)
	var namespace v1.Namespace
	if err := r.Client.Get(ctx, req.NamespacedName, &namespace); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get namespace %s: %w", req, err)
	}

	// Define desired state
	targetLabel := "java"
	desiredLabelValue := "sucks"
	if labelValue, found := namespace.Labels[targetLabel]; found && labelValue == desiredLabelValue {
		logger.Info("reconciliation success: noop")
		return reconcile.Result{}, nil
	}

	namespace.Labels[targetLabel] = desiredLabelValue

	// Apply your desired state to the system
	if err := r.Client.Update(ctx, &namespace); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to update namespace %s: %w", req, err)
	}

	logger.Info("reconciliation success: a new easter-egg was added to the labels of the namespace", "namespace", req.Name)
	return reconcile.Result{}, nil
}
