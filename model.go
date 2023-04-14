package main

import (
	"context"
	"errors"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) CreatePod(namespace, imageName, podName string) error {
	if len(imageName) == 0 {
		return errors.New("image_name cannot be empty")
	}

	if len(podName) == 0 {
		return errors.New("pod_name cannot be empty")
	}

	_, err := s.clientset.CoreV1().Pods(namespace).Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  podName,
					Image: imageName,
				},
			},
		},
		Status: corev1.PodStatus{},
	}, metav1.CreateOptions{})
	return err
}

func (s *Server) ListPods(namespace string) ([]corev1.Pod, error) {
	// List all pods in the default namespace
	pods, err := s.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

func (s *Server) StreamLogs(w http.ResponseWriter, namespace, podName string) error {
	if len(podName) == 0 {
		return errors.New("pod_name cannot be empty")
	}

	// create a request for the logs of the specified container
	req := s.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
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
}
