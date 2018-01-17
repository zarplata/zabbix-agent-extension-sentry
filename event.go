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
	sentryOrg string,
	sentryToken string,
	discovery bool,
	hostname string,
	zabbix string,
	port int,
	zabbixPrefix string,
) error {

	client, err := sentry.NewClient(sentryToken, &sentryApi, nil)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't connect sentry %s with token %s.",
			sentryApi,
			sentryToken,
		)
	}

	projects, err := client.GetProjects()
	if err != nil {
		return hierr.Errorf(err, "can't fetch all projects from sentry.")
	}

	if discovery {
		if err := discoveryProjects(sentryOrg, projects); err != nil {
			return hierr.Errorf(err, "can't discovery projects.")
		}
		os.Exit(0)
	}

	organization, err := client.GetOrganization(sentryOrg)
	if err != nil {
		return hierr.Errorf(err, "can't fetch organization.")
	}

	var metrics []*zsend.Metric

	// 59 seconds interval, 6 time series value retrive
	now := (time.Now().Unix()/10)*10 - 10
	preview := time.Duration(59) * time.Second
	later := now - int64(preview.Seconds())

	organizationStats, err := client.GetOrganizationStats(
		organization,
		"received",
		later,
		now,
		nil)
	if err != nil {
		return hierr.Errorf(err, "can't get organization stats.")
	}
	metrics = createMetrics(
		hostname,
		sentryOrg,
		metrics,
		organizationStats,
		zabbixPrefix,
	)
	for _, project := range projects {
		projectStats, err := client.GetProjectStats(
			organization,
			project,
			"received",
			later,
			now,
			nil)
		if err != nil {
			return hierr.Errorf(err, "can't get project	stats.")
		}
		metrics = createMetrics(
			hostname,
			project.Name,
			metrics,
			projectStats,
			zabbixPrefix,
		)
	}
	packet := zsend.NewPacket(metrics)
	sender := zsend.NewSender(
		zabbix,
		port,
	)
	sender.Send(packet)
	return nil
}
