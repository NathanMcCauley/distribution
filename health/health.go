package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	mutex            sync.RWMutex
	registeredChecks = make(map[string]Checker)
)

// Checker is the interface for a Health Checker
type Checker interface {
	// Check returns nil if the service is okay.
	Check() error
}

// CheckFunc is a convenience type to create functions that implement
// the Checker interface
type CheckFunc func() error

// Implements the Checker interface to allow for any func() error method
// to be passed as a Checker
func (cf CheckFunc) Check() error {
	return cf()
}

// Updater implements a health check that is explicitly set.
type Updater interface {
	Checker

	// Update updates the current status of the health check.
	Update(status error)
}

// NewStatusUpdater returns an explicitly settable status check.
func NewStatusUpdater() Updater {
	return &updater{}
}

// CheckStatus returns a map with all the current health check errors
func CheckStatus() map[string]string {
	mutex.RLock()
	defer mutex.RUnlock()
	statusKeys := make(map[string]string)
	for k, v := range registeredChecks {
		err := v.Check()
		if err != nil {
			statusKeys[k] = err.Error()
		}
	}

	return statusKeys
}

// Updater implements Checker and provides an asynchronous Update method.
// This allows us to have a Checker that returns the Check() call immediately
// not blocking on a potentially expensive check.
type updater struct {
	mu     sync.Mutex
	status error
}

// Implements the Checker interface
func (u *updater) Check() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.status
}

// Update allows asynchronous access to the status of a Checker.
func (u *updater) Update(status error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.status = status
}

// PeriodicChecker wraps an updater to provide a periodic checker
func PeriodicChecker(check Checker, period time.Duration) Checker {
	u := NewStatusUpdater()
	go func() {
		t := time.NewTicker(period)
		for {
			<-t.C
			u.Update(check.Check())
		}
	}()

	return u
}

// Register associates the checker with the provided name. We allow
// overwrites to a specific check status.
func Register(name string, check Checker) {
	mutex.RLock()
	defer mutex.RUnlock()

	registeredChecks[name] = check
}

// RegisterFunc allows the convenience of registering a checker directly
// from an arbitrary func() error
func RegisterFunc(name string, check func() error) {
	Register(name, CheckFunc(check))
}

// RegisterPeriodicFunc allows the convenience of registering a PeriodicChecker
// from an arbitrary func() error
func RegisterPeriodicFunc(name string, check func() error, period time.Duration) {
	Register(name, PeriodicChecker(CheckFunc(check), period))
}

// StatusHandler returns a JSON blob with all the currently registered Health Checks
// and their corresponding status.
// Returns 500 if any Error status exists, 200 otherwise
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		checksStatus := CheckStatus()
		if len(checksStatus) != 0 {
			w.WriteHeader(http.StatusInternalServerError)
		}
		jsonResponse, _ := json.Marshal(checksStatus)
		w.Write(jsonResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Registers global /debug/health api endpoint
func init() {
	http.HandleFunc("/debug/health", StatusHandler)
}
