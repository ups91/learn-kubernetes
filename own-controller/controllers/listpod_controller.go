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
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1alpha1 "learn-kubernetes/own-controller/api/v1alpha1"
)

// ListPodReconciler reconciles a ListPod object
type ListPodReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.begin,resources=listpods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.begin,resources=listpods/status,verbs=get;update;patch

func (r *ListPodReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	bkgcntx := context.Background()
	_ = r.Log.WithValues("listpod", req.NamespacedName)

	// your logic here
	fmt.Println("############ Node List Goes Next ################")
	//	_ = r.Log.WithValues("listpod", "############################################")
	printPods := v1.PodList{}
	printNodes := v1.NodeList{}

	if err := r.List(bkgcntx, &printNodes); err == nil {
		for _, i := range printNodes.Items {
			fmt.Println(i.ObjectMeta.Name)
		}
	}
	fmt.Println("-----------  Pod List Goes Next ----------------")

	if err := r.List(bkgcntx, &printPods); err == nil {
		for _, i := range printPods.Items {
			fmt.Printf("%s\t==>\t%s\n", i.Status.PodIP, i.ObjectMeta.Name)
		}
	}
	fmt.Println("#####################################################")
	return ctrl.Result{}, nil
}

func (r *ListPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.ListPod{}).
		Complete(r)
}
