package health

import (
  "net/http"
  "encoding/json"
  "sync"
)

var (
  mutex sync.RWMutex
  statusKeys = make(map[string]HealthStatus)
)

type Status string
const (
  StatusOK = "StatusOK"
  StatusWarning = "StatusWarning"
  StatusError = "StatusError"
)

type HealthStatus struct {
  Name string
  CurrentStatus Status
}

func UpdateStatus(status HealthStatus) {
  mutex.Lock()
  defer mutex.Unlock()
  statusKeys[status.Name] = status
}

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

func statusHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  if CheckStatus() == StatusError {
    w.WriteHeader(http.StatusInternalServerError)
  }
  jsonResponse, _ := json.Marshal(statusKeys)
  w.Write(jsonResponse)
}

func downHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    UpdateStatus(HealthStatus{"manual_status", StatusError})
  }
}

func upHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    UpdateStatus(HealthStatus{"manual_status", StatusOK})
  }
}

func init() {
  http.HandleFunc("/debug/health", statusHandler)
  http.HandleFunc("/debug/health/down", downHandler)
  http.HandleFunc("/debug/health/up", upHandler)
}
