STEPS
=====

Based on: https://learn.openshift.com/operatorframework/go-operator-podset/

$ operator-sdk version
operator-sdk version: "v1.10.1-ocp", commit: "58bf7213d22688d7d0cd4ccc7f4faf4ec5b33e3f", kubernetes version: "v1.21", go version: "go1.16.6", GOOS: "linux", GOARCH: "amd64"

$ operator-sdk init --domain=example.com --repo=github.com/redhat/podset-operator

Create the GVK, creates /api, /config and /controllers directories (by default
this is a single group api):
$ operator-sdk create api --group=app --version=v1alpha1 --kind=Podset --resource --controller

Define the API, modifying the spec and the status on:
$ cat api/v1alpha1/podset_types.go

Create the initial git repo:
$ git init .
$ git branch -M main
$ git remote add origin https://github.com/rauhersu/podset-operator.git

make generate generates code, like runtime.Object/DeepCopy implementations (See XXXdeepcopy.go)
$ make generate

make manifests generates Kubernetes object YAML, like CustomResourceDefinitions,
WebhookConfigurations, and RBAC roles (See config/crd/bases/XXX.yaml)
$ make manifests

See the validation through generated from the comment markers on _types.go:
~/gorepo/memcached-operator/config/crd/bases/cache.example.com_memcacheds.yaml

spec:
  ...
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        description: Memcached is the Schema for the memcacheds API

You can check the CRD is running ok:

$ oc new-project myproject
$ oc apply -f config/crd/bases/cache.example.com_memcacheds.yaml
$ oc get crd memcacheds.cache.example.com -o yaml
￼
So far we have created the API, let's go for the

1. 'reconcile' on controllers/XXX_controller.go
2. $ go mod tidy


Explain how the controller watches resources and how the reconcile loop is triggered.
If you’d like to skip this section, head to the deploy section to see how to run the
operator.Specify and understand the 'SetupWithManager' function (resources watched by
the controller)

The 'SetupWithManager()' function in controllers/memcached_controller.go specifies
how the controller is built to watch a CR and other resources that are owned and
managed by that controller.

import (
  ...
    appsv1 "k8s.io/api/apps/v1"
      ...
      )

func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
  return ctrl.NewControllerManagedBy(mgr).
      For(&cachev1alpha1.Memcached{}).
          Owns(&appsv1.Deployment{}).
              Complete(r)
              }

The 'NewControllerManagedBy()' provides a controller builder that allows various controller configurations.

'For(&cachev1alpha1.Memcached{})' specifies the Memcached type as the primary resource to watch. For each
Memcached type Add/Update/Delete event the reconcile loop will be sent a reconcile 'Request' (a namespace/name key)
for that Memcached object.

'Owns(&appsv1.Deployment{})' specifies the Deployments type as the secondary resource to watch. For each Deployment
type Add/Update/Delete event, the event handler will map each event to a reconcile Request for the owner of the Deployment.
Which in this case is the Memcached object for which the Deployment was created.

There are a number of other useful configurations that can be made when initialzing a controller.
For more details on these configurations consult the upstream builder and controller godocs.

Request: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Request


builder:
https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/builder#example-Builder

controller:
https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller

1.Set the max number of concurrent Reconciles for the controller via the MaxConcurrentReconciles option. Defaults to 1.
  https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller#Options

2.Filter watch events using predicates
  https://sdk.operatorframework.io/docs/building-operators/golang/references/event-filtering/

3.Choose the type of EventHandler to change how a watch event will translate to \
  reconcile requests for the reconcile loop. For operator relationships that are
  more complex than primary and secondary resources, the EnqueueRequestsFromMapFunc handler can be used to transform a watch event into an arbitrary set of reconcile requests.
https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/handler#hdr-EventHandlers
https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/handler#EnqueueRequestsFromMapFunc


Reconcile loop:

Reconcile loop
The reconcile function is responsible for enforcing the desired CR state on the actual state of the system. It runs each time an event occurs on a watched CR or resource, and will return some value depending on whether those states match or not.

In this way, every Controller has a Reconciler object with a Reconcile() method that implements the reconcile loop. The reconcile loop is passed the 'Request' argument which is a Namespace/Name key used to lookup the primary resource object, Memcached, from the cache:

import (
  ctrl "sigs.k8s.io/controller-runtime"

  cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
    ...
    )

func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
  _ = context.Background()
    ...

  // Lookup the Memcached instance for this reconcile request
    memcached := &cachev1alpha1.Memcached{}
      err := r.Get(ctx, req.NamespacedName, memcached)
        ...
        }

'Request': https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Request
'Client API doc': https://sdk.operatorframework.io/docs/building-operators/golang/references/client/

Filter watch events using predicates

Choose the type of EventHandler to change how a watch event will translate to reconcile requests for the reconcile loop. For operator relationships that are more complex than primary and secondary resources, the EnqueueRequestsFromMapFunc handler can be used to transform a watch event into an arbitrary set of reconcile requests.

Return options for a Reconciler:
The following are a few possible return options for a Reconciler:

With the error:
return ctrl.Result{}, err
Without an error:
return ctrl.Result{Requeue: true}, nil
Therefore, to stop the Reconcile, use:
return ctrl.Result{}, nil
Reconcile again after X time:
 return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil

Reconcile: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile

>> RBAC

Specify permissions and generate RBAC manifests
The controller needs certain RBAC permissions to interact with the resources it manages. These are specified via RBAC markers like the following:

//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
  ...
  }

The ClusterRole manifest at config/rbac/role.yaml is generated from the above markers via controller-gen with the following command: $ make manifests

Local:

$ WATCH_NAMESPACE=myproject make install run uninstall

Check the CRD was properly installed through the 'install' target above
oc $ get crd XXX -o yaml

$ oc project myproject
$ oc apply -f config/samples/app_v1alpha1_podset.yaml
$ oc get podset
NAME            DESIRED   AVAILABLE
podset-sample   3         3

$ oc patch memcached memcached-sample --type='json' -p '[{"op": "replace", "path": "/spec/size", "value":5}]'

Our Memcached controller creates pods containing OwnerReferences in their metadata
section. This ensures they will be removed upon deletion of the memcached-sample CR.

oc get pods -o yaml | grep ownerReferences -A10

oc delete memcached memcached-sample
oc get pods

And finally delete:

$ make uninstall

In cluster:

$ docker login
$ IMAGE_TAG_BASE=docker.io/rauherna/podset-operator make VERSION=0.0.1 docker-build docker-push
$ oc project myproject
$ IMAGE_TAG_BASE=docker.io/rauherna/podset-operator VERSION=0.0.1 make deploy

See the resources created:

$ oc get clusterserviceversion
$ oc get crd | grep podset
$ oc get sa | grep podset
$ oc get roles | grep podset
$ oc get rolebindings | grep podset
$ oc get deployments

