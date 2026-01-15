package lookup

import (
	"encoding/json"
	"net/http"
)

// Handler exposes the lookup HTTP endpoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	if !enforceAccess(w, r) {
		return
	}

	msisdn := r.URL.Query().Get("msisdn")
	if msisdn == "" {
		http.Error(w, "missing msisdn parameter", http.StatusBadRequest)
		return
	}

	resp := Analyze(msisdn)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
