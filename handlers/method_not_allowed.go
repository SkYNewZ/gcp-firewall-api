package handlers

import (
	"fmt"
	"net/http"

	"github.com/adeo/iwc-gcp-firewall-api/models"
)

// MethodNotAllowedHandler return a 405 error response
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	e := models.NewMethodNotAllowedError()
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	fmt.Fprint(w, e.Error())
}
