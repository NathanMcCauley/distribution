package api

import (
	"github.com/docker/distribution/health"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, recorder.Code, 404, "Code should be 404")
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

	assert.Equal(t, recorder.Code, 404, "Code should be 404")
}

// TestPOSTDownHandlerChangeStatus ensures the endpoint /debug/health/down changes
// the status code of the response to 500
// This test is order dependent, and should come before TestPOSTUpHandlerChangeStatus
func TestPOSTDownHandlerChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "https://fakeurl.com/debug/health/down", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	DownHandler(recorder, req)

	assert.Equal(t, recorder.Code, 200, "Code should be 200")
	assert.Equal(t, len(health.CheckStatus()), 1, "Calling downhandler should make health.CheckStatus() return an error check")
}

// TestPOSTUpHandlerChangeStatus ensures the endpoint /debug/health/up changes
// the status code of the response to 200
func TestPOSTUpHandlerChangeStatus(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "https://fakeurl.com/debug/health/up", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	UpHandler(recorder, req)

	assert.Equal(t, recorder.Code, 200, "Code should be 200")
	assert.Equal(t, len(health.CheckStatus()), 0, "Calling uphandler should make health.CheckStatus() return no error check")
}
