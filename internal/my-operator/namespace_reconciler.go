package myoperator

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type NamespaceReconciler struct {
	// Client is the APIServer client
	Client client.Client
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("got new reconciling request event", "reconcileRequest", req)

	// Getting current system state (be mindful that this client uses caching)
	var namespace v1.Namespace
	if err := r.Client.Get(ctx, req.NamespacedName, &namespace); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get namespace %s: %w", req, err)
	}
	logger.Info("got object", "namespace", namespace)

	// Define desired state
	if _, found := namespace.Labels["java"]; found {
		logger.Info("reconciliation success: noop")
		return reconcile.Result{}, nil
	}

	namespace.Labels["java"] = "sucks"

	// Apply your desired state to the system
	if err := r.Client.Update(ctx, &namespace); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to update namespace %s: %w", req, err)
	}

	logger.Info("reconciliation success: a new easter-egg was added to the labels of the namespace", "namespace", req.Name)
	return reconcile.Result{}, nil
}
