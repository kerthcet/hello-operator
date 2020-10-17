/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	crdv1alpha1 "hello-operator/api/v1alpha1"
)

// HelloReconciler reconciles a Hello object
type HelloReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=crd.hello.com,resources=hellos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crd.hello.com,resources=hellos/status,verbs=get;update;patch

func (r *HelloReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("hello", req.NamespacedName)

	// your logic here
	hello := crdv1alpha1.Hello{}

	if err := r.Get(ctx, hello.NamespacedName, hello); err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "RDS resource not found, ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1alpha1.Hello{}).
		Complete(r)
}
