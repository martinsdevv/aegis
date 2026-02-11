package gtwhttp

import (
	"fmt"
	"net/http"
)

func HandleNilPointer(w http.ResponseWriter, r *http.Request) {
	var x *int
	fmt.Println(*x)
}
