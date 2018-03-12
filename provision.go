package main

import (
	"fmt"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var options struct {
	Verbose    []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Namespace  string `short:"n" long:"namespace" description:"Namespace to create" required:"true"`
	Kubeconfig string `short:"k" long:"kubeconfig" description:"absolute path to the kubeconfig file (default: ~/.kube/config)"`
	Username   string `short:"u" long:"username" description:"User to add to namespace" required:"true"`
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func createNamespace(clientset *kubernetes.Clientset) {
	nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: options.Namespace}}

	_, err := clientset.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created namespace %s\n", options.Namespace)
}

func createRoleBinding(clientset *kubernetes.Clientset) {
	rbSpec := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Username + "-rolebinding",
			Namespace: options.Namespace,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Name: options.Username,
				Kind: "User",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "admin",
		},
	}

	_, err := clientset.RbacV1().RoleBindings(options.Namespace).Create(rbSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created rolebinding %s-rolebinding in namespace %s\n", options.Username, options.Namespace)
}

func main() {
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}
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

	_, err = clientset.CoreV1().Namespaces().Get(options.Namespace, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Namespace %s not found, creating...\n", options.Namespace)
		createNamespace(clientset)
	} else {
		fmt.Printf("Namespace already exists, creating rolebinding...\n")
	}
	createRoleBinding(clientset)

}
