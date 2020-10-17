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
	"strconv"
	"time"

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
	hello := &crdv1alpha1.Hello{}
	var err error

	// 如果没有找到对应的资源，则表明已经删除，直接返回。否则，返回错误，重新进入调协过程
	if err = r.Get(ctx, req.NamespacedName, hello); err != nil {
		if errors.IsNotFound(err) {
			log.Info("RDS resource not found, ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}

		log.Error(err, "Error when getting resource")
		return ctrl.Result{}, err
	}

	// DeletionTimestamp不为0表示是非删除操作
	if hello.ObjectMeta.DeletionTimestamp.IsZero() {
		// condition为空，表示还没有开始倒计时, 开始初始化工作
		if hello.Status.Condition == "" {
			log.Info("Create Hello")
			t := strconv.Itoa(hello.Spec.Times)
			log.Info(t)
			hello.ObjectMeta.Finalizers = append(hello.ObjectMeta.Finalizers, "hello.finalizer.io")
			if err = r.Update(ctx, hello); err != nil {
				log.Error(err, "create resource error")
				return ctrl.Result{}, err
			}

			hello.Status.Condition = "Running"
			if err = r.Status().Update(ctx, hello); err != nil {
				log.Error(err, "update status error")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		// condition不为空，则表示已经开始倒计时
		log.Info("Update Hello Spec times")
		// 如果次数达到10次，就执行删除操作
		if hello.Spec.Times == 10 {
			hello.Status.Condition = "Deleting"
			if err = r.Status().Update(ctx, hello); err != nil {
				log.Error(err, "update status error")
				return ctrl.Result{}, err
			}

			log.Info("Delete Hello")
			if err = r.Delete(ctx, hello); err != nil {
				log.Error(err, "delete hello error")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		// 次数每次递增1
		hello.Spec.Times++
		time.Sleep(5 * time.Second)

		log.Info("Show Hello Spec times", "times", hello.Spec.Times)
		if err = r.Update(ctx, hello); err != nil {
			log.Error(err, "create resource error")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// 方便截图
	time.Sleep(5 * time.Second)

	// 将finalizers清空
	hello.ObjectMeta.Finalizers = []string{}
	if err = r.Update(ctx, hello); err != nil {
		log.Error(err, "creating resource error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *HelloReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1alpha1.Hello{}).
		Complete(r)
}
