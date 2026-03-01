package gtwhttp

import (
	"fmt"
	"net/http"

	"github.com/martinsdevv/aegis/internal/gateway/middleware"
)

type AdminHandler struct {
	Store *middleware.APIKeyStore
}

func NewAdminHandler(store *middleware.APIKeyStore) *AdminHandler {
	return &AdminHandler{Store: store}
}

// DELETE /admin/cache/apikey/{hash}
func (a *AdminHandler) InvalidateAPIKey(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "missing hash", http.StatusBadRequest)
		return
	}

	err := a.Store.DeleteFromCache(r.Context(), hash)
	if err != nil {
		http.Error(w, "failed to delete cache", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleNilPointer(w http.ResponseWriter, r *http.Request) {
	var x *int
	fmt.Println(*x)
}

func HandleRLTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true,"route":"rltest"}`))
}
