package main

import (
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mycontroller "github.com/hbobenicio/k8s-examples/my-controller/internal/my-controller"
)

const (
	controllerName = "my-controller"
)

// # TODO missing topics
//   - graceful shutdown: does controller-runtime Start already does that?
//   - leader election:
//     Enabling this will ensure there is only one active controller manager.
//     Highly recommended reading: https://sklar.rocks/kubernetes-leader-election/
func main() {
	if err := run(); err != nil {
		log.Log.Error(err, controllerName+" failed")
		os.Exit(1)
	}
}

func run() error {
	loggingSetup()
	logger := log.Log

	// Setup a Manager
	logger.Info("setting up manager...")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Logger: logger})
	if err != nil {
		return fmt.Errorf("unable to set up overall controller manager: %w", err)
	}

	// Setup a new controller to reconcile ReplicaSets
	logger.Info("setting up controller...")
	ctrl, err := controller.New("namespace-controller", mgr, controller.Options{
		Reconciler: &mycontroller.NamespaceReconciler{Client: mgr.GetClient()},
	})
	if err != nil {
		return fmt.Errorf("unable to set up individual controller: %w", err)
	}

	// Watch ReplicaSets and enqueue ReplicaSet object key
	logger.Info("setting up watches...")
	if err := ctrl.Watch(source.Kind(mgr.GetCache(), &v1.Namespace{}, &handler.TypedEnqueueRequestForObject[*v1.Namespace]{})); err != nil {
		return fmt.Errorf("unable to watch namespaces: %w", err)
	}

	logger.Info("setting up health/liveness probe...")
	if err := mgr.AddHealthzCheck("liveness", healthz.Ping); err != nil {
		return fmt.Errorf("unable to add liveness checking: %w", err)
	}
	logger.Info("setting up readiness probe...")
	if err := mgr.AddReadyzCheck("readiness", healthz.Ping); err != nil {
		return fmt.Errorf("unable to add readiness checking: %w", err)
	}

	logger.Info("setup is done.")
	logger.Info("starting manager...")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	return nil
}

func loggingSetup() {
	rootLoggerOpts := &zap.Options{
		Development: true,
	}

	//TODO get replica ID from uname syscall

	rootLogger := zap.New(zap.UseFlagOptions(rootLoggerOpts)) //.
	// WithValues("namespace", os.Getenv("K8S_NAMESPACE")).
	// WithValues("pod.name", os.Getenv("K8S_POD_NAME"))

	controllerruntime.SetLogger(rootLogger)
}
