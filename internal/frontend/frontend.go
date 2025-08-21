// Package frontend provides embedded web dashboard for GHbex
package frontend

import (
	_ "embed"
	"net/http"
)

// Embedded dashboard HTML

//go:embed dashboard.html
var DashboardHTML []byte

// ServeDashboard serves the embedded dashboard
func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.WriteHeader(http.StatusOK)
	w.Write(DashboardHTML)
}
