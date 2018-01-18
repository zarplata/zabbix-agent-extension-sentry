package main

import (
	"encoding/json"
	"fmt"

	sentry "github.com/atlassian/go-sentry-api"
)

func discoveryOrgsProjects(
	organizations []sentry.Organization,
	projects []sentry.Project,
) error {
	discoveryData := make(map[string][]map[string]string)
	var discoveredItems []map[string]string

	for _, organization := range organizations {
		discoveredItem := make(map[string]string)
		discoveredItem["{#ORGANIZATION}"] = organization.Name
		discoveredItem["{#TYPE}"] = "organization"
		discoveredItems = append(discoveredItems, discoveredItem)
	}

	for _, project := range projects {
		discoveredItem := make(map[string]string)
		discoveredItem["{#PROJECT}"] = project.Name
		discoveredItem["{#ORGANIZATION}"] = project.Organization.Name
		discoveredItem["{#TYPE}"] = "project"
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

func discoveryQueue() error {
	discoveryData := make(map[string][]map[string]string)
	var discoveredItems []map[string]string

	queueNames, err := getQueueEvent()
	if err != nil {
		return err
	}

	for _, queueName := range queueNames {
		discoveredItem := make(map[string]string)
		discoveredItem["{#QUEUE}"] = queueName["queue"]
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
