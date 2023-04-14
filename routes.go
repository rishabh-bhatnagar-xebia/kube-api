package main

import "net/http"

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/create", s.handleCreatePod)
	s.router.HandleFunc("/list", s.handleListPods)
	s.router.HandleFunc("/logs", s.handleGetLogs)
}

func (s *Server) Serve(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
