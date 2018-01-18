package main

import (
	"os"
	"time"

	sentry "github.com/atlassian/go-sentry-api"
	zsend "github.com/blacked/go-zabbix"
	hierr "github.com/reconquest/hierr-go"
)

func event(
	sentryApi string,
	sentryToken string,
	discovery bool,
	hostname string,
	zabbix string,
	port int,
	zabbixPrefix string,
) error {

	var statTypes = []sentry.StatQuery{
		"received",
		"rejected",
		"blacklisted",
	}

	client, err := sentry.NewClient(sentryToken, &sentryApi, nil)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't connect sentry %s with token %s.",
			sentryApi,
			sentryToken,
		)
	}

	organizations, _, err := client.GetOrganizations()
	if err != nil {
		return hierr.Errorf(err, "can't fetch organizations.")
	}

	projects, err := client.GetProjects()
	if err != nil {
		return hierr.Errorf(err, "can't fetch projects.")
	}

	if discovery {
		if err := discoveryOrgsProjects(organizations, projects); err != nil {
			return hierr.Errorf(
				err,
				"can't discovery organizations and projects.",
			)
		}
		os.Exit(0)
	}

	var metrics []*zsend.Metric

	// 59 seconds interval, 6 time series value retrive
	now := (time.Now().Unix()/10)*10 - 10
	preview := time.Duration(59) * time.Second
	later := now - int64(preview.Seconds())

	for _, statType := range statTypes {

		for _, organization := range organizations {

			organizationStats, err := client.GetOrganizationStats(
				organization,
				statType,
				later,
				now,
				nil)
			if err != nil {
				return hierr.Errorf(err, "can't get organization stats.")
			}

			metrics = createOrganizationMetrics(
				hostname,
				organization.Name,
				statType,
				metrics,
				organizationStats,
				zabbixPrefix,
			)

			for _, project := range projects {
				projectStats, err := client.GetProjectStats(
					organization,
					project,
					statType,
					later,
					now,
					nil)
				if err != nil {
					return hierr.Errorf(err, "can't get project	stats.")
				}

				metrics = createProjectMetrics(
					hostname,
					project.Organization.Name,
					project.Name,
					statType,
					metrics,
					projectStats,
					zabbixPrefix,
				)
			}
		}
	}
	packet := zsend.NewPacket(metrics)
	sender := zsend.NewSender(
		zabbix,
		port,
	)
	sender.Send(packet)
	return nil
}
