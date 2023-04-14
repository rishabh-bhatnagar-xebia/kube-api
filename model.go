package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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
		return nil, err
	}

	return pods.Items, nil
}

func StreamLogs(w http.ResponseWriter, namespace, podName string) error {
	clientset, err := getClientSet()
	if err != nil {
		return err
	}

	// create a request for the logs of the specified container
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow: true,
	})

	// stream the logs
	stream, err := req.Stream(context.Background())
	if err != nil {
		return err
	}
	defer stream.Close()

	// write the logs to the client
	return forwardStream(stream, w)
}

// forwardStream reads the content from stream and forwards it to the
// ResponseWriter at a buffer rate of 1MB per flush
func forwardStream(stream io.ReadCloser, w http.ResponseWriter) error {
	// flush the last fetched data to the request
	f, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming unsupported")
	}

	buf := make([]byte, 1<<10) // todo: magic number can be exported to constants
	for {
		n, err := stream.Read(buf)
		if err != nil {
			return err
		}
		_, err = w.Write(buf[0:n])
		f.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func FilterPodFields(pods []corev1.Pod) (ret []PodResponse) {
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
