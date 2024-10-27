package main

import (
	"os"

	v1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	myoperator "myoperator/internal/my-operator"
)

func main() {
	//TODO missing topics:
	//     - health-checking: manager AddHealthzCheck() and AddReadyzCheck()
	//     - graceful shutdown: does controller-runtime Start already does that?
	//     - leader-elect: Enabling this will ensure there is only one active controller manager
	rootLoggerSetup()
	entryLog := log.Log.WithName("my-operator")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Logger: entryLog,
	})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up controller")
	c, err := controller.New("namespace-controller", mgr, controller.Options{
		Reconciler: &myoperator.NamespaceReconciler{Client: mgr.GetClient()},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch ReplicaSets and enqueue ReplicaSet object key
	if err := c.Watch(source.Kind(mgr.GetCache(), &v1.Namespace{}, &handler.TypedEnqueueRequestForObject[*v1.Namespace]{})); err != nil {
		entryLog.Error(err, "unable to watch Namespaces")
		os.Exit(1)
	}

	// Webhook example:
	// if err := builder.WebhookManagedBy(mgr).
	// 	For(&corev1.Pod{}).
	// 	// WithDefaulter(&podAnnotator{}).
	// 	// WithValidator(&podValidator{}).
	// 	Complete(); err != nil {
	// 	entryLog.Error(err, "unable to create webhook", "webhook", "Pod")
	// 	os.Exit(1)
	// }

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func rootLoggerSetup() {
	rootLoggerOpts := &zap.Options{
		Development: true,
	}

	//TODO get replica ID from uname syscall

	rootLogger := zap.New(zap.UseFlagOptions(rootLoggerOpts)).
		WithValues("namespace", os.Getenv("K8S_NAMESPACE")).
		WithValues("pod.name", os.Getenv("K8S_POD_NAME"))

	controllerruntime.SetLogger(rootLogger)
}
