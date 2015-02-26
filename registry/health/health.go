package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	mutex      sync.RWMutex
	statusKeys = make(map[string]HealthStatus)
)

// Status type allows us to use semantically meaningful Codes.
type Status string

// Represents the possible server states based on the currently recorded
// healthchecks.
const (
	StatusOK      = "StatusOK"
	StatusWarning = "StatusWarning"
	StatusError   = "StatusError"
)

// HealthStatus represents a named status check and it's current status.
type HealthStatus struct {
	Name          string
	CurrentStatus Status
}

// UpdateStatus updates the status of a status check
func UpdateStatus(status HealthStatus) {
	mutex.Lock()
	defer mutex.Unlock()
	statusKeys[status.Name] = status
}

// CheckStatus returns the status of the worst of all the currently registered
// health checks.
// StatusError < StatusWarning < StatusOK
func CheckStatus() Status {
	warning := false
	mutex.RLock()
	defer mutex.RUnlock()
	for _, v := range statusKeys {
		switch v.CurrentStatus {
		case StatusError:
			return StatusError
		case StatusWarning:
			warning = true
		}
	}
	if warning {
		return StatusWarning
	}
	return StatusOK
}

// StatusHandler returns a JSON blob with all the currently registered Health Checks
// and their corresponding status.
// Returns 500 if any Error status exists, 200 otherwise
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if CheckStatus() == StatusError {
			w.WriteHeader(http.StatusInternalServerError)
		}
		jsonResponse, _ := json.Marshal(statusKeys)
		w.Write(jsonResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Registers global /debug/health api endpoint
func init() {
	http.HandleFunc("/debug/health", StatusHandler)
}
