package main

import (
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	clientset, err := getClientSet()
	if err != nil {
		panic("error getting client-set for kubernetes")
	}

	s := Server{
		router:    http.NewServeMux(),
		clientset: clientset,
	}
	s.RegisterRoutes()

	port := ResolvePort()
	log.Println("will server on port:", port)

	panic(s.Serve(port))
}

func getClientSet() (*kubernetes.Clientset, error) {
	// Load the Kubernetes configuration from the default location
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
