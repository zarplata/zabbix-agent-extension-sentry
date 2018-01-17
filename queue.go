package main

import (
	"os"
	"os/exec"
	"strings"

	zsend "github.com/blacked/go-zabbix"
	hierr "github.com/reconquest/hierr-go"
)

func getQueueRawOut() (string, error) {

	out, err := exec.Command(sudo, sentryBin, "queues", "list").CombinedOutput()
	if err != nil {
		return "", hierr.Errorf(
			out,
			"can't run command %s %s queues list.",
			sudo,
			sentryBin,
		)
	}
	return string(out), nil
}

func getQueueEvent() ([]map[string]string, error) {

	var queueNames []map[string]string

	getQueueStats, err := getQueueRawOut()
	if err != nil {
		return queueNames, err
	}

	lines := strings.Split(getQueueStats, "\n")

	for _, line := range lines {
		if len(line) > 0 {
			line := strings.Split(line, " ")
			queueName := make(map[string]string)
			queueName["queue"] = line[0]
			queueName["event"] = line[1]
			queueNames = append(queueNames, queueName)
		}
	}
	return queueNames, nil
}

func queue(
	discovery bool,
	hostname string,
	zabbix string,
	port int,
	zabbixPrefix string,
) error {

	if discovery {
		err := discoveryQueue()
		if err != nil {
			return hierr.Errorf(err, "can't discovery projects.")
		}
		os.Exit(0)
	}

	queueNames, err := getQueueEvent()
	if err != nil {
		return err
	}

	var metrics []*zsend.Metric

	for _, queueName := range queueNames {

		metrics = createQueueMetrics(
			hostname,
			metrics,
			queueName,
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
