package main

import (
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Server holds all the config, resources and handlers required to run/serve
// a web-server and handle all the endpoints
type Server struct {
	router    *http.ServeMux
	clientset *kubernetes.Clientset
}

// PodResponse is a filtered view of Pod information returned by the
// kubernetes api for the user
type PodResponse struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid"`
	Labels    map[string]string `json:"labels"`
	Status    PodStatus         `json:"status"`
}

type PodStatus struct {
	State        corev1.ContainerState `json:"state"`
	Ready        bool                  `json:"ready"`
	RestartCount int32                 `json:"restart-count"`
	Image        string                `json:"image"`
	Started      *bool                 `json:"started"`
}

// CreatePodRequest is the input user must send for /create endpoint
type CreatePodRequest struct {
	Image     string `json:"image"`
	Namespace string `json:"namespace"`
	PodName   string `json:"pod_name"`
}

// GetLogsRequest is the input user must send for getting logs of a container
type GetLogsRequest struct {
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
}
