package models

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// FirewallRule descibe a firewall rule
type FirewallRule struct {
	Rule       compute.Firewall `json:"item"`
	CustomName string           `json:"custom_name"`
}

// FirewallRules describe a set of firewall rule
type FirewallRules []FirewallRule

// ApplicationRule describe and end-user response
type ApplicationRule struct {
	Project        string        `json:"project"`
	ServiceProject string        `json:"service_project"`
	Application    string        `json:"application"`
	Rules          FirewallRules `json:"data"`
}

// FirewallRuleManager contains methods to manage firewall rules
type FirewallRuleManager interface {
	ListFirewallRule(project string) ([]*compute.Firewall, error)
	GetFirewallRule(project, name string) (*compute.Firewall, error)
	CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error)
	DeleteFirewallRule(project, name string) error
}

// FirewallRuleClient provides primitives to collect rules from Google Cloud Platform. Implements FirewallRuleManager
type FirewallRuleClient struct {
	computeService *compute.Service
}

// NewFirewallRuleClient FirewallRuleClient contructor
func NewFirewallRuleClient() (*FirewallRuleClient, error) {
	c, err := google.DefaultClient(context.Background(), compute.CloudPlatformScope)
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	computeService, err := compute.New(c)
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	manager := FirewallRuleClient{}
	manager.computeService = computeService
	return &manager, nil
}

// ListFirewallRule returns given project's firewall rule
func (f *FirewallRuleClient) ListFirewallRule(project string) ([]*compute.Firewall, error) {
	req := f.computeService.Firewalls.List(project)

	var firewallRuleList []*compute.Firewall

	err := req.Pages(context.Background(), func(page *compute.FirewallList) error {
		for _, firewall := range page.Items {
			firewallRuleList = append(firewallRuleList, firewall)
		}
		return nil
	})
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	return firewallRuleList, nil
}

// GetFirewallRule returns firewall rule matching given project and name
func (f *FirewallRuleClient) GetFirewallRule(project, name string) (*compute.Firewall, error) {
	rules, err := f.computeService.Firewalls.Get(project, name).Context(context.Background()).Do()
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}
	return rules, nil
}

// CreateFirewallRule create given firewall rule on given project
func (f *FirewallRuleClient) CreateFirewallRule(project string, rule *compute.Firewall) (*compute.Firewall, error) {
	_, err := f.computeService.Firewalls.Insert(project, rule).Context(context.Background()).Do()
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	return f.GetFirewallRule(project, rule.Name)
}

// DeleteFirewallRule delete firewall rule matching given project and name
func (f *FirewallRuleClient) DeleteFirewallRule(project string, name string) error {
	_, err := f.computeService.Firewalls.Delete(project, name).Context(context.Background()).Do()
	if e, ok := err.(*googleapi.Error); ok {
		return NewGoogleApplicationError(e)
	}
	return err
}
