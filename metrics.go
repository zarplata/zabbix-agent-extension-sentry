package main

import (
	"fmt"
	"strconv"

	sentry "github.com/atlassian/go-sentry-api"
	zsend "github.com/blacked/go-zabbix"
)

func sumEvents(stats []sentry.Stat) int {
	var sum int
	sum = 0
	for _, stat := range stats {
		sum = sum + int(stat[1])
	}
	return sum
}

func makePrefix(prefix, key string) string {
	return fmt.Sprintf(
		"%s.%s", prefix, key,
	)
}

func createOrganizationMetrics(
	hostname string,
	name string,
	statType sentry.StatQuery,
	metrics []*zsend.Metric,
	stats []sentry.Stat,
	prefix string,
) []*zsend.Metric {

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("organization.event.%s.[%s]", statType, name),
			),
			strconv.Itoa(sumEvents(stats)),
		),
	)
	return metrics
}

func createProjectMetrics(
	hostname string,
	organization string,
	name string,
	statType sentry.StatQuery,
	metrics []*zsend.Metric,
	stats []sentry.Stat,
	prefix string,
) []*zsend.Metric {

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf(
					"project.event.%s.[%s][%s]",
					statType,
					organization,
					name),
			),
			strconv.Itoa(sumEvents(stats)),
		),
	)
	return metrics
}

func createQueueMetrics(
	hostname string,
	metrics []*zsend.Metric,
	queueName map[string]string,
	prefix string,
) []*zsend.Metric {

	metrics = append(
		metrics,
		zsend.NewMetric(
			hostname,
			makePrefix(
				prefix,
				fmt.Sprintf("queue.[%s]", queueName["queue"]),
			),
			queueName["event"],
		),
	)
	return metrics
}
