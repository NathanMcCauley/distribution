package health

import (
	"errors"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, recorder.Code, 200, "Code should be 200")
}

// TestReturns500IfThereAreErrorChecks ensures that the result code of the
// health endpoint is 500 if there are health checks with errors
func TestReturns500IfThereAreErrorChecks(t *testing.T) {
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "https://fakeurl.com/debug/health", nil)
	if err != nil {
		t.Errorf("Failed to create request.")
	}

	// Create a manual error
	Register("some_check", CheckFunc(func() error {
		return errors.New("This Check did not succeed")
	}))

	StatusHandler(recorder, req)

	assert.Equal(t, recorder.Code, 500, "Code should be 500")
}
