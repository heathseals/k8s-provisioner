package main

import (
	"fmt"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var options struct {
	Verbose    []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Namespace  string `short:"n" long:"namespace" description:"Namespace to create" required:"true"`
	Kubeconfig string `short:"k" long:"kubeconfig" description:"absolute path to the kubeconfig file (default: ~/.kube/config)"`
}

func main() {
	flags.Parse(&options)
	if options.Kubeconfig == "" {
		home := homeDir()
		options.Kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", options.Kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: options.Namespace}}

	_, err = clientset.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created namespace %s\n", options.Namespace)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
