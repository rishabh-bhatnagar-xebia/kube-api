package main

import (
	"encoding/json"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
)

func ResolvePort() string {
	port := PORT
	if len(os.Args) > 1 {
		port = os.Args[1]
		log.Println("got an inline value for the port:", port)
	}
	return port
}

func Wrap(err error) string {
	content, _ := json.Marshal(struct {
		Error string `json: err`
	}{
		Error: err.Error(),
	})
	return string(content)
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
