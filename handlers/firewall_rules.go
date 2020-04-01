package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adeo/iwc-gcp-firewall-api/helpers"
	"github.com/adeo/iwc-gcp-firewall-api/models"
	"github.com/adeo/iwc-gcp-firewall-api/services"
	"github.com/sirupsen/logrus"
	compute "google.golang.org/api/compute/v1"
)

var (
	manager      models.FirewallRuleManager
	googleClient models.GoogleClientInterface
)

func init() {
	manager, _ = models.NewFirewallRuleClient()
	googleClient, _ = models.NewGoogleClient()
}

// ListFirewallRuleHandler returns a set of firewall rules
func ListFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	// Validate needed permissions
	err := validate(r)
	if err != nil {
		handleError(err, w)
		return
	}

	project, serviceProject, application, _ := helpers.GetMuxVars(r)
	applicationRule, err := services.ListFirewallRule(manager, project, serviceProject, application)
	if err != nil {
		handleError(err, w)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		handleError(err, w)
		return
	}

	fmt.Fprint(w, string(res))
}

// GetFirewallRuleHandler return mathing firewall rule
func GetFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(r)
	if err != nil {
		handleError(err, w)
		return
	}

	project, serviceProject, application, rule := helpers.GetMuxVars(r)
	applicationRule, err := services.GetFirewallRule(manager, project, serviceProject, application, rule)
	if err != nil {
		handleError(err, w)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		handleError(err, w)
		return
	}

	fmt.Fprint(w, string(res))
}

// CreateFirewallRuleHandler create a given rule
func CreateFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	// Decode given rule in order to create it
	var body compute.Firewall
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		handleError(models.NewBadRequestError("Fail to decode body"), w)
		return
	}

	// Validate needed permissions
	err = validate(r)
	if err != nil {
		handleError(err, w)
		return
	}

	project, serviceProject, application, rule := helpers.GetMuxVars(r)
	applicationRule, err := services.CreateFirewallRule(manager, project, serviceProject, application, rule, body)
	if err != nil {
		handleError(err, w)
		return
	}

	res, err := json.Marshal(applicationRule)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(res))
}

// DeleteFirewallRuleHandler delete the given firewall rule
func DeleteFirewallRuleHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(r)
	if err != nil {
		handleError(err, w)
		return
	}

	project, serviceProject, application, rule := helpers.GetMuxVars(r)
	err = services.DeleteFirewallRule(manager, project, serviceProject, application, rule)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// The function valid if
// - provided Bearer token is okay
// - provided service project is a host project's service project
// - consumer is owner of the service project
func validate(r *http.Request) error {
	project, serviceProject, _, _ := helpers.GetMuxVars(r)

	user, err := services.GetUserEmailFromJWT(r.Header.Get("Authorization"))
	if err != nil {
		return err
	}

	// Test owner rights
	err = googleClient.IsProjectOwner(user, serviceProject)
	if err != nil {
		return err
	}

	// Test if service project/project
	err = googleClient.IsAServiceProjectOf(serviceProject, project)
	if err != nil {
		return err
	}

	return nil
}

func handleError(err error, w http.ResponseWriter) {
	if v, ok := err.(*models.ApplicationError); ok {
		w.WriteHeader(v.Code)
		fmt.Fprint(w, v.Error())
		return
	}

	// Else, throw error 500
	logrus.WithFields(logrus.Fields{
		"go-err": err.Error(),
	}).Error("Unexpected error")
	handleError(models.NewInternalError(), w)
}
