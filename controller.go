package main

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleCreatePod(w http.ResponseWriter, r *http.Request) {
	var createParams CreatePodRequest
	err := json.NewDecoder(r.Body).Decode(&createParams)
	if err != nil {
		http.Error(w, Wrap(err), http.StatusBadRequest)
		return
	}

	err = s.CreatePod(getNamespace(createParams.Namespace), createParams.Image, createParams.PodName)
	if err != nil {
		http.Error(w, Wrap(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ns := r.URL.Query().Get(HTTP_PARAM_NAMESPACE)
	pods, err := s.ListPods(getNamespace(ns))
	if err != nil {
		http.Error(w, Wrap(err), http.StatusInternalServerError)
		return
	}

	jsonPods, err := json.Marshal(FilterPodFields(pods))
	if err != nil {
		http.Error(w, Wrap(err), http.StatusInternalServerError)
		return
	}

	w.Write(jsonPods)
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	var getLogsParams GetLogsRequest
	err := json.NewDecoder(r.Body).Decode(&getLogsParams)
	if err != nil {
		http.Error(w, Wrap(err), http.StatusBadRequest)
		return
	}

	err = s.StreamLogs(w, getNamespace(getLogsParams.Namespace), getLogsParams.PodName)
	if err != nil {
		http.Error(w, Wrap(err), http.StatusInternalServerError)
		return
	}
}
