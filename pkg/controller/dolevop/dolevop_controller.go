package dolevop

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	
	dolevgroupv1alpha1 "github.com/opdolev/pkg/apis/dolevgroup/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"


	appsv1 "k8s.io/api/apps/v1"

	routev1 "github.com/openshift/api/route/v1"
)

var log = logf.Log.WithName("controller_dolevop")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DolevOp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDolevOp{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dolevop-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DolevOp
	err = c.Watch(&source.Kind{Type: &dolevgroupv1alpha1.DolevOp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dolevgroupv1alpha1.DolevOp{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to ConfigMap
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dolevgroupv1alpha1.DolevOp{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to Route
	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &dolevgroupv1alpha1.DolevOp{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner DolevOp
//	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
//		IsController: true,
//		OwnerType:    &dolevgroupv1alpha1.DolevOp{},
//	})
//	if err != nil {
//		return err
//	}


	return nil
}

// blank assignment to verify that ReconcileDolevOp implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDolevOp{}

// ReconcileDolevOp reconciles a DolevOp object
type ReconcileDolevOp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DolevOp object and makes changes based on the state read
// and what is in the DolevOp.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDolevOp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Message")

	// Fetch the instance
	instance := &dolevgroupv1alpha1.DolevOp{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// Check if configmap for websites list already exists, if not create a new one
	reconcileResult, err := r.manageConfigMap(instance, reqLogger)
	if err != nil {
		instance.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}

	//Check if service already exists, if not create a new one
	reconcileResult, err = r.manageService(instance, reqLogger)
	if err != nil {
		instance.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}
	//Check if route already exists, if not create a new one
	reconcileResult, err = r.manageRoute(instance, reqLogger)
	if err != nil {
		instance.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}
	//Check if deployment already exists, if not create a new one
	reconcileResult, err = r.manageDeployment(instance, reqLogger)
	if err != nil {
		instance.Status.Message = fmt.Sprintf("%v", err)
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			reqLogger.Error(err, "Failed to update CR status")
		}
		return *reconcileResult, err
	} else if err == nil && reconcileResult != nil {
		// In case requeue required
		return *reconcileResult, nil
	}

	instance.Status.Message = "All good"
	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		reqLogger.Error(err, "Failed to update CR status")
	}
	return reconcile.Result{}, nil
}

// Reconcile loop resources managers functions
func (r *ReconcileDolevOp) manageDeployment(instance *dolevgroupv1alpha1.DolevOp, reqLogger logr.Logger) (*reconcile.Result, error) {
	deployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		serverDeployment, err := r.deploymentForWebServer(instance)
		if err != nil {
			reqLogger.Error(err, "error getting server deployment")
			return &reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new server deployment.", "Deployment.Namespace", serverDeployment.Namespace, "Deployment.Name", serverDeployment.Name)
		err = r.client.Create(context.TODO(), serverDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Server Deployment.", "Deployment.Namespace", serverDeployment.Namespace, "Deployment.Name", serverDeployment.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get server deployment.")
		return &reconcile.Result{}, err
	}
	return nil, nil
}

func (r *ReconcileDolevOp) manageRoute(instance *dolevgroupv1alpha1.DolevOp, reqLogger logr.Logger) (*reconcile.Result, error) {
	//Check if route already exists, if not create a new one
	route := &routev1.Route{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, route)
	if err != nil && errors.IsNotFound(err) {
		serverRoute, err := r.routeForWebServer(instance)
		if err != nil {
			reqLogger.Error(err, "error getting server route")
			return &reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new route.", "Route.Namespace", serverRoute.Namespace, "Router.Name", serverRoute.Name)
		err = r.client.Create(context.TODO(), serverRoute)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Server Route.", "Route.Namespace", serverRoute.Namespace, "Route.Name", serverRoute.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get server route.")
		return &reconcile.Result{}, err
	}
	return nil, nil
}

func (r *ReconcileDolevOp) manageService(instance *dolevgroupv1alpha1.DolevOp, reqLogger logr.Logger) (*reconcile.Result, error) {
	service := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, service)
	if err != nil && errors.IsNotFound(err) {
		err := r.serviceForWebServer(instance, service)
		if err != nil {
			reqLogger.Error(err, "error getting server service")
			return &reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Server Service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get server service.")
		return &reconcile.Result{}, err
	} else {
		err := r.serviceForWebServer(instance, service)
		if err != nil {
			reqLogger.Error(err, "error getting server service")
			return &reconcile.Result{}, err
		}
		err = r.client.Update(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Server Service.", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return &reconcile.Result{}, err
		}
	}
	return nil, nil
}

func (r *ReconcileDolevOp) manageConfigMap(instance *dolevgroupv1alpha1.DolevOp, reqLogger logr.Logger) (*reconcile.Result, error) {
	cm := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, cm)
	if err != nil && errors.IsNotFound(err) {
		websitesCm, err := r.configMapForWebServer(instance)
		if err != nil {
			reqLogger.Error(err, "error getting websites ConfigMap")
			return &reconcile.Result{}, err
		}
		reqLogger.Info("Creating a new cm.", "ConfigMap.Namespace", websitesCm.Namespace, "ConfigMap.Name", websitesCm.Name)
		err = r.client.Create(context.TODO(), websitesCm)
		if err != nil {
			reqLogger.Error(err, "Failed to create new ConfigMap.", "ConfigMap.Namespace", websitesCm.Namespace, "ConfigMap.Name", websitesCm.Name)
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get configmap.")
		return &reconcile.Result{}, err
	}

	// Check if CM sync is required
	syncRequired, err := r.syncConfigMapForWebServer(instance, cm)
	if err != nil {
		reqLogger.Error(err, "Error during syncing ConfigMap.")
		return &reconcile.Result{}, err
	}
	// If CM website sync required, sync the CM
	if syncRequired {
		err = r.client.Update(context.TODO(), cm)
		if err != nil {
			reqLogger.Error(err, "Error during updating CM")
			return &reconcile.Result{}, err
		}
	}
	return nil, nil
}

// Resources creation functions
func (r *ReconcileDolevOp) deploymentForWebServer(instance *dolevgroupv1alpha1.DolevOp) (*appsv1.Deployment, error) {
	var replicas int32
	replicas = 1
	labels := map[string]string{
		"app": instance.Name,
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": instance.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   instance.Name,
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           "docker.io/dimssss/nginx-for-ocp:0.1",
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "website",
									MountPath: "/opt/app-root/src",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "website",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: instance.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, dep, r.scheme); err != nil {
		log.Error(err, "Error set controller reference for server deployment")
		return nil, err
	}
	return dep, nil
}

func (r *ReconcileDolevOp) serviceForWebServer(instance *dolevgroupv1alpha1.DolevOp, service *corev1.Service) error {
	labels := map[string]string{
		"app": instance.Name,
	}
	ports := []corev1.ServicePort{
		{
			Name: "https",
			Port: 8080,
		},
	}
	service.ObjectMeta.Name = instance.Name
	service.ObjectMeta.Namespace = instance.Namespace
	service.ObjectMeta.Labels = labels
	service.Spec.Selector = map[string]string{"app": instance.Name}
	service.Spec.Ports = ports
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		log.Error(err, "Error set controller reference for server service")
		return err
	}
	return nil
}

func (r *ReconcileDolevOp) routeForWebServer(instance *dolevgroupv1alpha1.DolevOp) (*routev1.Route, error) {
	labels := map[string]string{
		"app": instance.Name,
	}
	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationEdge,
			},
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: instance.Name,
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		log.Error(err, "Error set controller reference for server route")
		return nil, err
	}
	return route, nil
}

func (r *ReconcileDolevOp) configMapForWebServer(instance *dolevgroupv1alpha1.DolevOp) (*corev1.ConfigMap, error) {
	labels := map[string]string{
		"app": instance.Name,
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{"index.html": "Welcome to my website! <br> My name is " + instance.Spec.Name + " And my message to you all is " + instance.Spec.Message},
	}

	if err := controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
		log.Error(err, "Error set controller reference for configmap")
		return nil, err
	}
	return cm, nil
}

func (r *ReconcileDolevOp) syncConfigMapForWebServer(instance *dolevgroupv1alpha1.DolevOp, cm *corev1.ConfigMap) (syncRequired bool, err error) {
	if cm.Data["index.html"] != "Welcome to my website! <br> My name is " + instance.Spec.Name + " And my message to you all is " + instance.Spec.Message {
		log.Info("Message in CR spec not the same as in CM, gonna update website cm")
		cm.Data["index.html"] = "Welcome to my website! <br> My name is " + instance.Spec.Name + " And my message to you all is " + instance.Spec.Message
		return true, nil
	}
	log.Info("No sync required, the message didn't changed")
	return false, nil
}
