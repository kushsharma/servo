package routes

import (
	"encoding/json"
	"net/http"

	"github.com/kushsharma/servo/internal"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(internal.AppStats)
}
