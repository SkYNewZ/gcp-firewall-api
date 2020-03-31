package handlers

import (
	"fmt"
	"net/http"

	"github.com/adeo/iwc-gcp-firewall-api/models"
)

// MethodNotFound return a 404 error response
func MethodNotFound(w http.ResponseWriter, r *http.Request) {
	e := models.RouterError{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	fmt.Fprint(w, e.Error())
}
