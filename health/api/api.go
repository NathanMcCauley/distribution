package api

import (
	"errors"
	"github.com/docker/distribution/health"
	"net/http"
)

// UpHandler registers a manual_http_status that always returns an Error
func DownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		health.Register("manual_http_status", health.CheckFunc(func() error {
			return errors.New("Manual Check")
		}))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// UpHandler registers a manual_http_status that always returns nil
func UpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		health.Register("manual_http_status", health.CheckFunc(func() error {
			return nil
		}))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// init sets up the two endpoints to bring the service up and down
func init() {
	http.HandleFunc("/debug/health/down", DownHandler)
	http.HandleFunc("/debug/health/up", UpHandler)
}
