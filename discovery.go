package main

import (
	"encoding/json"
	"fmt"

	sentry "github.com/atlassian/go-sentry-api"
)

func discoveryOrganizations(
	organizations []sentry.Organization,
) error {
	discoveryData := make(map[string][]map[string]string)
	var discoveredItems []map[string]string

	for _, organization := range organizations {
		discoveredItem := make(map[string]string)
		discoveredItem["{#ORGANIZATION}"] = organization.Name
		discoveredItems = append(discoveredItems, discoveredItem)
	}

	discoveryData["data"] = discoveredItems

	out, err := json.Marshal(discoveryData)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", out)
	return nil
}

func discoveryProjects(
	projects []sentry.Project,
) error {
	discoveryData := make(map[string][]map[string]string)
	var discoveredItems []map[string]string

	for _, project := range projects {
		discoveredItem := make(map[string]string)
		discoveredItem["{#PROJECT}"] = project.Name
		discoveredItem["{#ORGANIZATION}"] = project.Organization.Name
		discoveredItems = append(discoveredItems, discoveredItem)
	}

	discoveryData["data"] = discoveredItems

	out, err := json.Marshal(discoveryData)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", out)
	return nil
}
