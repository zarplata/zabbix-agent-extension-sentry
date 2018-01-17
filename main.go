package main

import (
	"fmt"
	"os"
	"strconv"

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
  --event                     Web event check.
  --queue                     Cli queue check.

Common options:
  -z --zabbix <zabbix>        Hostname or IP address of zabbix server
                                  [default: 127.0.0.1].
  -p --port <port>            Port of zabbix server [default: 10051].
  --zabbix-prefix <prefix>    Add part of your prefix for key [default: None].
  -d --discovery              Run low-level discovery for determine disks.

Event options:
  -s --sentry <dsn>           Sentry DSN [default: http://localhost].
  -e --endpoint <epn>         Endpoint API [default: /api/0/].
  -o --organization <org>     Sentry organization [default: sentry].
  -t --token <tkn>            Sentry access token.

Other:
  -h --help                   Show this screen.
`
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

	if args["--discovery"].(bool) {
		discovery = true
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if args["--event"].(bool) {
		sentryDsn := args["--sentry"].(string)
		sentryEpn := args["--endpoint"].(string)
		sentryOrg := args["--organization"].(string)
		sentryToken, ok := args["--token"].(string)
		if !ok {
			fmt.Println("need setup --token <tkn>")
			os.Exit(1)
		}

		sentryApi := fmt.Sprintf("%s%s", sentryDsn, sentryEpn)

		err = event(
			sentryApi,
			sentryOrg,
			sentryToken,
			discovery,
			hostname,
			zabbix,
			port,
			zabbixPrefix,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if args["--queue"].(bool) {
		err = queue(
			discovery,
			hostname,
			zabbix,
			port,
			zabbixPrefix,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("OK")
}
