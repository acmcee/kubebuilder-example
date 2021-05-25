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
	"log"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasev1 "fordba.com/kubebuilder-example/api/v1"
)

// MySQLReconciler reconciles a MySQL object
type MySQLReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=database.fordba.com,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=database.fordba.com,resources=mysqls/status,verbs=get;update;patch

func (r *MySQLReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("mysql", req.NamespacedName)

	// your logic here

	obj := &databasev1.MySQL{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Println(err, "Unable to fetch object")
	} else {
		log.Println("get mysql from Kubebuilder ", obj.Spec.InstanceName, obj.Spec.IP, obj.Spec.Port)
	}

	// 初始化 CR 的 Status 为 Running
	obj.Status.Status = "Running"
	if err := r.Status().Update(ctx, obj); err != nil {
		log.Println(err, "unable to update status")
	}

	return ctrl.Result{}, nil
}

func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1.MySQL{}).
		Complete(r)
}
