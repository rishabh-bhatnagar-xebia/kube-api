package main

import (
	"encoding/json"
	"net/http"
)

func HandleCreatePod(w http.ResponseWriter, r *http.Request) {
}

func HandleListPods(w http.ResponseWriter, r *http.Request) {
	setJson(w)

	pods, err := ListPods(readNamespace(r))
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
	}

	jsonPods, err := json.Marshal(ExtractFields(pods))
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
	}

	w.Write(jsonPods)
}

func HandleGetLogs(w http.ResponseWriter, r *http.Request) {
}

func setJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func readNamespace(r *http.Request) string {
	if r.URL.Query().Has(HTTP_PARAM_NAMESPACE) {
		return r.URL.Query().Get(HTTP_PARAM_NAMESPACE)
	}
	return "default"
}
