package main

import (
	"encoding/json"
	"fmt"

	sentry "github.com/atlassian/go-sentry-api"
)

func discoveryProjects(sentryOrg string, projects []sentry.Project) error {
	discoveryData := make(map[string][]map[string]string)
	var discoveredItems []map[string]string

	for _, project := range projects {
		discoveredItem := make(map[string]string)
		discoveredItem["{#PROJECT}"] = project.Name
		//discoveredItem["{#SLUG}"] = strings.ToLower(project.CallSign)
		discoveredItems = append(discoveredItems, discoveredItem)
	}

	discoveredItem := make(map[string]string)
	discoveredItem["{#PROJECT}"] = sentryOrg
	discoveredItems = append(discoveredItems, discoveredItem)

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
