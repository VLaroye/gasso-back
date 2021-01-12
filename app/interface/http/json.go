package http

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error encoding response")
	}
}
