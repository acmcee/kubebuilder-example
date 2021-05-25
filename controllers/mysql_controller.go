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
	databasev1 "fordba.com/kubebuilder-example/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	log := r.Log.WithValues("mysql", req.NamespacedName)

	// your logic here

	objMySQL := &databasev1.MySQL{}
	if err := r.Get(ctx, req.NamespacedName, objMySQL); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("MySQL Kind not found")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Get mysql from Kubebuilder ", "InstanceName", objMySQL.Spec.InstanceName,
		"IP", objMySQL.Spec.IP, "Port", objMySQL.Spec.Port)

	objdeployment := &appsv1.Deployment{}
	dmFound := true
	if err := r.Get(ctx, req.NamespacedName, objdeployment); err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		dmFound = false
	}

	if dmFound {
		// update
	} else {
		if err := r.CreateMySQLDeployMent(ctx, objMySQL, objdeployment); err != nil {
			log.Error(err, "CreateMySQLDeployMent error")
		}
	}

	// 初始化 CR 的 Status 为 Running
	objMySQL.Status.Status = "Running"
	if err := r.Status().Update(ctx, objMySQL); err != nil {
		log.Error(err, "unable to update status")
	}

	return ctrl.Result{}, nil
}

func (r *MySQLReconciler) CreateMySQLDeployMent(ctx context.Context, vMySQL *databasev1.MySQL, vDeployMent *appsv1.Deployment) error {
	metadata := vMySQL.ObjectMeta.DeepCopy()
	vDeployMent.ObjectMeta = metav1.ObjectMeta{
		Name:      metadata.GetName(),
		Namespace: metadata.GetNamespace(),
	}
	// deployment labels


	mysqlLabel := vMySQL.ObjectMeta.GetLabels()
	vDeployMent.ObjectMeta.SetLabels(mysqlLabel)

	vDeployMent.Spec = appsv1.DeploymentSpec{
		Replicas: &vMySQL.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels:      mysqlLabel,
		},
		Template: v1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{
			Name:                       "",
			GenerateName:               "",
			Namespace:                  "",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     mysqlLabel,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ClusterName:                "",
			ManagedFields:              nil,
		},
			Spec: v1.PodSpec{
				Volumes:        nil,
				InitContainers: nil,
				Containers: []v1.Container{v1.Container{
					Name:                     vMySQL.ObjectMeta.GetName(),
					Image:                    vMySQL.Spec.Image,
					Command:                  nil,
					Args:                     nil,
					WorkingDir:               "",
					Ports:                    nil,
					EnvFrom:                  nil,
					Env:                      nil,
					Resources:                v1.ResourceRequirements{},
					VolumeMounts:             nil,
					VolumeDevices:            nil,
					LivenessProbe:            nil,
					ReadinessProbe:           nil,
					StartupProbe:             nil,
					Lifecycle:                nil,
					TerminationMessagePath:   "",
					TerminationMessagePolicy: "",
					ImagePullPolicy:          "",
					SecurityContext:          nil,
					Stdin:                    false,
					StdinOnce:                false,
					TTY:                      false,
				}},
				EphemeralContainers:           nil,
				RestartPolicy:                 "",
				TerminationGracePeriodSeconds: nil,
				ActiveDeadlineSeconds:         nil,
				DNSPolicy:                     "",
				NodeSelector:                  nil,
				ServiceAccountName:            "",
				DeprecatedServiceAccount:      "",
				AutomountServiceAccountToken:  nil,
				NodeName:                      "",
				HostNetwork:                   false,
				HostPID:                       false,
				HostIPC:                       false,
				ShareProcessNamespace:         nil,
				SecurityContext:               nil,
				ImagePullSecrets:              nil,
				Hostname:                      "",
				Subdomain:                     "",
				Affinity:                      nil,
				SchedulerName:                 "",
				Tolerations:                   nil,
				HostAliases:                   nil,
				PriorityClassName:             "",
				Priority:                      nil,
				DNSConfig:                     nil,
				ReadinessGates:                nil,
				RuntimeClassName:              nil,
				EnableServiceLinks:            nil,
				PreemptionPolicy:              nil,
				Overhead:                      nil,
				TopologySpreadConstraints:     nil,
			}},
		Strategy:                appsv1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
	// 在删除MySQL 的时候，可以自动删除下面的deployment
	if err := ctrl.SetControllerReference(vMySQL, vDeployMent, r.Scheme); err != nil {
		return err
	}

	if err := r.Create(ctx, vDeployMent); err != nil {
		return err
	}

	return nil
}

func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1.MySQL{}).
		Complete(r)
}
