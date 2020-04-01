package handlers

import (
	"fmt"
	"net/http"

	"github.com/adeo/iwc-gcp-firewall-api/models"
)

// NotFoundHandler return a 404 error response
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	e := models.NewNotFoundError()
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(e.Code)
	fmt.Fprint(w, e.Error())
}
