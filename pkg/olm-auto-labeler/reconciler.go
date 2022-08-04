package olm_auto_labeler

import (
	"context"
	"fmt"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type OLMAutoLabelerReconciler struct {
	client.Client
	logger logr.Logger
}

type OLMAutoLabellerReconcilerOption func(reconciler OLMAutoLabelerReconciler)

func WithLogger(logger logr.Logger) OLMAutoLabellerReconcilerOption {
	return func(reconciler OLMAutoLabelerReconciler) {
		reconciler.logger = logger
	}
}

func NewOLMAutoLabelerReconciler(client client.Client, options ...OLMAutoLabellerReconcilerOption) *OLMAutoLabelerReconciler {
	reconciler := OLMAutoLabelerReconciler{
		Client: client,
	}
	for _, option := range options {
		option(reconciler)
	}

	return &reconciler
}

func (a *OLMAutoLabelerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := builder.
		ControllerManagedBy(mgr).
		For(&v1.Namespace{}, ).
		Watches(&source.Kind{Type: &v1alpha1.ClusterServiceVersion{}}, handler.EnqueueRequestsFromMapFunc(func(o client.Object)[] reconcile.Request{
			csv, ok := o.(*v1alpha1.ClusterServiceVersion)
			if !ok {
				return nil
			}
			a.logger.Info("CSV operation: "+csv.GetNamespace())
			return []reconcile.Request{}
		})).WithEventFilter(predicate.NewPredicateFuncs(func (o client.Object) bool {
		if o.GetObjectKind().GroupVersionKind().Kind == "namespace" {
			return false
		}
		if !strings.HasPrefix(o.GetNamespace(), "openshift") {
			return false
		}
		a.Get(context.TODO(), types.NamespacedName{Name: o.GetNamespace(), })
		a.List(context.TODO(), nil, )
		return strings.HasPrefix(o.GetNamespace(), "openshift")
	})).Complete(a)

	return err

}

func (a *OLMAutoLabelerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	//rs := &appsv1.ReplicaSet{}
	//err := a.Get(ctx, req.NamespacedName, rs)
	//if err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//pods := &corev1.PodList{}
	//err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
	//if err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//rs.Labels["pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	//err = a.Update(ctx, rs)
	//if err != nil {
	//	return reconcile.Result{}, err
	//}
	a.logger.Info(fmt.Sprintf("Got a request: %s", req.NamespacedName))

	return reconcile.Result{}, nil
}
