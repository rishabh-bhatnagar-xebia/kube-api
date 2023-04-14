package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func HandleCreatePod(w http.ResponseWriter, r *http.Request) {
	imageName, err := readImageName(r)
	if err != nil {
		http.Error(w, wrap(err), http.StatusBadRequest)
		return
	}
	namespace := readNamespace(r)
	log.Println("creating the pod in the default namespace")

	err = CreatePod(namespace, imageName)
	if err != nil {
		http.Error(w, wrap(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleListPods(w http.ResponseWriter, r *http.Request) {
	setJson(w)

	pods, err := ListPods(readNamespace(r))
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

func readImageName(r *http.Request) (string, error) {
	if r.URL.Query().Has(HTTP_PARAM_IMAGE) {
		return r.URL.Query().Get(HTTP_PARAM_IMAGE), nil
	}
	return "", errors.New("missing image name")
}
