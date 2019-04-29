package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	sentry "github.com/atlassian/go-sentry-api"
	zsend "github.com/blacked/go-zabbix"
	docopt "github.com/docopt/docopt-go"
	hierr "github.com/reconquest/hierr-go"
)

var (
	version   = "[manual build]"
	discovery bool
	sudo      = "/usr/bin/sudo"
	sentryBin = "/usr/bin/sentry"
)

func main() {
	usage := `zabbix-agent-extension-sentry

Usage:
  zabbix-agent-extension-sentry [options]

Options:
  -z --zabbix <zabbix>        Hostname or IP address of zabbix server
                                  [default: 127.0.0.1].
  -p --port <port>            Port of zabbix server [default: 10051].
  --zabbix-prefix <prefix>    Add part of your prefix for key [default: None].
  -d --discovery              Run low-level discovery for determine disks.

Discovery options:
  --organizations             Discovery organization.
  --projects                  Discovery projects.

Event options:
  -s --sentry <dsn>           Sentry DSN [default: http://localhost].
  -e --endpoint <epn>         Endpoint API [default: /api/0/].
  -t --token <tkn>            Sentry access token.

Other:
  -h --help                   Show this screen.
`

	var statTypes = []sentry.StatQuery{
		"received",
		"rejected",
		"blacklisted",
	}

	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		fmt.Println(hierr.Errorf(err, "can't parse docopt"))
		os.Exit(1)
	}

	zabbix := args["--zabbix"].(string)
	port, err := strconv.Atoi(args["--port"].(string))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	zabbixPrefix := args["--zabbix-prefix"].(string)
	if zabbixPrefix == "None" {
		zabbixPrefix = "sentry"
	} else {
		zabbixPrefix = fmt.Sprintf("%s.%s", zabbixPrefix, "sentry")
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sentryDsn := args["--sentry"].(string)
	sentryEpn := args["--endpoint"].(string)
	sentryToken, ok := args["--token"].(string)
	if !ok {
		fmt.Println("need setup --token <tkn>")
		os.Exit(1)
	}

	sentryApi := fmt.Sprintf("%s%s", sentryDsn, sentryEpn)

	client, err := sentry.NewClient(sentryToken, &sentryApi, nil)
	if err != nil {
		fmt.Println(
			hierr.Errorf(
				err,
				"can't connect sentry %s with token %s.",
				sentryApi,
				sentryToken,
			),
		)
	}

	organizations, _, err := client.GetOrganizations()
	if err != nil {
		fmt.Println(hierr.Errorf(err, "can't fetch organizations"))
		os.Exit(1)
	}

	projects, err := client.GetProjects()
	if err != nil {
		fmt.Println(hierr.Errorf(err, "can't fetch projects"))
		os.Exit(1)
	}

	if args["--discovery"].(bool) {
		if args["--organizations"].(bool) {
			if err := discoveryOrganizations(organizations); err != nil {
				fmt.Println(hierr.Errorf(err, "can't discovery organizations"))
			}
		}
		if args["--projects"].(bool) {
			if err := discoveryProjects(projects); err != nil {
				fmt.Println(hierr.Errorf(err, "can't discovery projects"))
			}
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
				fmt.Println(
					hierr.Errorf(err, "can't get organization stats"),
				)
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
					fmt.Println(
						hierr.Errorf(err, "can't get project stats"),
					)
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

	fmt.Println("OK")
}
