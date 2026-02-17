package gtwhttp

import (
	"fmt"
	"net/http"
)

func HandleNilPointer(w http.ResponseWriter, r *http.Request) {
	var x *int
	fmt.Println(*x)
}

func HandleRLTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true,"route":"rltest"}`))
}
