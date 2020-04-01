package models

import (
	"context"
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GoogleClient describe some operations with Google
type GoogleClient struct {
	projectService *cloudresourcemanager.ProjectsService
	computeService *compute.ProjectsService
}

// GoogleClientInterface describe GoogleClient's operations
type GoogleClientInterface interface {
	IsProjectOwner(user string, projectID string) error
	IsAServiceProjectOf(projectA, projetB string) error
}

// NewGoogleClient GoogleClient constructor
func NewGoogleClient() (*GoogleClient, error) {
	p, err := cloudresourcemanager.NewService(context.Background(), option.WithScopes(cloudresourcemanager.CloudPlatformReadOnlyScope))
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	c, err := compute.NewService(context.Background(), option.WithScopes(cloudresourcemanager.CloudPlatformScope))
	if e, ok := err.(*googleapi.Error); ok {
		return nil, NewGoogleApplicationError(e)
	}

	return &GoogleClient{
		projectService: p.Projects,
		computeService: c.Projects,
	}, nil
}

// IsProjectOwner return if a given user is owner of the given project
// Return nil if user is owner
// https://cloud.google.com/resource-manager/reference/rest/v1/projects/getIamPolicy
func (c *GoogleClient) IsProjectOwner(user string, projectID string) error {
	// Get project's policy
	policy, err := c.projectService.GetIamPolicy(projectID, &cloudresourcemanager.GetIamPolicyRequest{}).Context(context.Background()).Do()
	if e, ok := err.(*googleapi.Error); ok {
		return NewGoogleApplicationError(e)
	}

	user = "user:" + user

	// Browse each binding and return if we have a match for roles/owner
	for _, b := range policy.Bindings {
		if b.Role == "roles/owner" {
			for _, member := range b.Members {
				if member == user {
					return nil
				}
			}
		}
	}

	return NewForbiddenError(fmt.Sprintf("User [%s] does not have permission to access project [%s]. The resource may not exist or you don't have roles/owner", user, projectID))
}

// IsAServiceProjectOf test if given projectA is a service project of projectB
// Return nil if validated
// https://cloud.google.com/compute/docs/reference/rest/v1/projects/getXpnHost
func (c *GoogleClient) IsAServiceProjectOf(projectA, projetB string) error {
	projectAHostProject, err := c.computeService.GetXpnHost(projectA).Context(context.Background()).Do()
	if e, ok := err.(*googleapi.Error); ok {
		return NewGoogleApplicationError(e)
	}

	e := NewForbiddenError(fmt.Sprintf("Project [%s] is not a [%s]'s service project or it may not exist", projectA, projetB))

	// Not host project found for projectA, this is not a service project
	if projectAHostProject == nil {
		return e
	}

	// If project math
	if projectAHostProject.Name == projetB {
		return nil
	}

	return e
}
