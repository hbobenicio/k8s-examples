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
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Logger: logger,

		// this leader election will create a k8s resource called lease (not actually a lock)
		// this lease is shared between instances and is used to elect a leader.
		// the leader will be the only one receiving reconcile events for declared watches of your manager (by default)
		// the other instances will work as standby's
		LeaderElection:   true,
		LeaderElectionID: "my-controller-lease", // check it with kubectl -n my-controller get leases
	})
	if err != nil {
		return fmt.Errorf("unable to set up overall controller manager: %w", err)
	}

	// Setup a new controller to reconcile Namespaces
	logger.Info("setting up controller...")
	ctrl, err := controller.New("namespace-controller", mgr, controller.Options{
		Reconciler: &mycontroller.NamespaceReconciler{
			Client: mgr.GetClient(),
		},

		// MaxConcurrentReconciles is the maximum number of concurrent Reconciles which can be run. Defaults to 1
		// MaxConcurrentReconciles: ...,

		// NeedLeaderElection indicates whether the controller needs to use leader election. Defaults to true
		// NeedLeaderElection: ...,
	})
	if err != nil {
		return fmt.Errorf("unable to set up individual controller: %w", err)
	}

	// Watch Namespaces and enqueue Namespace object key
	logger.Info("setting up watches...")
	{
		namespaceSource := source.Kind(mgr.GetCache(), &v1.Namespace{}, &handler.TypedEnqueueRequestForObject[*v1.Namespace]{})
		if err := ctrl.Watch(namespaceSource); err != nil {
			return fmt.Errorf("unable to watch namespaces: %w", err)
		}
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
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	log.Log.Error(err, "failed to get hostname")
	// }

	rootLogger := zap.New(zap.UseFlagOptions(rootLoggerOpts)) //.
	// WithValues("namespace", os.Getenv("K8S_NAMESPACE")).
	// WithValues("pod.name", os.Getenv("K8S_POD_NAME"))

	controllerruntime.SetLogger(rootLogger)
}
