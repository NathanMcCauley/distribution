package health

import (
  "net/http"
  "encoding/json"
)

var statusKeys = make(map[string]HealthStatus)

type Status string
const (
  OK = "Ok"
  WARNING = "Warning"
  ERROR = "Error"
)

type HealthStatus struct {
  Name string
  CurrentStatus Status
}

func UpdateStatus(status HealthStatus) {
  statusKeys[status.Name] = status
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  jsonResponse, _ := json.Marshal(statusKeys)
  w.Write(jsonResponse)
}

func downHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    UpdateStatus(HealthStatus{"manual_status", ERROR})
  }
}

func upHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {
    UpdateStatus(HealthStatus{"manual_status", OK})
  }
}

func init() {
  http.HandleFunc("/debug/health", statusHandler)
  http.HandleFunc("/debug/health/down", downHandler)
  http.HandleFunc("/debug/health/up", upHandler)
}
