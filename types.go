package main

// PodResponse is a filtered view of Pod information returned by the
// kubernetes api for the user
type PodResponse struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid"`
	Labels    map[string]string `json:"labels"`
}

// CreatePodRequest is the input user must send for /create endpoint
type CreatePodRequest struct {
	Image     string `json:"image"`
	Namespace string `json:"namespace"`
}

// GetLogsRequest is the input user must send for getting logs of a container
type GetLogsRequest struct {
	PodName   string `json:"pod_name"`
	Namespace string `json:"namespace"`
}
