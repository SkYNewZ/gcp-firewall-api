package services

import (
	"fmt"
	"strings"

	"github.com/adeo/iwc-gcp-firewall-api/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
)

// ListFirewallRule returns a set of firewall rules related to an application
func ListFirewallRule(manager models.FirewallRuleManager, project, serviceProject, application string) (*models.ApplicationRule, error) {
	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
	}).Debugln("Listing rules")

	// List all firewall rule in given project
	gRules, err := manager.ListFirewallRule(project)
	if err != nil {
		return nil, err
	}

	// Create reponse
	endUserResult := models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
	}

	endUserResultRules := make(models.FirewallRules, 0)

	// For each obtains Google rules
	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	for _, gRule := range gRules {
		// Filter with managed rules with this application
		if strings.HasPrefix(gRule.Name, prefix) {
			customName := gRule.Name[len(prefix):]
			endUserResultRules = append(endUserResultRules, models.FirewallRule{
				Rule:       *gRule,
				CustomName: customName,
			})
		}
	}

	endUserResult.Rules = endUserResultRules
	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
	}).Debugf("Found %d rules", len(endUserResultRules))
	return &endUserResult, nil
}

// CreateFirewallRule create given firewall rule on given project
func CreateFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string, rule compute.Firewall) (*models.ApplicationRule, error) {
	customNameAndTargetTag := fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)
	rule.Name = customNameAndTargetTag
	rule.TargetTags = []string{customNameAndTargetTag}

	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
		"rule_name":       customNameAndTargetTag,
		"target_tag":      customNameAndTargetTag,
	}).Debugln("Creating rule")

	gRule, err := manager.CreateFirewallRule(project, &rule)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	createdRule := models.FirewallRule{
		Rule:       *gRule,
		CustomName: gRule.Name[len(prefix):],
	}

	return &models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
		Rules:          models.FirewallRules{createdRule},
	}, nil
}

// GetFirewallRule return matching firewall rule
func GetFirewallRule(manager models.FirewallRuleManager, project string, serviceProject string, application string, ruleName string) (*models.ApplicationRule, error) {
	n := fmt.Sprintf("%s-%s-%s", serviceProject, application, ruleName)

	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
		"rule_name":       n,
	}).Debugln("Searching rule")

	gRule, err := manager.GetFirewallRule(project, n)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
		"rule_name":       n,
	}).Debugf("Rule found with ID %d", gRule.Id)

	prefix := fmt.Sprintf("%s-%s-", serviceProject, application)
	createdRule := models.FirewallRule{
		Rule:       *gRule,
		CustomName: gRule.Name[len(prefix):],
	}

	return &models.ApplicationRule{
		Application:    application,
		Project:        project,
		ServiceProject: serviceProject,
		Rules:          models.FirewallRules{createdRule},
	}, nil
}

// DeleteFirewallRule delete firewall rule mathing project, service project, application name and rule name
func DeleteFirewallRule(manager models.FirewallRuleManager, project, serviceProject, application, customName string) error {
	ruleName := fmt.Sprintf("%s-%s-%s", serviceProject, application, customName)
	logrus.WithFields(logrus.Fields{
		"project":         project,
		"service_project": serviceProject,
		"application":     application,
		"rule_name":       ruleName,
	}).Debugln("Deleting rule")

	return manager.DeleteFirewallRule(project, ruleName)
}
