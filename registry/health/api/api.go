package api

import (
	"github.com/docker/distribution/registry/health"
	"net/http"
)

func DownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		health.UpdateStatus(health.HealthStatus{"manual_status", health.StatusError})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func UpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		health.UpdateStatus(health.HealthStatus{"manual_status", health.StatusOK})
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func init() {
	http.HandleFunc("/debug/health/down", DownHandler)
	http.HandleFunc("/debug/health/up", UpHandler)
}
