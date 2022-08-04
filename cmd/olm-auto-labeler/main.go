package main

import (
	. "github.com/openshift/operator-framework-olm/pkg/olm-auto-labeler"
	operatorscheme "github.com/operator-framework/api/pkg/operators/install"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

func main() {
	logf.SetLogger(zap.New())

	var log = logf.Log.WithName("olm-auto-labeler")

	scheme := runtime.NewScheme()
	operatorscheme.Install(scheme)
	utilruntime.Must(v1.AddToScheme(scheme))

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "could not create manager")
		os.Exit(1)
	}

	err = builder.
		ControllerManagedBy(mgr).
		For(&v1.Namespace{}).
		Watches(&source.Kind{Type: &v1alpha1.ClusterServiceVersion{}}, handler.EnqueueRequestsFromMapFunc(func(o client.Object)[] reconcile.Request{
			csv, ok := o.(*v1alpha1.ClusterServiceVersion)
			if !ok {
				return nil
			}
			log.Info("CSV operation: "+csv.GetNamespace())
			return []reconcile.Request{}
		})).WithEventFilter(predicate.NewPredicateFuncs(func (o client.Object) bool {
			if o.GetObjectKind().GroupVersionKind().Kind == "namespace" {
				return true
			}
			return strings.HasPrefix(o.GetNamespace(), "openshift")
		})).Complete(NewOLMAutoLabelerReconciler(mgr.GetClient(), WithLogger(log)))

	if err != nil {
		log.Error(err, "could not create controller")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "could not start manager")
		os.Exit(1)
	}
}
