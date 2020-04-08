package main

import (
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"os"
	"path/filepath"
)

var (
	machinesetGVR = schema.GroupVersionResource{
		Group:    "machine.openshift.io",
		Version:  "v1beta1",
		Resource: "machines",
	}
)
var (
	masterURL  string
	kubeconfig string
)

func main() {

	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig(), "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

	klog.InitFlags(nil)

	flag.Parse()

	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			klog.Fatal("Error building kubeconfig: %s", err.Error())
		}
	}
	dynClient, errClient := dynamic.NewForConfig(cfg)
	if errClient != nil {
		klog.Fatal("Error received creating client %v", errClient)
	}

	crdClient := dynClient.Resource(machinesetGVR)

	crd, errCrd := crdClient.Namespace("openshift-machine-api").Get("aputtur-worker-0-rz6v5", metav1.GetOptions{})
	//	crd, errCrd := crdClient.Namespace("openshift-machine-api").List( metav1.ListOptions{})

	if errCrd != nil {
		klog.Fatal("Error getting CRD %v", errCrd)
	}
	klog.Info("Got CRD: %v", crd)
}

func defaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
