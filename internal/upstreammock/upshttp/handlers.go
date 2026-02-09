package upshttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type pongResponse struct {
	Pong bool `json:"pong"`
}

func HandlePing(w http.ResponseWriter, _ *http.Request) {
	pong := true
	if !pong {
		http.Error(w, "something got really wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := pongResponse{
		Pong: pong,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
	// fmt.Printf("%d bytes written\n", wc)
}

func HandleEcho(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed, only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.ContentLength == 0 {
		http.Error(w, "a body is required", http.StatusUnprocessableEntity)
		return
	}

	var pld map[string]any

	err := json.NewDecoder(r.Body).Decode(&pld)
	if err != nil {
		slog.Error("invalid json body", "err", err)
		http.Error(w, "invalid json body", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]any{
		"echo": pld,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}
