package main

import (
	"encoding/json"
	"net/http"
)

func HandleCreatePod(w http.ResponseWriter, r *http.Request) {
	var createParams CreatePodRequest
	err := json.NewDecoder(r.Body).Decode(&createParams)
	if err != nil {
		http.Error(w, wrap(err), http.StatusBadRequest)
		return
	}

	err = CreatePod(getNamespace(createParams.Namespace), createParams.Image)
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleListPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ns := r.URL.Query().Get(HTTP_PARAM_NAMESPACE)
	pods, err := ListPods(getNamespace(ns))
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
		return
	}

	jsonPods, err := json.Marshal(ExtractFields(pods))
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
		return
	}

	w.Write(jsonPods)
}

func HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	var getLogsParams GetLogsRequest
	err := json.NewDecoder(r.Body).Decode(&getLogsParams)
	if err != nil {
		http.Error(w, wrap(err), http.StatusBadRequest)
		return
	}

	err = StreamLogs(w, getLogsParams.Namespace, getLogsParams.PodName)
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
		return
	}
}

// getNamespace returns the input namespace if it is set,
// otherwise returns a default one
func getNamespace(namespace string) string {
	if len(namespace) == 0 {
		return "default"
	}
	return namespace
}
