package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreatePod(imageName string) error {
	return nil
}

func ListPods(namespace string) ([]corev1.Pod, error) {
	// Load the Kubernetes configuration from the default location
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// List all pods in the default namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pods.Items, nil
}

func GetLogs(podName string) (string, error) {
	return "", nil
}

func ExtractFields(pods []corev1.Pod) (ret []PodResponse) {
	for _, pod := range pods {
		ret = append(ret, PodResponse{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			UID:       string(pod.UID),
			Labels:    pod.Labels,
		})
	}
	return
}

type PodResponse struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid"`
	Labels    map[string]string `json:"labels"`
}
