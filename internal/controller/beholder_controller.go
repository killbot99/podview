/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	monsterv1 "github.com/killbot99/podview/api/v1"
)

// BeholderReconciler reconciles a Beholder object
type BeholderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monster.podview.killbot99.io,resources=beholders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monster.podview.killbot99.io,resources=beholders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monster.podview.killbot99.io,resources=beholders/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Beholder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *BeholderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	Logger := log.FromContext(ctx)
	// get the current state of the beholder resource
	var beholder monsterv1.Beholder
	if err := r.Get(ctx, req.NamespacedName, &beholder); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Build Deployment Desired state from janky template yaml
	dy := fmt.Sprintf(PodviwerDeploymentTemplate,
		beholder.GetName(),
		beholder.GetNamespace(),
		beholder.Spec.ReplicaCount,
		beholder.GetName(),
		beholder.GetName(),
		beholder.Spec.Image.Registry,
		beholder.Spec.Image.Tag,
		beholder.GetName(),
		beholder.Spec.UI.Color,
		beholder.Spec.UI.Message,
		beholder.Spec.Resources.CPULimit.String(),
		beholder.Spec.Resources.MemoryLimit.String(),
		beholder.Spec.Resources.CPURequest.String(),
		beholder.Spec.Resources.MemoryRequest.String(),
	)
	d := &appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(strings.NewReader(dy), 100).Decode(d); err != nil {
		return ctrl.Result{}, err
	}
	controllerutil.SetControllerReference(&beholder, d, r.Scheme)

	Logger.Info("updating podview deployment", "name", d.GetName())
	if err := r.Update(ctx, d); err != nil {
		if apierrors.IsNotFound(err) {
			Logger.Info("podview deployment not found. creating.", "name", d.GetName())
			if err := r.Create(ctx, d); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	sy := fmt.Sprintf(PodviewerServiceTemplate,
		beholder.GetName(),
		beholder.GetNamespace(),
		beholder.GetName(),
	)
	s := &corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(strings.NewReader(sy), 100).Decode(s); err != nil {
		return ctrl.Result{}, err
	}
	controllerutil.SetControllerReference(&beholder, s, r.Scheme)
	Logger.Info("updating podview service", "name", s.GetName())
	if err := r.Update(ctx, s); err != nil {
		if apierrors.IsNotFound(err) {
			Logger.Info("podview service not found. creating.", "name", s.GetName())
			if err := r.Create(ctx, s); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BeholderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monsterv1.Beholder{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&appsv1.Deployment{}, builder.WithPredicates(ReadyReplicasChangedPredicate{})).
		Complete(r)
}

// I know this is ugly but it's a quick way to get this deployment up, and i
// care more about my controller logic than how to import a 3rd party deployment template.
var PodviwerDeploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podinfo-%s
  namespace: %s
spec:
  minReadySeconds: 3
  revisionHistoryLimit: 5
  progressDeadlineSeconds: 60
  replicas: %d
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  selector:
    matchLabels:
      app: podinfo-%s
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9797"
      labels:
        app: podinfo-%s
    spec:
      containers:
      - name: podinfod
        image: %s:%s
        imagePullPolicy: IfNotPresent
        ports:
        - name: http
          containerPort: 9898
          protocol: TCP
        - name: http-metrics
          containerPort: 9797
          protocol: TCP
        - name: grpc
          containerPort: 9999
          protocol: TCP
        command:
        - ./podinfo
        - --port=9898
        - --port-metrics=9797
        - --grpc-port=9999
        - --grpc-service-name=podinfo-%s
        - --level=info
        - --random-delay=false
        - --random-error=false
        env:
        - name: PODINFO_UI_COLOR
          value: "%s"
        - name: PODINFO_UI_MESSAGE
          value: %s
        livenessProbe:
          exec:
            command:
            - podcli
            - check
            - http
            - localhost:9898/healthz
          initialDelaySeconds: 5
          timeoutSeconds: 5
        readinessProbe:
          exec:
            command:
            - podcli
            - check
            - http
            - localhost:9898/readyz
          initialDelaySeconds: 5
          timeoutSeconds: 5
        resources:
          limits:
            cpu: '%s'
            memory: '%s'
          requests:
            cpu: '%s'
            memory: '%s'
`

var PodviewerServiceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: podinfo-%s
  namespace: %s
spec:
  type: ClusterIP
  selector:
    app: podinfo-%s
  ports:
    - name: http
      port: 9898
      protocol: TCP
      targetPort: http
    - port: 9999
      targetPort: grpc
      protocol: TCP
      name: grpc
`
