package health

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	HealthOK bool `json:"ok"`
}

func HealthHandler(c *Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed, only GET is available", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if c.Ready() {
			resp := healthResponse{
				HealthOK: true,
			}
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "error writing response", http.StatusInternalServerError)
				return
			}
			return
		}

		resp := healthResponse{
			HealthOK: false,
		}

		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "error writing response", http.StatusInternalServerError)
			return
		}
		http.Error(w, "not ready", http.StatusServiceUnavailable)
	}
}
