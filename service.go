package main

import (
	"fmt"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

// Perform a service check
func checkService(c *cli.Context) error {

	monitoringData := nagiosPlugin.NewMonitoring()

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), nagiosPlugin.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set --cluster-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}
	if c.String("service-name") == "" {
		return cli.NewExitError("You must set --service-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}

	// Get Ambari connection
	ambariClient := client.New(ambariUrl+"/api/v1", ambariLogin, ambariPassword)
	ambariClient.DisableVerifySSL()

	// Check service alertes
	alerts, err := ambariClient.AlertsInService(c.String("cluster-name"), strings.ToUpper(c.String("service-name")))
	if err != nil {
		monitoringData.AddMessage("Somethink wrong when try to check service alerts on %s: %v", c.String("service-name"), err)
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.ToSdtOut()
	}

	// Remove all UNKNOWN alert before to compute them
	filterAlerts := make([]client.Alert, 0, 1)
	for _, alert := range alerts {
		if alert.AlertInfo.State != "UNKNOWN" {
			filterAlerts = append(filterAlerts, alert)
		}
	}

	if c.Bool("exclude-node-alerts") {
		monitoringData = computeState(filterAlerts, monitoringData, []string{"SERVICE"})
	} else {
		monitoringData = computeState(filterAlerts, monitoringData, []string{})
	}

	monitoringData.ToSdtOut()
	return nil

}
