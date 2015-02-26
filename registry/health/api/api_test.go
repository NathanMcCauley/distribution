package api

import (
	"github.com/docker/distribution/registry/health"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGETDownHandlerDoesNotChangeStatus ensures that calling the endpoint
// /debug/health/down with METHOD GET returns a 404
func TestGETDownHandlerDoesNotChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health/down", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	DownHandler(recorder, req)

	if recorder.Code != 404 {
		t.Errorf("Did not get a 404.")
	}
}

// TestGETUpHandlerDoesNotChangeStatus ensures that calling the endpoint
// /debug/health/down with METHOD GET returns a 404
func TestGETUpHandlerDoesNotChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health/up", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	DownHandler(recorder, req)

	if recorder.Code != 404 {
		t.Errorf("Did not get a 404.")
	}
}

// TestPOSTUpHandlerChangeStatus ensures the endpoint /debug/health/up changes
// the status of the check to StatusOK
func TestPOSTUpHandlerChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "https://fakeurl.com/debug/health/up", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	health.UpdateStatus(health.HealthStatus{"manual_status", health.StatusError})
	UpHandler(recorder, req)

	if health.CheckStatus() != health.StatusOK {
		t.Errorf("Did not get a StatusOK.")
	}
}

// TestPOSTDownHandlerChangeStatus ensures the endpoint /debug/health/down changes
// the status of the check to StatusError
func TestPOSTDownHandlerChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "https://fakeurl.com/debug/health/down", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	health.UpdateStatus(health.HealthStatus{"manual_status", health.StatusOK})
	DownHandler(recorder, req)

	if health.CheckStatus() != health.StatusError {
		t.Errorf("Did not get StatusError.")
	}
}
