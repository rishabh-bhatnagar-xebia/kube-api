package main

import "net/http"

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("/create", HandleCreatePod)
	router.HandleFunc("/list", HandleListPods)
	router.HandleFunc("/logs", HandleGetLogs)
}
