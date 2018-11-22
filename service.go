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
		return cli.NewExitError("You must set the --cluster-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}
	if c.String("service-name") == "" {
		return cli.NewExitError("You must set the --service-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}

	// Get Ambari connection
	ambariClient := client.New(ambariUrl+"/api/v1", ambariLogin, ambariPassword)
	ambariClient.DisableVerifySSL()

	// Check service alertes
	alerts, err := ambariClient.AlertsInService(c.String("cluster-name"), strings.ToUpper(c.String("service-name")))
	if err != nil {
		monitoringData.AddMessage("Something went wrong when trying to check service alerts on %s: %v", c.String("service-name"), err)
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

	params := &OptionnalComputeState{}
	if c.Bool("exclude-node-alerts") {
		params.Scopes = []string{"SERVICE"}
	}
	if c.String("include-alerts") != "" {
		params.IncludeAlerts = strings.Split(c.String("include-alerts"), ",")
	}
	if c.String("exclude-alerts") != "" {
		params.ExcludeAlerts = strings.Split(c.String("exclude-alerts"), ",")
	}
	monitoringData, err = computeState(filterAlerts, monitoringData, params)
	if err != nil {
		return err
	}

	monitoringData.ToSdtOut()
	return nil

}
