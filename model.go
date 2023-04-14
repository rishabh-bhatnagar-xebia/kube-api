package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type PodResponse struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	UID       string            `json:"uid"`
	Labels    map[string]string `json:"labels"`
}

func CreatePod(namespace, imageName string) error {
	clientset, err := getClientSet()
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().Pods(namespace).Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateContainerName(imageName),
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  generateContainerName(imageName),
					Image: imageName,
				},
			},
		},
		Status: corev1.PodStatus{},
	}, metav1.CreateOptions{})
	return err
}

func ListPods(namespace string) ([]corev1.Pod, error) {
	clientset, err := getClientSet()
	if err != nil {
		return nil, err
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

func generateContainerName(imageName string) string {
	return fmt.Sprintf("%s-%s", strings.ReplaceAll(imageName, "/", "-"), uuid.New().String())
}
