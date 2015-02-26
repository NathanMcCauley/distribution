package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestReturns200IfThereAreNoChecks ensures that the result code of the health
// endpoint is 200 if there are not currently registered checks.
func TestReturns200IfThereAreNoChecks(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	StatusHandler(recorder, req)

	if recorder.Code != 200 {
		t.Errorf("Did not get a 200.")
	}
}

// TestReturns200IfThereAreOkOrWarningChecks ensures that the result code of the
// health endpoint is 200 if there are health checks with the StatusWarning code
func TestReturns200IfThereAreOkOrWarningChecks(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	UpdateStatus(HealthStatus{"WarningStatus", StatusWarning})
	UpdateStatus(HealthStatus{"OKStatus", StatusOK})
	StatusHandler(recorder, req)

	if recorder.Code != 200 {
		t.Errorf("Did not get a 200.")
	}
}

// TestReturns500IfThereAreErrorChecks ensures that the result code of the
// health endpoint is 500 if there are health checks with the StatusError code
func TestReturns500IfThereAreErrorChecks(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	UpdateStatus(HealthStatus{"OKStatus", StatusOK})
	UpdateStatus(HealthStatus{"WarningStatus", StatusWarning})
	UpdateStatus(HealthStatus{"ErrorStatus", StatusError})
	StatusHandler(recorder, req)

	if recorder.Code != 500 {
		t.Errorf("Did not get a 500.")
	}
}
