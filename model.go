package main

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
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

func (s *Server) StreamLogs(w http.ResponseWriter, r *http.Request, namespace, podName string) error {
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
	return relayStreamUsingWebSocket(stream, w, r)
}

func relayStreamUsingWebSocket(stream io.ReadCloser, w http.ResponseWriter, r *http.Request) error {
	// upgrade HTTP request to WebSocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  WEBSOCKET_BUFFER_SIZE,
		WriteBufferSize: WEBSOCKET_BUFFER_SIZE,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.New("Failed to upgrade connection: " + err.Error())
	}
	defer conn.Close()

	// stream the logs from kubernetes to client
	buf := make([]byte, WEBSOCKET_BUFFER_SIZE)
	for {
		// read from the kubernetes logs
		n, err := stream.Read(buf)
		if err != nil {
			return err
		}

		// write to the client's websocket
		err = conn.WriteMessage(websocket.TextMessage, buf[:n])
		if err != nil {
			return err
		}
	}
}
