package kube

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/listers/apps/v1"
	v1 "k8s.io/client-go/listers/core/v1"
	discoverylisters "k8s.io/client-go/listers/discovery/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfig *rest.Config
var clientset *kubernetes.Clientset

var PodLister v1.PodLister
var NodeLister v1.NodeLister
var ReplicaSetLister appsv1.ReplicaSetLister
var DeploymentLister appsv1.DeploymentLister
var NamespaceLister v1.NamespaceLister
var EndpointSliceLister discoverylisters.EndpointSliceLister

func InitKubeConfig() {
	var err error
	kubeConfig, err = rest.InClusterConfig()
	if err != nil {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", "/home/ubuntu/.kube/config.yaml")
		if err != nil {
			log.Fatal("Failed to create kubeconfig from local file:", err)
		}
		log.Println("Using kubeconfig from local file")
	} else {
		log.Println("Using in-cluster kubeconfig")
	}
}

func InitClientset() {
	var err error
	clientset, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Kubernetes clientset initialized successfully")
}

func InitLister(stopCh <-chan struct{}) {
	factory := informers.NewSharedInformerFactory(clientset, 0)
	podInformer := factory.Core().V1().Pods()
	PodLister = podInformer.Lister()

	nodeInformer := factory.Core().V1().Nodes()
	NodeLister = nodeInformer.Lister()

	replicaSetInformer := factory.Apps().V1().ReplicaSets()
	ReplicaSetLister = replicaSetInformer.Lister()

	deploymentInformer := factory.Apps().V1().Deployments()
	DeploymentLister = deploymentInformer.Lister()

	namespaceInformer := factory.Core().V1().Namespaces()
	NamespaceLister = namespaceInformer.Lister()

	endpointSliceInformer := factory.Discovery().V1().EndpointSlices()
	EndpointSliceLister = endpointSliceInformer.Lister()

	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)
}

func ListCollectorIP() {
	endpoints, _ := clientset.CoreV1().Endpoints("metrics-server-ns").List(context.TODO(), metav1.ListOptions{})
	for _, endpoint := range endpoints.Items {
		for _, subset := range endpoint.Subsets {
			for _, address := range subset.Addresses {
				log.Printf("Collector IP: %s", address.IP)
			}
		}
	}
}
